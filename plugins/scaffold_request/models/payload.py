from dataclasses import dataclass
import re
from typing import Any

from scaffold_request import read_int, read_optional_str, read_required_str  # noqa: E402
from utils import read_required_env  # noqa: E402


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
    argocd_repo_url: str
    cd_target_revision: str
    cd_namespace: str
    team_owner: str
    scm_external_url: str
    argocd_repository_registry_host: str

    @classmethod
    def from_dict(cls, payload: dict[str, Any]) -> "ScaffoldPayload":
        service_name = re.sub(r"[\s_]+", "-", read_required_str(payload, "service_name").lower()).strip("-")

        return cls(
            environment=read_required_str(payload, "environment"),
            service_name=service_name,
            team_owner=read_required_str(payload, "team_owner"),
            port=read_int(payload, "port", default=0, min_value=1, max_value=65535),
            database=read_required_str(payload, "database"),
            image_tag=read_optional_str(payload, "image_tag", "latest"),
            module_path=read_required_str(payload, "module_path"),
            ci_registry_host=read_required_env("CI_REGISTRY_HOST"),
            ci_server_url=read_required_env("CI_SERVER_URL"),
            cd_project_name=read_required_env("CD_PROJECT_NAME"),
            cd_target_revision=read_required_env("CD_TARGET_REVISION"),
            cd_namespace=read_required_env("CD_NAMESPACE"),
            scm_external_url=read_required_env("SCM_EXTERNAL_URL"),
            argocd_repo_url=read_required_env("ARGOCD_REPO_URL"),
            argocd_repository_registry_host=read_required_env("ARGOCD_REPOSITORY_REGISTRY_HOST"),
        )

    @property
    def cd_repo_url(self) -> str:
        return "/".join([self.argocd_repo_url.rstrip("/"), self.team_owner, self.service_name]) + ".git"

    @property
    def scm_repo_url(self) -> str:
        return "/".join([self.scm_external_url.rstrip("/"), self.team_owner, self.service_name]) + ".git"

    @property
    def cd_image_repository(self) -> str:
        return f"{self.argocd_repository_registry_host.rstrip('/')}/{self.service_name}"

    def to_template(self) -> dict[str, str]:
        return {
            "SERVICE_NAME": self.service_name,
            "MODULE_PATH": self.module_path,
            "PORT": str(self.port),
            "IMAGE_TAG": self.image_tag,
            "ENVIRONMENT": self.environment,
            "CD_PROJECT_NAME": self.cd_project_name,
            "CD_IMAGE_REPOSITORY": self.cd_image_repository,
            "CD_REPO_URL": self.cd_repo_url,
            "SCM_REPO_URL": self.scm_repo_url,
            "CD_TARGET_REVISION": self.cd_target_revision,
            "CD_NAMESPACE": self.cd_namespace,
            "CI_REGISTRY_HOST": self.ci_registry_host,
            "CI_SERVER_URL": self.ci_server_url,
        }

    def to_dict(self) -> dict[str, Any]:
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
            "cd_target_revision": self.cd_target_revision,
            "cd_namespace": self.cd_namespace,
            "cd_repo_url": self.cd_repo_url,
            "scm_repo_url": self.scm_repo_url,
            "cd_image_repository": self.cd_image_repository,
            "scm_external_url": self.scm_external_url,
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
    def from_env(cls) -> "GitOpsConfig":
        return cls(
            api_base_url=read_required_env("SCM_API_URL"),
            token=read_required_env("SCM_TOKEN"),
            gitops_owner=read_required_env("GITOPS_REPO_OWNER"),
            gitops_repo=read_required_env("GITOPS_REPO_NAME"),
            branch=read_required_env("GITOPS_BRANCH"),
            base_path=read_required_env("GITOPS_BASE_PATH"),
            author_name=read_required_env("GITOPS_COMMIT_USER_NAME"),
            author_email=read_required_env("GITOPS_COMMIT_USER_EMAIL"),
        )
