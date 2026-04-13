import json
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
from common import (  # noqa: E402
    build_scaffold_output,
    read_payload,
    read_required_str,
    resolve_container_image,
    resolve_service_dir,
    scaffold_from_template,
    split_container_image,
    success,
    validate_service_name,
)

SCHEMA_PATH = Path(__file__).with_name("schema.json")


def load_schema() -> dict:
    return json.loads(SCHEMA_PATH.read_text(encoding="utf-8"))


def main() -> None:
    schema = load_schema()
    payload = read_payload(required_fields=schema.get("required", ["service_name"]))

    service_name = read_required_str(payload, "service_name")
    project_id = str(payload.get("project_id", "")).strip()
    output_dir_raw = str(payload.get("output_dir", "")).strip()
    queue_name = str(payload.get("queue_name", f"{service_name}-jobs")).strip() or f"{service_name}-jobs"

    if output_dir_raw == "":
        base_dir = Path("/app/generated")
        output_dir_raw = str(base_dir / project_id) if project_id else str(base_dir)

    validate_service_name(service_name)
    service_dir = resolve_service_dir(output_dir_raw, service_name)
    image = resolve_container_image(payload, service_name)
    image_repository, image_tag = split_container_image(image)
    scaffold_from_template(
        service_dir,
        "python-worker",
        {
            "SERVICE_NAME": service_name,
            "QUEUE_NAME": queue_name,
            "IMAGE": image,
            "IMAGE_REPOSITORY": image_repository,
            "IMAGE_TAG": image_tag,
            "PORT": "8080",
        },
    )

    success(build_scaffold_output(service_dir, service_name, payload))


if __name__ == "__main__":
    main()
