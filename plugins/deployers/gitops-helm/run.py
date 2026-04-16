import base64
import json
import os
import subprocess
import sys
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
)

SCHEMA_PATH = Path(__file__).with_name("schema.json")

DEFAULT_GITOPS_BRANCH = "main"
DEFAULT_GITOPS_BASE_PATH = "envs"
DEFAULT_COMMIT_USER_NAME = "devhub-bot"
DEFAULT_COMMIT_USER_EMAIL = "devhub-bot@local"


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


def load_schema() -> dict[str, Any]:
    return json.loads(SCHEMA_PATH.read_text(encoding="utf-8"))


def parse_payload(schema: dict[str, Any]) -> dict[str, Any]:
    required_fields = schema.get("required", ["service", "environment", "version"])
    return read_payload(required_fields=required_fields)


def build_values_path(payload: DeploymentPayload) -> str:
    return f"{payload.gitops_base_path.strip('/')}/{payload.environment}/{payload.service}.yaml"


def build_commit_message(payload: DeploymentPayload) -> str:
    return f"deploy {payload.service} to {payload.environment} with image tag {payload.version}"


def http_request(method: str, url: str, token: str, body: dict[str, Any] | None = None) -> tuple[int, str]:
    data = None
    headers = {
        "Accept": "application/json",
    }

    if token:
        headers["Authorization"] = f"token {token}"

    if body is not None:
        data = json.dumps(body).encode("utf-8")
        headers["Content-Type"] = "application/json"

    req = urllib.request.Request(url, data=data, method=method, headers=headers)

    with urllib.request.urlopen(req) as resp:
        return resp.getcode(), resp.read().decode("utf-8")


def get_file(payload: DeploymentPayload, path: str) -> dict[str, Any]:
    url = (
        f"{payload.scm_api_url}/repos/"
        f"{payload.gitops_repo_owner}/{payload.gitops_repo_name}/contents/{path}"
        f"?ref={payload.gitops_branch}"
    )
    status, body = http_request("GET", url, payload.scm_token)
    if status < 200 or status >= 300:
        raise RuntimeError(f"get file failed: {status} {body}")
    return json.loads(body)


def update_file(payload: DeploymentPayload, path: str, sha: str, content_b64: str, message: str) -> dict[str, Any]:
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

    status, response = http_request("PUT", url, payload.scm_token, body)
    if status < 200 or status >= 300:
        raise RuntimeError(f"update file failed: {status} {response}")
    return json.loads(response)


def decode_content(encoded: str) -> str:
    return base64.b64decode(encoded.replace("\n", "")).decode("utf-8")


def encode_content(content: str) -> str:
    return base64.b64encode(content.encode("utf-8")).decode("utf-8")


def update_image_tag_in_yaml(yaml_text: str, new_tag: str) -> str:
    lines = yaml_text.splitlines()
    updated: list[str] = []
    in_image = False

    for line in lines:
        stripped = line.strip()

        if stripped.startswith("image:"):
            in_image = True
            updated.append(line)
            continue

        if in_image and stripped.startswith("tag:"):
            indent = line[: line.index("tag")]
            updated.append(f"{indent}tag: {new_tag}")
            in_image = False
            continue

        updated.append(line)

    return "\n".join(updated) + "\n"


def maybe_sync_argocd(payload: DeploymentPayload, app_name: str) -> None:
    if payload.argocd_server == "" or payload.argocd_auth_token == "":
        return

    base_cmd = [
        "argocd",
        "--server",
        payload.argocd_server,
        "--auth-token",
        payload.argocd_auth_token,
    ]

    if payload.argocd_insecure:
        base_cmd.append("--insecure")

    refresh_cmd = base_cmd + ["app", "get", app_name, "--refresh"]
    sync_cmd = base_cmd + ["app", "sync", app_name]

    try:
        subprocess.run(
            refresh_cmd,
            capture_output=True,
            text=True,
            check=False,
            timeout=60,
        )

        sync_proc = subprocess.run(
            sync_cmd,
            capture_output=True,
            text=True,
            check=False,
            timeout=300,
        )

        if sync_proc.returncode != 0:
            stderr = (sync_proc.stderr or "").strip()
            stdout = (sync_proc.stdout or "").strip()
            raise RuntimeError(
                f"argocd app sync failed for app={app_name!r}; "
                f"stdout={stdout}; stderr={stderr}"
            )

    except subprocess.TimeoutExpired as exc:
        raise RuntimeError(f"argocd command timed out for app={app_name!r}") from exc
    except FileNotFoundError as exc:
        raise RuntimeError(
            "argocd CLI not found. Ensure it is installed in the plugin runtime/container."
        ) from exc


def run() -> None:
    schema = load_schema()
    payload_dict = parse_payload(schema)
    payload = DeploymentPayload.from_dict(payload_dict)

    values_path = build_values_path(payload)
    scm_file = get_file(payload, values_path)

    sha = str(scm_file.get("sha", "")).strip()
    content = str(scm_file.get("content", "")).strip()

    if sha == "" or content == "":
        raise RuntimeError("invalid gitops file response: missing sha or content")

    decoded_yaml = decode_content(content)
    updated_yaml = update_image_tag_in_yaml(decoded_yaml, payload.version)
    encoded_yaml = encode_content(updated_yaml)

    commit_message = build_commit_message(payload)
    update_result = update_file(payload, values_path, sha, encoded_yaml, commit_message)

    commit_sha = str(update_result.get("commit", {}).get("sha", "")).strip()
    external_ref = f"{payload.service}-{payload.environment}"

    maybe_sync_argocd(payload, external_ref)

    success(
        {
            "external_ref": external_ref,
            "commit_sha": commit_sha,
            "finished_at": datetime.now(timezone.utc).isoformat(),
        }
    )


if __name__ == "__main__":
    run()