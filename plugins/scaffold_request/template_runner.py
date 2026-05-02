import sys
import tempfile
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from scaffolders import fail, render_template, resolve_service_dir, scaffold_from_directory, success
from scaffolders.clients.git_repository_publisher import GitRepositoryPublisher
from scaffolders.clients.gitops_values_publisher import GitOpsValuesPublisher
from scaffolders.clients.scm_repository_client import SCMRepositoryClient
from scaffolders.models.payload import GitOpsConfig, ScaffoldPayload
from scaffolders.models.response import ScaffoldResponse
from scaffolders.services.scaffold_bootstrapper import ScaffoldBootstrapper
from scaffolders import read_payload


DEFAULT_REQUIRED_FIELDS = [
    "environment",
    "service_name",
    "port",
    "database",
    "image_tag",
    "module_path",
    "ci_registry_host",
    "ci_server_url",
    "cd_project_name",
    "cd_repo_url",
    "cd_target_revision",
    "cd_namespace",
    "cd_image_repository",
]


def run_template_plugin(plugin_dir: Path) -> None:
    config = GitOpsConfig.from_env()
    if not config:
        fail("gitops config is not set")

    template_dir = plugin_dir / "template"
    values_template_path = template_dir / "deploy" / "helm" / "values.yaml"
    payload = ScaffoldPayload.from_dict(read_payload(required_fields=DEFAULT_REQUIRED_FIELDS))

    temp_dir = tempfile.TemporaryDirectory(prefix=f"{payload.service_name}-")
    service_dir = resolve_service_dir(temp_dir.name, payload.service_name)

    try:
        scaffold_from_directory(service_dir, template_dir, payload.to_template())

        values_content = render_template(
            values_template_path.read_text(encoding="utf-8"),
            payload.to_template(),
        ).strip()
        repo_values_path = service_dir / "deploy" / "helm" / "values.yaml"
        if repo_values_path.exists():
            repo_values_path.unlink()

        scm_client = SCMRepositoryClient(config)
        bootstrapper = ScaffoldBootstrapper(
            scm_client=scm_client,
            git_publisher=GitRepositoryPublisher(
                branch=config.branch,
                author_name=config.author_name,
                author_email=config.author_email,
            ),
            gitops_values_publisher=GitOpsValuesPublisher(scm_client),
        )
        bootstrapper.bootstrap(payload, service_dir, values_content)

        success(ScaffoldResponse.from_dict({
            "repo_url": payload.cd_repo_url,
            "path": str(service_dir),
        }).to_dict())
    finally:
        temp_dir.cleanup()
