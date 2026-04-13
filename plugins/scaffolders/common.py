import json
import os
import re
import shutil
import subprocess
import sys
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path
from typing import Any


def fail(reason: str) -> None:
    print(json.dumps({"status": "error", "error": {"reason": reason}}))
    sys.exit(1)


def success(output: dict[str, Any]) -> None:
    print(json.dumps({"status": "ok", "output": output}))
    sys.exit(0)


def read_payload(required_fields: list[str]) -> dict[str, Any]:
    raw = sys.stdin.read()
    if not raw.strip():
        fail("empty stdin payload")

    try:
        envelope = json.loads(raw)
    except json.JSONDecodeError as exc:
        fail(f"invalid JSON input: {exc}")

    payload = envelope.get("payload")
    if not isinstance(payload, dict):
        fail("payload is required and must be an object")

    for key in required_fields:
        if key not in payload:
            fail(f"missing required field: payload.{key}")

    return payload


def read_required_str(payload: dict[str, Any], key: str) -> str:
    value = str(payload.get(key, "")).strip()
    if value == "":
        fail(f"{key} must not be empty")
    return value


def read_int(payload: dict[str, Any], key: str, default: int, min_value: int, max_value: int) -> int:
    raw = payload.get(key, default)
    try:
        value = int(raw)
    except (TypeError, ValueError):
        fail(f"{key} must be an integer")

    if value < min_value or value > max_value:
        fail(f"{key} must be between {min_value} and {max_value}")
    return value


def normalize_module_path(module_path: str, service_name: str) -> str:
    base = module_path.strip().rstrip("/")
    if base == "":
        return service_name
    return f"{base}/{service_name}"


def validate_service_name(name: str) -> None:
    if re.match(r"^[a-z0-9][a-z0-9-]*$", name) is None:
        fail("service_name must match ^[a-z0-9][a-z0-9-]*$")


def resolve_service_dir(output_dir_raw: str, service_name: str) -> Path:
    output_dir = Path(output_dir_raw).expanduser()
    service_dir = output_dir / service_name
    service_dir.mkdir(parents=True, exist_ok=True)
    return service_dir


def scaffold_from_template(service_dir: Path, template_name: str, replacements: dict[str, str]) -> None:
    templates_root = Path(__file__).resolve().parents[2] / "templates"
    template_dir = templates_root / template_name
    common_chart_dir = templates_root / "charts" / "app"

    if not template_dir.is_dir():
        fail(f"template directory does not exist: {template_dir}")

    if service_dir.exists():
        shutil.rmtree(service_dir)
    service_dir.mkdir(parents=True, exist_ok=True)

    copy_template_tree(template_dir, service_dir, replacements)

    if common_chart_dir.is_dir():
        chart_target_dir = service_dir / "charts" / "app"
        chart_target_dir.mkdir(parents=True, exist_ok=True)
        copy_template_tree(common_chart_dir, chart_target_dir, replacements)


def copy_template_tree(source_dir: Path, target_dir: Path, replacements: dict[str, str]) -> None:
    for path in source_dir.rglob("*"):
        relative_path = path.relative_to(source_dir)
        destination = target_dir / relative_path

        if path.is_dir():
            destination.mkdir(parents=True, exist_ok=True)
            continue

        destination.parent.mkdir(parents=True, exist_ok=True)
        copy_template_file(path, destination, replacements)


def copy_template_file(source: Path, destination: Path, replacements: dict[str, str]) -> None:
    try:
        raw = source.read_text(encoding="utf-8")
    except UnicodeDecodeError:
        shutil.copy2(source, destination)
        return

    destination.write_text(render_template(raw, replacements), encoding="utf-8")


def render_template(content: str, replacements: dict[str, str]) -> str:
    rendered = content
    for key, value in replacements.items():
        rendered = rendered.replace(f"{{{{{key}}}}}", value)
    return rendered


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


def build_scaffold_output(service_dir: Path, service_name: str, payload: dict[str, Any]) -> dict[str, Any]:
    repo_url = read_required_str(payload, "repo_url")
    push_service_to_repo(service_dir, repo_url)
    return {"repo_url": repo_url, "path": str(service_dir)}


def maybe_publish_to_gitea(service_dir: Path, service_name: str, payload: dict[str, Any]) -> str:
    username = os.getenv("GITEA_USERNAME", "").strip()
    token = os.getenv("GITEA_TOKEN", "").strip()

    if username == "" or token == "":
        return ""

    api_base_url = os.getenv("GITEA_URL", "http://gitea:3000").strip().rstrip("/")
    external_base_url = os.getenv("GITEA_EXTERNAL_URL", "https://gitea.devhub.local").strip().rstrip("/")
    owner = os.getenv("GITEA_OWNER", "").strip() or username
    default_branch = os.getenv("GITEA_DEFAULT_BRANCH", "main").strip() or "main"
    is_private = os.getenv("GITEA_PRIVATE", "false").strip().lower() in {"1", "true", "yes", "on"}
    description = str(payload.get("description", "")).strip()

    create_gitea_repo(
        api_base_url=api_base_url,
        owner=owner,
        repo_name=service_name,
        token=token,
        username=username,
        description=description,
        default_branch=default_branch,
        is_private=is_private,
    )
    push_service_to_gitea(
        service_dir=service_dir,
        api_base_url=api_base_url,
        owner=owner,
        repo_name=service_name,
        username=username,
        token=token,
        default_branch=default_branch,
    )

    return f"{external_base_url}/{owner}/{service_name}.git"


def create_gitea_repo(
    api_base_url: str,
    owner: str,
    repo_name: str,
    token: str,
    username: str,
    description: str,
    default_branch: str,
    is_private: bool,
) -> None:
    if owner == username:
        path = "/api/v1/user/repos"
    else:
        path = f"/api/v1/orgs/{urllib.parse.quote(owner)}/repos"

    body = json.dumps(
        {
            "name": repo_name,
            "description": description,
            "default_branch": default_branch,
            "private": is_private,
            "auto_init": False,
        }
    ).encode("utf-8")
    request = urllib.request.Request(
        f"{api_base_url}{path}",
        data=body,
        headers={
            "Authorization": f"token {token}",
            "Content-Type": "application/json",
            "Accept": "application/json",
        },
        method="POST",
    )

    try:
        with urllib.request.urlopen(request) as response:
            if response.status not in (200, 201):
                fail(f"failed to create Gitea repo {owner}/{repo_name}: HTTP {response.status}")
    except urllib.error.HTTPError as exc:
        try:
            error_body = exc.read().decode("utf-8", errors="replace")
        except Exception:
            error_body = str(exc)
        fail(f"failed to create Gitea repo {owner}/{repo_name}: HTTP {exc.code} {error_body}")
    except urllib.error.URLError as exc:
        fail(f"failed to reach Gitea API at {api_base_url}: {exc}")


def push_service_to_gitea(
    service_dir: Path,
    api_base_url: str,
    owner: str,
    repo_name: str,
    username: str,
    token: str,
    default_branch: str,
) -> None:
    auth_remote = build_authenticated_remote_url(api_base_url, owner, repo_name, username, token)
    git_author_name = os.getenv("GIT_AUTHOR_NAME", "DevHub Scaffold Bot").strip() or "DevHub Scaffold Bot"
    git_author_email = os.getenv("GIT_AUTHOR_EMAIL", "devhub@example.local").strip() or "devhub@example.local"

    run_git(service_dir, "init")
    run_git(service_dir, "checkout", "-B", default_branch)
    run_git(service_dir, "config", "user.name", git_author_name)
    run_git(service_dir, "config", "user.email", git_author_email)
    has_changes = stage_and_detect_changes(service_dir)
    if has_changes:
        run_git(service_dir, "commit", "-m", "Initial scaffold from DevHub")
    run_git(service_dir, "remote", "remove", "origin", check=False)
    run_git(service_dir, "remote", "add", "origin", auth_remote)
    run_git(service_dir, "push", "-u", "origin", default_branch)


def push_service_to_repo(service_dir: Path, repo_url: str) -> None:
    default_branch = os.getenv("GITEA_DEFAULT_BRANCH", "main").strip() or "main"
    git_author_name = os.getenv("GIT_AUTHOR_NAME", "DevHub Scaffold Bot").strip() or "DevHub Scaffold Bot"
    git_author_email = os.getenv("GIT_AUTHOR_EMAIL", "devhub@example.local").strip() or "devhub@example.local"
    remote_url = build_push_remote_url(repo_url)

    run_git(service_dir, "init")
    run_git(service_dir, "checkout", "-B", default_branch)
    run_git(service_dir, "config", "user.name", git_author_name)
    run_git(service_dir, "config", "user.email", git_author_email)
    has_changes = stage_and_detect_changes(service_dir)
    if has_changes:
        run_git(service_dir, "commit", "-m", "Initial scaffold from DevHub")
    run_git(service_dir, "remote", "remove", "origin", check=False)
    run_git(service_dir, "remote", "add", "origin", remote_url)
    run_git(service_dir, "push", "-u", "origin", default_branch)


def build_authenticated_remote_url(api_base_url: str, owner: str, repo_name: str, username: str, token: str) -> str:
    parsed = urllib.parse.urlparse(api_base_url)
    safe_username = urllib.parse.quote(username, safe="")
    safe_token = urllib.parse.quote(token, safe="")
    netloc = f"{safe_username}:{safe_token}@{parsed.netloc}"
    return urllib.parse.urlunparse(
        (parsed.scheme, netloc, f"/{owner}/{repo_name}.git", "", "", "")
    )


def build_push_remote_url(repo_url: str) -> str:
    repo_url = rewrite_repo_url_for_container_access(repo_url)
    parsed = urllib.parse.urlparse(repo_url)

    if parsed.scheme not in {"http", "https"}:
        return repo_url

    if parsed.username or parsed.password:
        return repo_url

    username = os.getenv("GITEA_USERNAME", "").strip()
    token = os.getenv("GITEA_TOKEN", "").strip()
    if username == "" or token == "":
        return repo_url

    safe_username = urllib.parse.quote(username, safe="")
    safe_token = urllib.parse.quote(token, safe="")
    netloc = f"{safe_username}:{safe_token}@{parsed.netloc}"
    return urllib.parse.urlunparse(
        (parsed.scheme, netloc, parsed.path, parsed.params, parsed.query, parsed.fragment)
    )


def rewrite_repo_url_for_container_access(repo_url: str) -> str:
    external_base_url = os.getenv("GITEA_EXTERNAL_URL", "").strip().rstrip("/")
    internal_base_url = os.getenv("GITEA_URL", "").strip().rstrip("/")

    if external_base_url == "" or internal_base_url == "":
        return repo_url

    external = urllib.parse.urlparse(external_base_url)
    internal = urllib.parse.urlparse(internal_base_url)
    repo = urllib.parse.urlparse(repo_url)

    if not external.netloc or not internal.netloc:
        return repo_url

    if repo.netloc != external.netloc:
        return repo_url

    return urllib.parse.urlunparse(
        (
            internal.scheme or repo.scheme,
            internal.netloc,
            repo.path,
            repo.params,
            repo.query,
            repo.fragment,
        )
    )


def run_git(service_dir: Path, *args: str, check: bool = True) -> None:
    completed = subprocess.run(
        ["git", *args],
        cwd=service_dir,
        check=False,
        capture_output=True,
        text=True,
    )
    if check and completed.returncode != 0:
        fail(
            f"git {' '.join(args)} failed: {completed.stderr.strip() or completed.stdout.strip() or f'exit {completed.returncode}'}"
        )


def stage_and_detect_changes(service_dir: Path) -> bool:
    run_git(service_dir, "add", ".")
    completed = subprocess.run(
        ["git", "diff", "--cached", "--quiet"],
        cwd=service_dir,
        check=False,
        capture_output=True,
        text=True,
    )
    if completed.returncode == 0:
        return False
    if completed.returncode == 1:
        return True
    fail(
        "git diff --cached --quiet failed: "
        + (completed.stderr.strip() or completed.stdout.strip() or f"exit {completed.returncode}")
    )
    return False
