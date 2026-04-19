from dataclasses import dataclass
from typing import Any

from releasers import read_optional_str, read_required_str


@dataclass(frozen=True)
class ReleasePayload:
    release_id: str
    service_id: str
    plugin_id: str
    tag: str
    repo_url: str
    target: str
    name: str
    notes: str

    @classmethod
    def from_dict(cls, payload: dict[str, Any]):
        tag = read_required_str(payload, "tag")

        return cls(
            release_id=read_optional_str(payload, "release_id", tag),
            service_id=read_optional_str(payload, "service_id"),
            plugin_id=read_optional_str(payload, "plugin_id"),
            tag=tag,
            repo_url=read_required_str(payload, "repo_url"),
            target=read_optional_str(payload, "target", "main"),
            name=read_optional_str(payload, "name", tag),
            notes=read_optional_str(payload, "notes"),
        )
