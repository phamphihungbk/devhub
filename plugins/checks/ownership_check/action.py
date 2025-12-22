from common.context import ActionContext
from common.result import ActionResult


def run() -> None:
    ctx = ActionContext.from_stdin()
    owners = ctx.payload.get("owners", [])

    if len(owners) < 2:
        ActionResult(
            status="error",
            error={"reason": "At least two owners required"}
        ).fail()

    ActionResult(status="ok").ok()


if __name__ == "__main__":
    run()