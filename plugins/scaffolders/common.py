import json
import re
import sys
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
