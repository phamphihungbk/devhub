from dataclasses import dataclass
from typing import Any

from scaffolders import read_required_str


@dataclass(frozen=True)
class ScaffoldResponse:
    repo_url: str
    path: str

    @classmethod
    def from_dict(cls, payload: dict[str, Any]):
        return cls(
            repo_url=read_required_str(payload, "repo_url"),
            path=read_required_str(payload, "path"),
        )

    def to_dict(self):
        return {
            "repo_url": self.repo_url, 
            "path": self.path
        }
