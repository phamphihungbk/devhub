import hashlib
import importlib.util
import json
import sys
from datetime import datetime, timezone
from pathlib import Path

SCHEMA_PATH = Path(__file__).with_name("schema.json")
IO_PATH = Path(__file__).resolve().parents[2] / "utils" / "io.py"

io_spec = importlib.util.spec_from_file_location("devhub_plugin_io", IO_PATH)
plugin_io = importlib.util.module_from_spec(io_spec)
io_spec.loader.exec_module(plugin_io)


def run() -> None:
    schema = json.loads(SCHEMA_PATH.read_text(encoding="utf-8"))
    payload = plugin_io.read_payload(required_fields=schema.get("required", ["tag", "repo_url"]))

    tag = plugin_io.read_required_str(payload, "tag")
    target = plugin_io.read_optional_str(payload, "target", "main")
    repo_url = plugin_io.read_required_str(payload, "repo_url")
    release_id = plugin_io.read_optional_str(payload, "release_id", tag)
    name = plugin_io.read_optional_str(payload, "name", f"Promote {tag}")
    notes = plugin_io.read_optional_str(payload, "notes", f"Promote image built from {target}")

    digest = hashlib.sha1(f"{release_id}:{repo_url}:{tag}:{target}".encode()).hexdigest()

    plugin_io.success(
        {
            "external_ref": f"image-promotion:{release_id}",
            "commit_sha": digest[:12],
            "tag": tag,
            "name": name,
            "notes": notes,
            "repo_url": repo_url,
            "finished_at": datetime.now(timezone.utc).isoformat(),
        }
    )


if __name__ == "__main__":
    run()
