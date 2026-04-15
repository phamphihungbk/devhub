import urllib.parse
import re
from pathlib import Path

from .io import fail


def normalize_module_path(module_path: str, service_name: str) -> str:
    base = module_path.strip().rstrip("/")
    if base == "":
        return service_name
    return f"{base}/{service_name}"


def infer_module_base_from_repo_url(repo_url: str) -> str:
    if repo_url == "":
        return ""

    parsed = urllib.parse.urlparse(repo_url)
    host = parsed.netloc.strip()
    path = parsed.path.strip().strip("/")
    if host == "" or path == "":
        return ""

    segments = [segment for segment in path.split("/") if segment]
    if not segments:
        return ""

    if segments[-1].endswith(".git"):
        segments[-1] = segments[-1][: -len(".git")]

    owner_segments = segments[:-1]
    if not owner_segments:
        return host

    return "/".join([host, *owner_segments])


def validate_service_name(name: str) -> None:
    if re.match(r"^[a-z0-9][a-z0-9-]*$", name) is None:
        fail("service_name must match ^[a-z0-9][a-z0-9-]*$")


def resolve_service_dir(output_dir_raw: str, service_name: str) -> Path:
    output_dir = Path(output_dir_raw).expanduser()
    service_dir = output_dir / service_name
    service_dir.mkdir(parents=True, exist_ok=True)
    return service_dir
