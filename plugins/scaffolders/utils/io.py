import json
import sys
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


def read_optional_str(payload: dict[str, Any], key: str, default: str = "") -> str:
    value = str(payload.get(key, "")).strip()
    if value != "":
        return value
    return default


def read_int(payload: dict[str, Any], key: str, default: int, min_value: int, max_value: int) -> int:
    raw = payload.get(key, default)
    try:
        value = int(raw)
    except (TypeError, ValueError):
        fail(f"{key} must be an integer")

    if value < min_value or value > max_value:
        fail(f"{key} must be between {min_value} and {max_value}")
    return value
