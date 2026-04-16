import json
import os
import sys
import tempfile
import urllib.error
import urllib.request
import base64
from dataclasses import dataclass
from pathlib import Path
from typing import Any

sys.path.insert(0, str(Path(__file__).resolve().parents[2]))
from scaffolders import (  # noqa: E402
    build_scaffold_output,
    infer_module_base_from_repo_url,
    normalize_module_path,
    read_int,
    read_optional_str,
    read_payload,
    read_required_str,
    resolve_container_image,
    resolve_service_dir,
    scaffold_from_directory,
    split_container_image,
    success,
    validate_service_name,
)

SCHEMA_PATH = Path(__file__).with_name("schema.json")
TEMPLATE_NAME = "template"
LOCAL_TEMPLATE_DIR = Path(__file__).with_name("template")
DEFAULT_ENVIRONMENT = "dev"
DEFAULT_NAMESPACE = "devhub"
DEFAULT_TARGET_REVISION = "main"
DEFAULT_ARGOCD_PROJECT = "default"
DEFAULT_REGISTRY_URL = "host.docker.internal:5001"
DEFAULT_SERVER_URL = "http://host.docker.internal:3000"
DEFAULT_GITOPS_BRANCH = "main"
DEFAULT_GITOPS_BASE_PATH = "envs"
DEFAULT_GITOPS_COMMIT_USER_NAME = "devhub-bot"
DEFAULT_GITOPS_COMMIT_USER_EMAIL = "devhub-bot@local"
API_PREFIX = "/api/v1"


@dataclass(frozen=True)
class ScaffoldPayload:
    service_name: str
    project_id: str
    repo_url: str
    template: str
    environment: str
    namespace: str
    target_revision: str
    argocd_project: str
    registry_url: str
    server_url: str
    module_path: str
    port: int
    image: str

    @classmethod
    def from_dict(cls, payload: dict[str, Any], properties: dict[str, Any]) -> "ScaffoldPayload":
        service_name = read_required_str(payload, "service_name")
        validate_service_name(service_name)
        repo_url = read_required_str(payload, "repo_url")

        return cls(
            service_name=service_name,
            project_id=read_optional_str(payload, "project_id"),
            repo_url=repo_url,
            template=read_optional_str(payload, "template", TEMPLATE_NAME),
            environment=read_optional_str(payload, "environment", DEFAULT_ENVIRONMENT),
            namespace=read_optional_str(payload, "namespace", DEFAULT_NAMESPACE),
            target_revision=read_optional_str(payload, "target_revision", DEFAULT_TARGET_REVISION),
            argocd_project=read_optional_str(payload, "argocd_project", DEFAULT_ARGOCD_PROJECT),
            registry_url=read_optional_str(payload, "registry_url", DEFAULT_REGISTRY_URL),
            server_url=read_optional_str(payload, "server_url", DEFAULT_SERVER_URL),
            module_path=read_module_base(payload, properties),
            port=read_port(payload, properties),
            image=resolve_container_image(payload, service_name),
        )

    def to_output_payload(self) -> dict[str, Any]:
        return {
            "service_name": self.service_name,
            "project_id": self.project_id,
            "repo_url": self.repo_url,
            "template": self.template,
            "environment": self.environment,
            "namespace": self.namespace,
            "target_revision": self.target_revision,
            "argocd_project": self.argocd_project,
            "registry_url": self.registry_url,
            "server_url": self.server_url,
            "module_path": self.module_path,
            "port": self.port,
            "image": self.image,
        }


def load_schema() -> dict[str, Any]:
    return json.loads(SCHEMA_PATH.read_text(encoding="utf-8"))


def read_port(payload: dict[str, Any], properties: dict[str, Any]) -> int:
    default_port = properties.get("port", {}).get("default", 8080)
    return read_int(
        payload,
        "port",
        default=default_port,
        min_value=1,
        max_value=65535,
    )


def read_module_base(payload: dict[str, Any], properties: dict[str, Any]) -> str:
    explicit_module = read_optional_str(payload, "module_path")
    if explicit_module != "":
        return explicit_module

    inferred_module = infer_module_base_from_repo_url(read_optional_str(payload, "repo_url"))
    if inferred_module != "":
        return inferred_module

    return str(properties.get("module_path", {}).get("default", "github.com/acme")).strip()


def parse_payload(schema: dict[str, Any]) -> dict[str, Any]:
    required_fields = schema.get("required", ["service_name"])
    return read_payload(required_fields=required_fields)


def build_template_context(payload: ScaffoldPayload) -> dict[str, str]:
    image_repository, image_tag = split_container_image(payload.image)

    return {
        "SERVICE_NAME": payload.service_name,
        "MODULE_PATH": normalize_module_path(payload.module_path, payload.service_name),
        "PORT": str(payload.port),
        "IMAGE": payload.image,
        "IMAGE_REPOSITORY": image_repository,
        "IMAGE_TAG": image_tag,
        "REPO_URL": payload.repo_url,
        "TARGET_REVISION": payload.target_revision,
        "NAMESPACE": payload.namespace,
        "ARGOCD_PROJECT": payload.argocd_project,
        "REGISTRY_URL": payload.registry_url,
        "SERVER_URL": payload.server_url,
        "ENVIRONMENT": payload.environment,
    }


def normalize_api_base_url(raw: str) -> str:
    value = raw.strip().rstrip("/")
    if value == "":
        return ""
    if value.endswith(API_PREFIX):
        return value
    return f"{value}{API_PREFIX}"


def http_request(method: str, url: str, token: str, body: dict[str, Any] | None = None) -> tuple[int, str]:
    data = None
    headers = {"Accept": "application/json"}

    if token:
        headers["Authorization"] = f"token {token}"

    if body is not None:
        data = json.dumps(body).encode("utf-8")
        headers["Content-Type"] = "application/json"

    request = urllib.request.Request(url, data=data, headers=headers, method=method)
    with urllib.request.urlopen(request) as response:
        return response.getcode(), response.read().decode("utf-8")


def get_file(api_base_url: str, token: str, owner: str, repo_name: str, path: str, branch: str) -> dict[str, Any] | list[dict[str, Any]] | None:
    url = f"{api_base_url}/repos/{owner}/{repo_name}/contents/{path}?ref={branch}"
    try:
        status, body = http_request("GET", url, token)
    except urllib.error.HTTPError as exc:
        if exc.code == 404:
            return None
        raise

    if status < 200 or status >= 300:
        raise RuntimeError(f"get file failed: {status} {body}")

    parsed = json.loads(body)
    if isinstance(parsed, dict) or isinstance(parsed, list):
        return parsed
    raise RuntimeError(f"unexpected contents response type for path={path}")


def upsert_file(
    api_base_url: str,
    token: str,
    owner: str,
    repo_name: str,
    path: str,
    branch: str,
    content: str,
    message: str,
    author_name: str,
    author_email: str,
) -> None:
    existing = get_file(api_base_url, token, owner, repo_name, path, branch)
    if isinstance(existing, list):
        raise RuntimeError(f"gitops path conflict: {path} resolves to a directory, expected a file")

    body: dict[str, Any] = {
        "branch": branch,
        "content": base64.b64encode(content.encode("utf-8")).decode("utf-8"),
        "message": message,
        "author": {"name": author_name, "email": author_email},
        "committer": {"name": author_name, "email": author_email},
    }
    if isinstance(existing, dict):
        body["sha"] = str(existing.get("sha", "")).strip()

    url = f"{api_base_url}/repos/{owner}/{repo_name}/contents/{path}"
    status, response = http_request("PUT", url, token, body)
    if status < 200 or status >= 300:
        raise RuntimeError(f"upsert file failed: {status} {response}")


def resolve_gitops_values_path(
    api_base_url: str,
    token: str,
    owner: str,
    repo_name: str,
    branch: str,
    base_path: str,
    service_name: str,
) -> str:
    primary_path = f"{base_path}/{service_name}.yaml"
    candidate_paths = [
        primary_path,
        f"{primary_path}/values.yaml",
        f"{primary_path}/app/values.yaml",
    ]

    last_directory_path = ""
    for index, candidate_path in enumerate(candidate_paths):
        existing = get_file(api_base_url, token, owner, repo_name, candidate_path, branch)
        if isinstance(existing, list):
            last_directory_path = candidate_path
            continue

        if existing is not None:
            return candidate_path

        if index == 0:
            return candidate_path

    if last_directory_path != "":
        raise RuntimeError(
            f"gitops path conflict: {last_directory_path} resolves to a directory, expected a file"
        )

    return candidate_paths[-1]


def build_gitops_values_content(payload: ScaffoldPayload) -> str:
    image_repository, image_tag = split_container_image(payload.image)
    return f"""nameOverride: "{payload.service_name}"
fullnameOverride: "{payload.service_name}"

replicaCount: 1

image:
  repository: "{image_repository}"
  tag: "{image_tag}"
  pullPolicy: IfNotPresent

service:
  enabled: true
  type: ClusterIP
  port: {payload.port}

containerPort: {payload.port}

ingress:
  enabled: false
  className: ""
  host: "{payload.service_name}.devhub.local"
  path: /

serviceAccount:
  create: true
  name: ""

config:
  enabled: false
  values: {{}}
"""


def maybe_bootstrap_gitops(payload: ScaffoldPayload) -> None:
    api_base_url = normalize_api_base_url(os.getenv("SCM_API_URL", ""))
    token = os.getenv("SCM_TOKEN", "").strip()
    owner = os.getenv("GITOPS_REPO_OWNER", "").strip()
    repo_name = os.getenv("GITOPS_REPO_NAME", "").strip()

    if api_base_url == "" or token == "" or owner == "" or repo_name == "":
        return

    branch = os.getenv("GITOPS_BRANCH", DEFAULT_GITOPS_BRANCH).strip() or DEFAULT_GITOPS_BRANCH
    base_path = os.getenv("GITOPS_BASE_PATH", DEFAULT_GITOPS_BASE_PATH).strip().strip("/") or DEFAULT_GITOPS_BASE_PATH
    author_name = os.getenv("GITOPS_COMMIT_USER_NAME", DEFAULT_GITOPS_COMMIT_USER_NAME).strip() or DEFAULT_GITOPS_COMMIT_USER_NAME
    author_email = os.getenv("GITOPS_COMMIT_USER_EMAIL", DEFAULT_GITOPS_COMMIT_USER_EMAIL).strip() or DEFAULT_GITOPS_COMMIT_USER_EMAIL

    values_path = resolve_gitops_values_path(
        api_base_url,
        token,
        owner,
        repo_name,
        branch,
        base_path,
        payload.service_name,
    )
    upsert_file(
        api_base_url,
        token,
        owner,
        repo_name,
        values_path,
        branch,
        build_gitops_values_content(payload),
        f"bootstrap gitops values for {payload.service_name}",
        author_name,
        author_email,
    )


def run() -> None:
    schema = load_schema()
    payload_dict = parse_payload(schema)
    payload = ScaffoldPayload.from_dict(payload_dict, schema.get("properties", {}))

    with tempfile.TemporaryDirectory(prefix=f"{payload.service_name}-") as temp_dir:
        service_dir = resolve_service_dir(temp_dir, payload.service_name)
        template_context = build_template_context(payload)

        scaffold_from_directory(service_dir, LOCAL_TEMPLATE_DIR, template_context)
        maybe_bootstrap_gitops(payload)

        success(build_scaffold_output(service_dir, payload.service_name, payload.to_output_payload()))


if __name__ == "__main__":
    run()
