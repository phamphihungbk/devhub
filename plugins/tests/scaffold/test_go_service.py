import json
import subprocess
from pathlib import Path


def test_go_service_scaffold(tmp_path):
    input_payload = {
        "action": "scaffold.go_service",
        "correlation_id": "test-123",
        "payload": {
            "service_name": "test-service",
            "owner": "team-a",
            "output_dir": str(tmp_path),
        },
    }

    result = subprocess.run(
        ["python", "scaffold/go_service/action.py"],
        input=json.dumps(input_payload),
        capture_output=True,
        text=True,
    )

    assert result.returncode == 0

    service_dir = tmp_path / "test-service"
    assert service_dir.exists()
    assert (service_dir / "main.go").exists()