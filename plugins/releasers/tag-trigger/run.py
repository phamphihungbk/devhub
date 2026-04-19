import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[2]))
from releasers import fail, success  # noqa: E402
from releasers.models.payload import ReleasePayload
from releasers.services import TagReleaseService
from releasers import read_payload
import json

SCHEMA_PATH = Path(__file__).with_name("schema.json")


def run() -> None:
    if not SCHEMA_PATH.exists():
        fail("schema.json is not found")

    schema = json.loads(SCHEMA_PATH.read_text(encoding="utf-8"))

    payload_dict = read_payload(required_fields=schema.get("required", ["tag"]))
    payload = ReleasePayload.from_dict(payload_dict)

    service = TagReleaseService()
    response = service.trigger(payload)

    success(response.to_dict())


if __name__ == "__main__":
    run()
