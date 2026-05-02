import os
from dataclasses import dataclass
from typing import Any

from scaffolders import read_int, read_optional_str, read_required_str  # noqa: E402


DEFAULT_GITOPS_BRANCH = "main"
DEFAULT_GITOPS_BASE_PATH = "envs"
DEFAULT_GITOPS_COMMIT_USER_NAME = "devhub-bot"
DEFAULT_GITOPS_COMMIT_USER_EMAIL = "devhub-bot@local"


@dataclass(frozen=True)
class ScaffoldPayload:
    environment: str
    service_name: str
    port: int
    database: str
    image_tag: str
    module_path: str
    ci_registry_host: str
    ci_server_url: str
    cd_project_name: str
    cd_repo_url: str
    cd_target_revision: str
    cd_namespace: str
    cd_image_repository: str

    @classmethod
    def from_dict(cls, payload: dict[str, Any]):
        return cls(
            environment=read_required_str(payload, "environment"),
            service_name=read_required_str(payload, "service_name"),
            port=read_int(payload, "port", default=0, min_value=1, max_value=65535),
            database=read_required_str(payload, "database"),
            image_tag=read_optional_str(payload, "image_tag", "latest"),
            module_path=read_required_str(payload, "module_path"),
            ci_registry_host=read_required_str(payload, "ci_registry_host"),
            ci_server_url=read_required_str(payload, "ci_server_url"),
            cd_project_name=read_required_str(payload, "cd_project_name"),
            cd_repo_url=read_required_str(payload, "cd_repo_url"),
            cd_target_revision=read_required_str(payload, "cd_target_revision"),
            cd_namespace=read_required_str(payload, "cd_namespace"),
            cd_image_repository=read_required_str(payload, "cd_image_repository"),
        )
    
    def to_template(self):
        return {
            "SERVICE_NAME": self.service_name,
            "MODULE_PATH": self.module_path,
            "PORT": str(self.port),
            "IMAGE_TAG": self.image_tag,
            "ENVIRONMENT": self.environment,
            "CD_PROJECT_NAME": self.cd_project_name,
            "CD_IMAGE_REPOSITORY": self.cd_image_repository,
            "CD_REPO_URL": self.cd_repo_url,
            "CD_TARGET_REVISION": self.cd_target_revision,
            "CD_NAMESPACE": self.cd_namespace,
            "CI_REGISTRY_HOST": self.ci_registry_host,
            "CI_SERVER_URL": self.ci_server_url,
        }

    def to_dict(self):
        return {
            "environment": self.environment,
            "service_name": self.service_name,
            "port": self.port,
            "database": self.database,
            "image_tag": self.image_tag,
            "module_path": self.module_path,
            "ci_registry_host": self.ci_registry_host,
            "ci_server_url": self.ci_server_url,
            "cd_project_name": self.cd_project_name,
            "cd_repo_url": self.cd_repo_url,
            "cd_target_revision": self.cd_target_revision,
            "cd_namespace": self.cd_namespace,
            "cd_image_repository": self.cd_image_repository,
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
