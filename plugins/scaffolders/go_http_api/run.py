import json
import sys
import tempfile
import urllib.error
import urllib.request
import base64
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[2]))
from payload import GitOpsConfig, ScaffoldPayload
from scaffolders import (  # noqa: E402
    build_scaffold_output,
    normalize_module_path,
    read_payload,
    render_template,
    resolve_service_dir,
    scaffold_from_directory,
    split_container_image,
    success,
    read_optional_str,
    infer_module_base_from_repo_url,
)

SCHEMA_PATH = Path(__file__).with_name("schema.json")
LOCAL_TEMPLATE_DIR = Path(__file__).with_name("template")
VALUES_TEMPLATE_PATH = LOCAL_TEMPLATE_DIR / "deploy" / "helm" / "values.yaml"

DEFAULT_ENVIRONMENT = "dev"
DEFAULT_NAMESPACE = "devhub"
DEFAULT_TARGET_REVISION = "main"
DEFAULT_ARGOCD_PROJECT = "default"
DEFAULT_SERVER_URL = "http://host.docker.internal:3000"
DEFAULT_REGISTRY_URL = "host.docker.internal:5001"
DEFAULT_HELM_SERVER_URL = "http://host.minikube.internal:3000"
DEFAULT_HELM_REGISTRY_URL = "host.minikube.internal:5001"

DEFAULT_GITOPS_BRANCH = "main"
DEFAULT_GITOPS_BASE_PATH = "envs"
DEFAULT_GITOPS_COMMIT_USER_NAME = "devhub-bot"
DEFAULT_GITOPS_COMMIT_USER_EMAIL = "devhub-bot@local"

class GitOpsClient:
    def __init__(self, config: GitOpsConfig):
        self.config = config

    def _request(self, method: str, path: str, body: dict | None = None):
        url = f"{self.config.api_base_url}/repos/{self.config.owner}/{self.config.repo_name}/contents/{path}"

        headers = {"Accept": "application/json"}
        if self.config.token:
            headers["Authorization"] = f"token {self.config.token}"

        data = None
        if body:
            data = json.dumps(body).encode()
            headers["Content-Type"] = "application/json"

        req = urllib.request.Request(url, data=data, headers=headers, method=method)

        try:
            with urllib.request.urlopen(req) as resp:
                return resp.getcode(), resp.read().decode()
        except urllib.error.HTTPError as e:
            error_body = e.read().decode()
            raise RuntimeError(f"HTTP {e.code}: {error_body}") from e

    def get_file(self, path: str):
        try:
            _, body = self._request("GET", f"{path}?ref={self.config.branch}")
        except RuntimeError as e:
            if "404" in str(e):
                return None
            raise

        data = json.loads(body)

        if isinstance(data, dict) and data.get("sha"):
            return data

        return None

    def create_file(self, path: str, content: str, message: str):
        body = {
            "branch": self.config.branch,
            "content": base64.b64encode(content.encode()).decode(),
            "message": message,
            "author": {
                "name": self.config.author_name,
                "email": self.config.author_email,
            },
            "committer": {
                "name": self.config.author_name,
                "email": self.config.author_email,
            },
        }

        return self._request("POST", path, body)

    def update_file(self, path: str, content: str, sha: str, message: str):
        body = {
            "branch": self.config.branch,
            "content": base64.b64encode(content.encode()).decode(),
            "message": message,
            "sha": sha,
            "author": {
                "name": self.config.author_name,
                "email": self.config.author_email,
            },
            "committer": {
                "name": self.config.author_name,
                "email": self.config.author_email,
            },
        }

        return self._request("PUT", path, body)

    def save_file(self, path: str, content: str, message: str):
        try:
            self.create_file(path, content, message)
            return
        except RuntimeError as e:
            if "[SHA]" not in str(e):
                raise

        existing = self.get_file(path)
        if not existing or "sha" not in existing:
            raise RuntimeError(f"Cannot resolve sha for {path}")

        self.update_file(path, content, existing["sha"], message)


def load_schema():
    return json.loads(SCHEMA_PATH.read_text())


def parse_payload(schema):
    return read_payload(required_fields=schema.get("required", ["service_name"]))


def build_template_context(payload: ScaffoldPayload):
    repo, tag = split_container_image(payload.image)

    return {
        "SERVICE_NAME": payload.service_name,
        "MODULE_PATH": normalize_module_path(payload.module_path, payload.service_name),
        "PORT": str(payload.port),
        "IMAGE": payload.image,
        "IMAGE_REPOSITORY": f'payload.registry_url/{repo}',
        "IMAGE_TAG": tag,
        "REPO_URL": payload.repo_url,
        "TARGET_REVISION": payload.target_revision,
        "NAMESPACE": payload.namespace,
        "ARGOCD_PROJECT": payload.argocd_project,
        "REGISTRY_URL": payload.registry_url,
        "SERVER_URL": payload.server_url,
        "ENVIRONMENT": payload.environment,
    }


def build_gitops_values_content(payload: ScaffoldPayload):
    template = VALUES_TEMPLATE_PATH.read_text()
    context = build_template_context(payload)
    
    return render_template(template, context).strip()


def remove_repo_values_file(service_dir: Path):
    repo_values_path = service_dir / "deploy" / "helm" / "values.yaml"
    if repo_values_path.exists():
        repo_values_path.unlink()


def bootstrap_gitops(payload: ScaffoldPayload):
    config = GitOpsConfig.from_env()
    if not config:
        return

    client = GitOpsClient(config)

    path = f"{config.base_path}/{payload.environment}/{payload.service_name}.yaml"
    content = build_gitops_values_content(payload)

    client.save_file(
        path=path,
        content=content,
        message=f"bootstrap gitops values for {payload.service_name}",
    )


def read_module_base(payload, properties):
    explicit = read_optional_str(payload, "module_path")
    if explicit:
        return explicit

    inferred = infer_module_base_from_repo_url(read_optional_str(payload, "repo_url"))
    if inferred:
        return inferred

    return properties.get("module_path", {}).get("default", "github.com/acme")


def run():
    schema = load_schema()
    payload_dict = parse_payload(schema)
    payload = ScaffoldPayload.from_dict(payload_dict, schema.get("properties", {}))

    with tempfile.TemporaryDirectory(prefix=f"{payload.service_name}-") as temp_dir:
        service_dir = resolve_service_dir(temp_dir, payload.service_name)

        context = build_template_context(payload)

        scaffold_from_directory(service_dir, LOCAL_TEMPLATE_DIR, context)
        remove_repo_values_file(service_dir)

        bootstrap_gitops(payload)

        success(build_scaffold_output(service_dir, payload.service_name, payload.to_output_payload()))


if __name__ == "__main__":
    run()
