from fastapi import FastAPI

app = FastAPI(title="[[SERVICE_NAME]]")


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok", "service": "[[SERVICE_NAME]]"}
