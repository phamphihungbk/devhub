import json
import sys
from dataclasses import dataclass
from typing import Any, Dict


@dataclass
class ActionContext:
    action: str
    correlation_id: str
    payload: Dict[str, Any]

    @staticmethod
    def from_stdin() -> "ActionContext":
        raw = sys.stdin.read()
        data = json.loads(raw)
        return ActionContext(
            action=data["action"],
            correlation_id=data["correlation_id"],
            payload=data["payload"],
        )