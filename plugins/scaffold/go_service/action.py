import json
from pathlib import Path
from jsonschema import validate

from common.context import ActionContext
from common.result import ActionResult
from common.logger import get_logger


def run() -> None:
    ctx = ActionContext.from_stdin()
    logger = get_logger(ctx.correlation_id)

    schema = json.loads(Path("schema.json").read_text())
    validate(instance=ctx.payload, schema=schema)

    service_name = ctx.payload["service_name"]
    owner = ctx.payload["owner"]
    output_dir = Path(ctx.payload["output_dir"])

    service_dir = output_dir / service_name
    service_dir.mkdir(parents=True, exist_ok=True)

    (service_dir / "main.go").write_text(
        f"""package main

// Owner: {owner}

func main() {{
    println("Hello from {service_name}")
}}
"""
    )

    logger.info(f"Scaffolded Go service {service_name}")

    ActionResult(
        status="ok",
        output={"path": str(service_dir)}
    ).ok()


if __name__ == "__main__":
    run()