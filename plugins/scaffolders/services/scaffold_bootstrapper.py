from scaffolders.models.payload import ScaffoldPayload
from scaffolders import fail


class ScaffoldBootstrapper:
    def __init__(self, scm_client, git_publisher, gitops_values_publisher):
        self.scm_client = scm_client
        self.git_publisher = git_publisher
        self.gitops_values_publisher = gitops_values_publisher

    def bootstrap(self, payload: ScaffoldPayload, service_dir, values_content: str) -> None:
        owner, repo_name = self.scm_client.parse_repo_coordinates(payload.repo_url)

        if self.scm_client.repo_exists(owner, repo_name):
            fail(f"repository {owner}/{repo_name} already exists")

        self.scm_client.create_repo(
            owner=owner,
            repo_name=repo_name,
            description=f"Generated service for {payload.service_name}",
        )

        remote_url = self.scm_client.build_authenticated_remote_url(owner, repo_name)
        self.git_publisher.publish_new_repository(service_dir, remote_url)

        self.gitops_values_publisher.bootstrap(
            environment=payload.environment,
            service_name=payload.service_name,
            values_content=values_content,
            base_path=self.scm_client.config.base_path,
        )
