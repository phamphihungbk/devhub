import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[2]))

from scaffolders.models.payload import GitOpsConfig, ScaffoldPayload
from scaffolders.models.response import ScaffoldResponse
from scaffolders import fail, success
from scaffolders.clients.scm_repository_client import SCMRepositoryClient
from scaffolders.clients.git_repository_publisher import GitRepositoryPublisher
from scaffolders.clients.gitops_values_publisher import GitOpsValuesPublisher
from scaffolders.services.scaffold_service_generator import ScaffoldServiceGenerator
from scaffolders.services.scaffold_bootstrapper import ScaffoldBootstrapper

SCHEMA_PATH = Path(__file__).with_name("schema.json")
LOCAL_TEMPLATE_DIR = Path(__file__).with_name("template")
VALUES_TEMPLATE_PATH = LOCAL_TEMPLATE_DIR / "deploy" / "helm" / "values.yaml"


def run() -> None:
    config = GitOpsConfig.from_env()
    if not config:
        fail("gitops config is not set")

    generator = ScaffoldServiceGenerator(
        schema_path=SCHEMA_PATH,
        template_dir=LOCAL_TEMPLATE_DIR,
        values_template_path=VALUES_TEMPLATE_PATH,
    )

    schema = generator.load_schema()
    payload_dict = generator.parse_payload(schema)
    payload = ScaffoldPayload.from_dict(payload_dict)

    scm_client = SCMRepositoryClient(config)
    git_publisher = GitRepositoryPublisher(
        branch=config.branch,
        author_name=config.author_name,
        author_email=config.author_email,
    )
    values_publisher = GitOpsValuesPublisher(scm_client)
    bootstrapper = ScaffoldBootstrapper(
        scm_client=scm_client,
        git_publisher=git_publisher,
        gitops_values_publisher=values_publisher,
    )

    temp_dir, service_dir = generator.create_source(payload)
    try:
        values_content = generator.build_gitops_values_content(payload)
        bootstrapper.bootstrap(payload, service_dir, values_content)

        response = ScaffoldResponse.from_dict(
            {
                "repo_url": payload.repo_url,
                "path": str(service_dir),
            }
        )
        success(response.to_dict())
    finally:
        temp_dir.cleanup()


if __name__ == "__main__":
    run()
