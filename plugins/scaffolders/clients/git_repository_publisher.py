import shutil
import subprocess
import tempfile
from pathlib import Path

from scaffolders import fail

INITIAL_COMMIT_MESSAGE = "Initial scaffold from DevHub"


class GitRepositoryPublisher:
    def __init__(self, branch: str, author_name: str, author_email: str):
        self.branch = branch
        self.author_name = author_name
        self.author_email = author_email

    def publish_new_repository(self, source_dir: Path, remote_url: str) -> None:
        with tempfile.TemporaryDirectory(prefix="repo-push-") as temp_dir:
            temp_path = Path(temp_dir)
            repo_dir = temp_path / "repo"

            self._run_git(temp_path, ["git", "clone", remote_url, repo_dir.name])
            self._run_git(repo_dir, ["git", "config", "user.name", self.author_name])
            self._run_git(repo_dir, ["git", "config", "user.email", self.author_email])
            self._run_git(repo_dir, ["git", "checkout", "-B", self.branch])

            self._sync_directory_contents(source_dir, repo_dir)

            if self._stage_and_detect_changes(repo_dir):
                self._run_git(repo_dir, ["git", "commit", "-m", INITIAL_COMMIT_MESSAGE])
                self._run_git(repo_dir, ["git", "push", "-u", "origin", self.branch])

    def _sync_directory_contents(self, source_dir: Path, destination_dir: Path) -> None:
        self._clear_directory_contents(destination_dir, preserve={".git"})

        for path in source_dir.iterdir():
            destination = destination_dir / path.name
            if path.is_dir():
                shutil.copytree(path, destination, dirs_exist_ok=True)
            else:
                shutil.copy2(path, destination)

    def _clear_directory_contents(self, directory: Path, preserve: set[str] | None = None) -> None:
        protected = preserve or set()
        for path in directory.iterdir():
            if path.name in protected:
                continue
            if path.is_dir():
                shutil.rmtree(path)
            else:
                path.unlink()

    def _stage_and_detect_changes(self, repo_dir: Path) -> bool:
        self._run_git(repo_dir, ["git", "add", "."])

        completed = subprocess.run(
            ["git", "diff", "--cached", "--quiet"],
            cwd=repo_dir,
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

    def _run_git(self, cwd: Path, args: list[str]) -> None:
        completed = subprocess.run(
            args,
            cwd=cwd,
            check=False,
            capture_output=True,
            text=True,
        )
        if completed.returncode != 0:
            fail(
                f"{' '.join(args)} failed: "
                + (completed.stderr.strip() or completed.stdout.strip() or f"exit {completed.returncode}")
            )