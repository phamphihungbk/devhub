import base64
import json
import sys
import urllib.error
import urllib.request
from dataclasses import dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

sys.path.insert(0, str(Path(__file__).resolve().parents[2]))
from scaffolders import (  # noqa: E402
    read_optional_str,
    read_payload,
    read_required_str,
    success,
    fail,
)

SCHEMA_PATH = Path(__file__).with_name("schema.json")

DEFAULT_GITOPS_BRANCH = "main"
DEFAULT_GITOPS_BASE_PATH = "envs"
DEFAULT_COMMIT_USER_NAME = "devhub-bot"
DEFAULT_COMMIT_USER_EMAIL = "devhub-bot@local"

def log(msg: str):
    print(msg, file=sys.stderr)


@dataclass(frozen=True)
class DeploymentPayload:
    deployment_id: str
    project_id: str
    service: str
    environment: str
    version: str
    repo_url: str
    scm_api_url: str
    scm_token: str
    gitops_repo_owner: str
    gitops_repo_name: str
    gitops_branch: str
    gitops_base_path: str
    commit_user_name: str
    commit_user_email: str
    argocd_server: str
    argocd_auth_token: str
    argocd_insecure: bool

    @classmethod
    def from_dict(cls, payload: dict[str, Any]) -> "DeploymentPayload":
        return cls(
            deployment_id=read_optional_str(payload, "deployment_id"),
            project_id=read_optional_str(payload, "project_id"),
            service=read_required_str(payload, "service"),
            environment=read_required_str(payload, "environment"),
            version=read_required_str(payload, "version"),
            repo_url=read_optional_str(payload, "repo_url"),
            scm_api_url=read_required_str(payload, "scm_api_url").rstrip("/"),
            scm_token=read_required_str(payload, "scm_token"),
            gitops_repo_owner=read_required_str(payload, "gitops_repo_owner"),
            gitops_repo_name=read_required_str(payload, "gitops_repo_name"),
            gitops_branch=read_optional_str(payload, "gitops_branch", DEFAULT_GITOPS_BRANCH),
            gitops_base_path=read_optional_str(payload, "gitops_base_path", DEFAULT_GITOPS_BASE_PATH),
            commit_user_name=read_optional_str(payload, "commit_user_name", DEFAULT_COMMIT_USER_NAME),
            commit_user_email=read_optional_str(payload, "commit_user_email", DEFAULT_COMMIT_USER_EMAIL),
            argocd_server=read_optional_str(payload, "argocd_server"),
            argocd_auth_token=read_optional_str(payload, "argocd_auth_token"),
            argocd_insecure=bool(payload.get("argocd_insecure", False)),
        )


def normalize_url(url: str) -> str:
    if not url:
        return url
    if not url.startswith("http://") and not url.startswith("https://"):
        return f"http://{url}"
    return url


def http_request(method: str, url: str, token: str, body: dict | None = None):
    headers = {"Accept": "application/json"}

    if token:
        headers["Authorization"] = f"token {token}"

    data = None
    if body:
        headers["Content-Type"] = "application/json"
        data = json.dumps(body).encode()

    req = urllib.request.Request(url, data=data, headers=headers, method=method)

    try:
        with urllib.request.urlopen(req) as resp:
            return resp.getcode(), resp.read().decode()
    except urllib.error.HTTPError as e:
        raise RuntimeError(f"HTTP {e.code}: {e.read().decode()}") from e


def get_file(payload: DeploymentPayload, path: str):
    url = (
        f"{payload.scm_api_url}/repos/"
        f"{payload.gitops_repo_owner}/{payload.gitops_repo_name}/contents/{path}"
        f"?ref={payload.gitops_branch}"
    )

    try:
        _, body = http_request("GET", url, payload.scm_token)
    except RuntimeError as e:
        if "404" in str(e):
            return None
        raise

    return json.loads(body)


def update_file(payload: DeploymentPayload, path: str, sha: str, content_b64: str, message: str):
    url = f"{payload.scm_api_url}/repos/{payload.gitops_repo_owner}/{payload.gitops_repo_name}/contents/{path}"

    body = {
        "branch": payload.gitops_branch,
        "content": content_b64,
        "sha": sha,
        "message": message,
        "author": {
            "name": payload.commit_user_name,
            "email": payload.commit_user_email,
        },
        "committer": {
            "name": payload.commit_user_name,
            "email": payload.commit_user_email,
        },
    }

    _, resp = http_request("PUT", url, payload.scm_token, body)
    return json.loads(resp)


def decode_content(encoded: str) -> str:
    return base64.b64decode(encoded.replace("\n", "")).decode()


def encode_content(content: str) -> str:
    return base64.b64encode(content.encode()).decode()


def update_image_tag(yaml_text: str, new_tag: str) -> str:
    lines = yaml_text.splitlines()
    result = []
    in_image = False

    for line in lines:
        stripped = line.strip()

        if stripped.startswith("image:"):
            in_image = True
            result.append(line)
            continue

        if in_image and stripped.startswith("tag:"):
            indent = line[: line.index("tag")]
            result.append(f"{indent}tag: {new_tag}")
            in_image = False
            continue

        result.append(line)

    return "\n".join(result) + "\n"


def sync_argocd(payload: DeploymentPayload, app_name: str):
    if not payload.argocd_server or not payload.argocd_auth_token:
        return

    server = normalize_url(payload.argocd_server)

    url = f"{server}/api/v1/applications/{app_name}/sync"

    headers = {
        "Authorization": f"Bearer {payload.argocd_auth_token}",
        "Content-Type": "application/json",
    }

    body = {"revision": payload.gitops_branch}

    req = urllib.request.Request(
        url,
        data=json.dumps(body).encode(),
        headers=headers,
        method="POST",
    )

    try:
        urllib.request.urlopen(req)
    except urllib.error.HTTPError as e:
        log(f"ArgoCD sync failed (ignored): {e.read().decode()}")


def run():
    schema = json.loads(SCHEMA_PATH.read_text())
    payload_dict = read_payload(required_fields=schema.get("required", []))
    payload = DeploymentPayload.from_dict(payload_dict)

    values_path = f"{payload.gitops_base_path}/{payload.environment}/{payload.service}.yaml"

    scm_file = get_file(payload, values_path)
    if not scm_file:
        fail(f"values file not found: {values_path}")

    sha = scm_file.get("sha")
    content = scm_file.get("content")

    if not sha or not content:
        fail("invalid gitops file response")

    decoded = decode_content(content)
    updated = update_image_tag(decoded, payload.version)
    encoded = encode_content(updated)

    commit_message = f"deploy {payload.service}:{payload.version}"

    result = update_file(payload, values_path, sha, encoded, commit_message)

    commit_sha = result.get("commit", {}).get("sha", "")
    app_name = f"{payload.service}-{payload.environment}"

    sync_argocd(payload, app_name)

    success(
        {
            "external_ref": app_name,
            "commit_sha": commit_sha,
            "finished_at": datetime.now(timezone.utc).isoformat(),
        }
    )


if __name__ == "__main__":
    run()