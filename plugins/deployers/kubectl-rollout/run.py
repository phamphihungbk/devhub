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
    payload = plugin_io.read_payload(required_fields=schema.get("required", []))

    service = plugin_io.read_required_str(payload, "service")
    environment = plugin_io.read_required_str(payload, "environment")
    version = plugin_io.read_required_str(payload, "version")

    digest = hashlib.sha1(f"kubernetes:{service}:{environment}:{version}".encode()).hexdigest()

    plugin_io.success(
        {
            "external_ref": f"k8s:{service}:{environment}",
            "commit_sha": digest[:12],
            "finished_at": datetime.now(timezone.utc).isoformat(),
        }
    )


if __name__ == "__main__":
    run()
