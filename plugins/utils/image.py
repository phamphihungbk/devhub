from typing import Any

from .io import fail


def resolve_container_image(payload: dict[str, Any], service_name: str) -> str:
    explicit_image = str(payload.get("image", "")).strip()
    if explicit_image != "":
        return explicit_image

    image_repository = str(payload.get("image_repository", "")).strip().rstrip("/")
    image_tag = str(payload.get("image_tag", "")).strip() or "latest"

    if image_repository == "":
        return f"{service_name}:{image_tag}"

    return f"{image_repository}/{service_name}:{image_tag}"


def split_container_image(image: str) -> tuple[str, str]:
    image = image.strip()
    if image == "":
        fail("image must not be empty")

    last_slash = image.rfind("/")
    last_colon = image.rfind(":")
    if last_colon > last_slash:
        return image[:last_colon], image[last_colon + 1 :]
    return image, "latest"
