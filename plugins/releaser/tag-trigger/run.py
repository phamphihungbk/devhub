import json
from datetime import datetime, timezone
from pathlib import Path
import sys
from typing import Any

sys.path.insert(0, str(Path(__file__).resolve().parents[2]))
from releaser import (  # noqa: E402
    read_optional_str,
    read_payload,
    read_required_str,
    push_tag_to_repo,
    success,
)

SCHEMA_PATH = Path(__file__).with_name("schema.json")


def load_schema() -> dict[str, Any]:
    return json.loads(SCHEMA_PATH.read_text(encoding="utf-8"))


def parse_payload(schema: dict[str, Any]) -> dict[str, Any]:
    required_fields = schema.get("required", ["tag"])
    return read_payload(required_fields=required_fields)


def run() -> None:
    schema = load_schema()
    payload = parse_payload(schema)

    tag = read_required_str(payload, "tag")
    repo_url = read_required_str(payload, "repo_url")
    target = read_optional_str(payload, "target", "main")
    name = read_optional_str(payload, "name", tag)
    notes = read_optional_str(payload, "notes")

    push_tag_to_repo(
        repo_url=repo_url,
        tag=tag,
        target=target,
        message=notes or f"Release {name}",
    )

    success(
        {
            "external_ref": read_optional_str(payload, "release_id", tag),
            "commit_sha": target,
            "tag": tag,
            "name": name,
            "notes": notes,
            "repo_url": repo_url,
            "finished_at": datetime.now(timezone.utc).isoformat(),
        }
    )


if __name__ == "__main__":
    run()
