import os
from dataclasses import dataclass
from typing import Any

from scaffolders import read_int, read_required_str  # noqa: E402


DEFAULT_GITOPS_BRANCH = "main"
DEFAULT_GITOPS_BASE_PATH = "envs"
DEFAULT_GITOPS_COMMIT_USER_NAME = "devhub-bot"
DEFAULT_GITOPS_COMMIT_USER_EMAIL = "devhub-bot@local"


@dataclass(frozen=True)
class ScaffoldPayload:
    service_name: str
    project_id: str
    repo_url: str
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
    def from_dict(cls, payload: dict[str, Any]):
        return cls(
            service_name=read_required_str(payload, "service_name"),
            project_id=read_required_str(payload, "project_id"),
            repo_url=read_required_str(payload, "repo_url"),
            environment=read_required_str(payload, "environment"),
            namespace=read_required_str(payload, "namespace"),
            target_revision=read_required_str(payload, "target_revision"),
            argocd_project=read_required_str(payload, "argocd_project"),
            registry_url=read_required_str(payload, "registry_url"),
            server_url=read_required_str(payload, "server_url"),
            module_path=read_required_str(payload, "module_path"),
            port=read_int(payload, "port", default=0, min_value=1, max_value=65535),
            image=read_required_str(payload, "image"),
        )

    def to_dict(self):
        return {
            "service_name": self.service_name,
            "project_id": self.project_id,
            "repo_url": self.repo_url,
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


@dataclass(frozen=True)
class GitOpsConfig:
    api_base_url: str
    token: str
    gitops_owner: str
    gitops_repo: str
    branch: str
    base_path: str
    author_name: str
    author_email: str

    @classmethod
    def from_env(cls):
        api_base_url = os.getenv("SCM_API_URL", "").strip().rstrip("/")
        token = os.getenv("SCM_TOKEN", "").strip()
        gitops_owner = os.getenv("GITOPS_REPO_OWNER", "").strip()
        gitops_repo = os.getenv("GITOPS_REPO_NAME", "").strip()

        if not api_base_url or not token or not gitops_owner or not gitops_repo:
            return None

        return cls(
            api_base_url=api_base_url,
            token=token,
            gitops_owner=gitops_owner,
            gitops_repo=gitops_repo,
            branch=os.getenv("GITOPS_BRANCH", DEFAULT_GITOPS_BRANCH).strip() or DEFAULT_GITOPS_BRANCH,
            base_path=os.getenv("GITOPS_BASE_PATH", DEFAULT_GITOPS_BASE_PATH).strip("/"),
            author_name=os.getenv("GITOPS_COMMIT_USER_NAME", DEFAULT_GITOPS_COMMIT_USER_NAME).strip()
            or DEFAULT_GITOPS_COMMIT_USER_NAME,
            author_email=os.getenv("GITOPS_COMMIT_USER_EMAIL", DEFAULT_GITOPS_COMMIT_USER_EMAIL).strip()
            or DEFAULT_GITOPS_COMMIT_USER_EMAIL,
        )
