import json
import tempfile
from pathlib import Path
from typing import Any

from scaffolders.models.payload import ScaffoldPayload
from scaffolders import (
    normalize_module_path,
    read_payload,
    render_template,
    resolve_service_dir,
    scaffold_from_directory,
    split_container_image,
)


class ScaffoldServiceGenerator:
    def __init__(self, schema_path: Path, template_dir: Path, values_template_path: Path):
        self.schema_path = schema_path
        self.template_dir = template_dir
        self.values_template_path = values_template_path

    def load_schema(self) -> dict[str, Any]:
        return json.loads(self.schema_path.read_text(encoding="utf-8"))

    def parse_payload(self, schema: dict[str, Any]) -> dict[str, Any]:
        return read_payload(required_fields=schema.get("required", ["service_name"]))

    def create_source(self, payload: ScaffoldPayload) -> tuple[tempfile.TemporaryDirectory, Path]:
        temp_dir = tempfile.TemporaryDirectory(prefix=f"{payload.service_name}-")
        service_dir = resolve_service_dir(temp_dir.name, payload.service_name)

        scaffold_from_directory(
            service_dir,
            self.template_dir,
            self.build_template_context(payload),
        )
        self.remove_repo_values_file(service_dir)
        return temp_dir, service_dir

    def build_gitops_values_content(self, payload: ScaffoldPayload) -> str:
        template = self.values_template_path.read_text(encoding="utf-8")
        return render_template(template, self.build_template_context(payload)).strip()

    def build_template_context(self, payload: ScaffoldPayload) -> dict[str, str]:
        image_repository, image_tag = split_container_image(payload.image)

        return {
            "SERVICE_NAME": payload.service_name,
            "MODULE_PATH": normalize_module_path(payload.module_path, payload.service_name),
            "PORT": str(payload.port),
            "IMAGE": payload.image,
            "IMAGE_REPOSITORY": image_repository,
            "IMAGE_TAG": image_tag,
            "REPO_URL": payload.repo_url,
            "TARGET_REVISION": payload.target_revision,
            "NAMESPACE": payload.namespace,
            "ARGOCD_PROJECT": payload.argocd_project,
            "REGISTRY_URL": payload.registry_url,
            "SERVER_URL": payload.server_url,
            "ENVIRONMENT": payload.environment,
        }

    def remove_repo_values_file(self, service_dir: Path) -> None:
        repo_values_path = service_dir / "deploy" / "helm" / "values.yaml"
        if repo_values_path.exists():
            repo_values_path.unlink()
