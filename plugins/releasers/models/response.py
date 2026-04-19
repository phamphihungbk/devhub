from dataclasses import dataclass
from datetime import datetime, timezone
from typing import Any

from releasers import read_optional_str, read_required_str


@dataclass(frozen=True)
class ReleaseResponse:
    external_ref: str
    commit_sha: str
    tag: str
    name: str
    notes: str
    repo_url: str
    finished_at: str

    @classmethod
    def from_payload(cls, payload) -> "ReleaseResponse":
        return cls(
            external_ref=payload.release_id or payload.tag,
            commit_sha=payload.target,
            tag=payload.tag,
            name=payload.name,
            notes=payload.notes,
            repo_url=payload.repo_url,
            finished_at=datetime.now(timezone.utc).isoformat(),
        )

    @classmethod
    def from_dict(cls, payload: dict[str, Any]):
        return cls(
            external_ref=read_optional_str(payload, "external_ref"),
            commit_sha=read_required_str(payload, "commit_sha"),
            tag=read_required_str(payload, "tag"),
            name=read_required_str(payload, "name"),
            notes=read_optional_str(payload, "notes"),
            repo_url=read_required_str(payload, "repo_url"),
            finished_at=read_required_str(payload, "finished_at"),
        )

    def to_dict(self):
        return {
            "external_ref": self.external_ref,
            "commit_sha": self.commit_sha,
            "tag": self.tag,
            "name": self.name,
            "notes": self.notes,
            "repo_url": self.repo_url,
            "finished_at": self.finished_at,
        }
