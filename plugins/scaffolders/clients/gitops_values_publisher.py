class GitOpsValuesPublisher:
    def __init__(self, scm_client):
        self.scm_client = scm_client

    def bootstrap(self, environment: str, service_name: str, values_content: str, base_path: str) -> None:
        self.scm_client.save_file(
            path=f"{base_path}/{environment}/{service_name}.yaml",
            content=values_content,
            message=f"bootstrap gitops values for {service_name}",
        )