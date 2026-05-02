import base64
import json
import urllib.error
import urllib.parse
import urllib.request
from typing import Any

from scaffolders import fail
from scaffolders.models.payload import GitOpsConfig


class SCMRepositoryClient:
    def __init__(self, config: GitOpsConfig):
        self.config = config

    def repo_exists(self, owner: str, repo_name: str) -> bool:
        try:
            self._request("GET", f"/repos/{urllib.parse.quote(owner)}/{urllib.parse.quote(repo_name)}")
            return True
        except RuntimeError as exc:
            if "HTTP 404" in str(exc):
                return False
            raise

    def create_repo(self, owner: str, repo_name: str, description: str) -> None:
        body = {
            "name": repo_name,
            "description": description,
            "default_branch": self.config.branch,
            "private": False,
            "auto_init": False,
        }

        try:
            self._request("POST", f"/orgs/{urllib.parse.quote(owner)}/repos", body)
            return
        except RuntimeError as exc:
            if "HTTP 404" not in str(exc):
                raise

        self._request("POST", "/user/repos", {**body, "organization": owner})

    def get_file(self, path: str) -> dict[str, Any] | None:
        try:
            _, body = self._contents_request("GET", f"{path}?ref={self.config.branch}")
        except RuntimeError as exc:
            if "404" in str(exc):
                return None
            raise

        data = json.loads(body)
        if isinstance(data, dict) and data.get("sha"):
            return data
        return None

    def create_file(self, path: str, content: str, message: str) -> None:
        self._contents_request("POST", path, self._build_file_body(content, message))

    def update_file(self, path: str, content: str, sha: str, message: str) -> None:
        body = self._build_file_body(content, message)
        body["sha"] = sha
        self._contents_request("PUT", path, body)

    def save_file(self, path: str, content: str, message: str) -> None:
        try:
            self.create_file(path, content, message)
            return
        except RuntimeError as exc:
            if "[SHA]" not in str(exc):
                raise

        existing = self.get_file(path)
        if not existing or "sha" not in existing:
            raise RuntimeError(f"cannot resolve sha for {path}")

        self.update_file(path, content, existing["sha"], message)

    def parse_repo_coordinates(self, repo_url: str) -> tuple[str, str]:
        parsed = urllib.parse.urlparse(repo_url.strip())
        path_segments = [segment for segment in parsed.path.split("/") if segment]

        if len(path_segments) < 2:
            fail(f"invalid repo_url: {repo_url}")

        owner = path_segments[-2]
        repo_name = path_segments[-1]
        if repo_name.endswith(".git"):
            repo_name = repo_name[:-4]

        if not owner or not repo_name:
            fail(f"invalid repo_url: {repo_url}")

        return owner, repo_name

    def build_authenticated_remote_url(self, owner: str, repo_name: str) -> str:
        username = self.config.gitops_owner
        if not owner:
            fail(f"invalid owner: {owner}")

        if not repo_name:
            fail("repo_name is required")

        if not username:
            fail("SCM_USERNAME (username) is required")

        if not self.config.token:
            fail("SCM_TOKEN is required")

        parsed = urllib.parse.urlparse(self.config.api_base_url)
        safe_username = urllib.parse.quote(username, safe="")
        safe_token = urllib.parse.quote(self.config.token, safe="")
        netloc = f"{safe_username}:{safe_token}@{parsed.netloc}"

        return urllib.parse.urlunparse(
            (parsed.scheme, netloc, f"/{owner}/{repo_name}.git", "", "", "")
        )

    def _build_file_body(self, content: str, message: str) -> dict[str, Any]:
        return {
            "branch": self.config.branch,
            "content": base64.b64encode(content.encode("utf-8")).decode("utf-8"),
            "message": message,
            "author": {
                "name": self.config.author_name,
                "email": self.config.author_email,
            },
            "committer": {
                "name": self.config.author_name,
                "email": self.config.author_email,
            },
        }

    def _contents_request(
        self,
        method: str,
        path: str,
        body: dict[str, Any] | None = None,
    ) -> tuple[int, str]:
        return self._request(
            method,
            f"/repos/{urllib.parse.quote(self.config.gitops_owner)}/{urllib.parse.quote(self.config.gitops_repo)}/contents/{path}",
            body,
        )

    def _request(self, method: str, path: str, body: dict[str, Any] | None = None) -> tuple[int, str]:
        normalized_path = path if path.startswith("/") else f"/{path}"
        url = f"{self.config.api_base_url}{normalized_path}"

        headers = {"Accept": "application/json"}
        if self.config.token:
            headers["Authorization"] = f"token {self.config.token}"

        data = None
        if body is not None:
            data = json.dumps(body).encode("utf-8")
            headers["Content-Type"] = "application/json"

        req = urllib.request.Request(url, data=data, headers=headers, method=method)

        try:
            with urllib.request.urlopen(req) as resp:
                return resp.getcode(), resp.read().decode("utf-8")
        except urllib.error.HTTPError as exc:
            error_body = exc.read().decode("utf-8")
            raise RuntimeError(f"HTTP {exc.code}: {error_body}") from exc
