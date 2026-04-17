import json
import os
import shutil
import subprocess
import tempfile
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path
from typing import Any

from .io import fail, read_required_str


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

    push_directory_to_remote(service_dir, auth_remote, default_branch, git_author_name, git_author_email)


def push_service_to_repo(service_dir: Path, repo_url: str) -> None:
    default_branch = os.getenv("GITEA_DEFAULT_BRANCH", "main").strip() or "main"
    git_author_name = os.getenv("GIT_AUTHOR_NAME", "DevHub Scaffold Bot").strip() or "DevHub Scaffold Bot"
    git_author_email = os.getenv("GIT_AUTHOR_EMAIL", "devhub@example.local").strip() or "devhub@example.local"
    remote_url = build_push_remote_url(repo_url)

    push_directory_to_remote(service_dir, remote_url, default_branch, git_author_name, git_author_email)


def push_tag_to_repo(repo_url: str, tag: str, target: str, message: str = "") -> None:
    remote_url = build_push_remote_url(repo_url)
    git_author_name = os.getenv("GIT_AUTHOR_NAME", "DevHub Release Bot").strip() or "DevHub Release Bot"
    git_author_email = os.getenv("GIT_AUTHOR_EMAIL", "devhub@example.local").strip() or "devhub@example.local"

    with tempfile.TemporaryDirectory(prefix="tag-trigger-") as temp_dir:
        repo_dir = Path(temp_dir)
        run_git(repo_dir, "clone", remote_url, ".")
        run_git(repo_dir, "config", "user.name", git_author_name)
        run_git(repo_dir, "config", "user.email", git_author_email)
        run_git(repo_dir, "fetch", "--tags", "origin")

        target_ref = target.strip() or os.getenv("GITEA_DEFAULT_BRANCH", "main").strip() or "main"
        run_git(repo_dir, "rev-parse", "--verify", target_ref)

        if message.strip() != "":
            run_git(repo_dir, "tag", "-a", tag, target_ref, "-m", message.strip())
        else:
            run_git(repo_dir, "tag", tag, target_ref)

        run_git(repo_dir, "push", "origin", f"refs/tags/{tag}")


def push_directory_to_remote(
    source_dir: Path,
    remote_url: str,
    default_branch: str,
    git_author_name: str,
    git_author_email: str,
) -> None:
    with tempfile.TemporaryDirectory(prefix="repo-push-") as temp_dir:
        temp_path = Path(temp_dir)
        repo_dir = temp_path / "repo"

        run_git(temp_path, "clone", remote_url, repo_dir.name)
        run_git(repo_dir, "config", "user.name", git_author_name)
        run_git(repo_dir, "config", "user.email", git_author_email)
        run_git(repo_dir, "checkout", "-B", default_branch, check=False)

        sync_directory_contents(source_dir, repo_dir)

        has_changes = stage_and_detect_changes(repo_dir)
        if has_changes:
            run_git(repo_dir, "commit", "-m", "Initial scaffold from DevHub")
            run_git(repo_dir, "push", "-u", "origin", default_branch)


def sync_directory_contents(source_dir: Path, destination_dir: Path) -> None:
    clear_directory_contents(destination_dir, preserve={".git"})

    for path in source_dir.iterdir():
        destination = destination_dir / path.name
        if path.is_dir():
            shutil.copytree(path, destination, dirs_exist_ok=True)
        else:
            shutil.copy2(path, destination)


def clear_directory_contents(directory: Path, preserve: set[str] | None = None) -> None:
    protected = preserve or set()
    for path in directory.iterdir():
        if path.name in protected:
            continue
        if path.is_dir():
            shutil.rmtree(path)
        else:
            path.unlink()


def initialize_repo(service_dir: Path, git_author_name: str, git_author_email: str, default_branch: str) -> None:
    run_git(service_dir, "init")
    run_git(service_dir, "checkout", "-B", default_branch)
    run_git(service_dir, "config", "user.name", git_author_name)
    run_git(service_dir, "config", "user.email", git_author_email)
    has_changes = stage_and_detect_changes(service_dir)
    if has_changes:
        run_git(service_dir, "commit", "-m", "Initial scaffold from DevHub")


def build_authenticated_remote_url(api_base_url: str, owner: str, repo_name: str, username: str, token: str) -> str:
    parsed = urllib.parse.urlparse(api_base_url)
    safe_username = urllib.parse.quote(username, safe="")
    safe_token = urllib.parse.quote(token, safe="")
    netloc = f"{safe_username}:{safe_token}@{parsed.netloc}"
    return urllib.parse.urlunparse((parsed.scheme, netloc, f"/{owner}/{repo_name}.git", "", "", ""))


def build_push_remote_url(repo_url: str) -> str:
    repo_url = rewrite_repo_url_for_container_access(repo_url)
    parsed = urllib.parse.urlparse(repo_url)

    if parsed.scheme not in {"http", "https"}:
        return repo_url

    if parsed.username or parsed.password:
        return repo_url

    username = resolve_git_username(repo_url)
    token = first_non_empty(
        os.getenv("GITEA_TOKEN", "").strip(),
        os.getenv("SCM_TOKEN", "").strip(),
    )
    if username == "" or token == "":
        return repo_url

    safe_username = urllib.parse.quote(username, safe="")
    safe_token = urllib.parse.quote(token, safe="")
    netloc = f"{safe_username}:{safe_token}@{parsed.netloc}"
    return urllib.parse.urlunparse(
        (parsed.scheme, netloc, parsed.path, parsed.params, parsed.query, parsed.fragment)
    )


def rewrite_repo_url_for_container_access(repo_url: str) -> str:
    repo = urllib.parse.urlparse(repo_url)
    if repo.scheme not in {"http", "https"} or repo.netloc == "":
        return repo_url

    internal_base_url = first_non_empty(
        os.getenv("GITEA_URL", "").strip().rstrip("/"),
        "http://gitea:3000",
        derive_base_url_from_api(os.getenv("SCM_API_URL", "").strip()),
    )
    internal = urllib.parse.urlparse(internal_base_url)
    if not internal.netloc:
        return repo_url

    external_candidates = collect_external_repo_hosts()
    if repo.netloc not in external_candidates:
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


def derive_base_url_from_api(api_url: str) -> str:
    api_url = api_url.strip()
    if api_url == "":
        return ""

    parsed = urllib.parse.urlparse(api_url)
    if not parsed.scheme or not parsed.netloc:
        return ""

    return urllib.parse.urlunparse((parsed.scheme, parsed.netloc, "", "", "", ""))


def collect_external_repo_hosts() -> set[str]:
    candidates = set()

    for raw in (
        os.getenv("GITEA_EXTERNAL_URL", "").strip().rstrip("/"),
        os.getenv("SCM_EXTERNAL_URL", "").strip().rstrip("/"),
        "https://gitea.devhub.local",
        "http://gitea.devhub.local",
        "http://localhost:3000",
        "https://localhost:3000",
    ):
        if raw == "":
            continue
        parsed = urllib.parse.urlparse(raw)
        if parsed.netloc:
            candidates.add(parsed.netloc)

    gitea_domain = os.getenv("GITEA_DOMAIN", "").strip()
    if gitea_domain != "":
        candidates.add(gitea_domain)

    return candidates


def resolve_git_username(repo_url: str) -> str:
    explicit = first_non_empty(
        os.getenv("GITEA_USERNAME", "").strip(),
        os.getenv("SCM_USERNAME", "").strip(),
    )
    if explicit:
        return explicit

    parsed = urllib.parse.urlparse(repo_url)
    path_segments = [segment for segment in parsed.path.split("/") if segment]
    if len(path_segments) >= 2:
        return path_segments[0]

    return ""


def first_non_empty(*values: str) -> str:
    for value in values:
        if value:
            return value
    return ""


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
