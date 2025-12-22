import json
import sys
from dataclasses import dataclass
from typing import Any, Dict


@dataclass
class ActionResult:
    status: str
    output: Dict[str, Any] | None = None
    error: Dict[str, Any] | None = None

    def ok(self) -> None:
        print(json.dumps({"status": "ok", "output": self.output or {}}))
        sys.exit(0)

    def fail(self) -> None:
        print(json.dumps({"status": "error", "error": self.error or {}}))
        sys.exit(1)