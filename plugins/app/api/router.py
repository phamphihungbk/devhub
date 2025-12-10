from fastapi.routing import APIRouter

from devhub_plugin.web.api import echo, kafka, monitoring

api_router = APIRouter()
api_router.include_router(monitoring.router)
api_router.include_router(echo.router, prefix="/echo", tags=["echo"])
api_router.include_router(kafka.router, prefix="/kafka", tags=["kafka"])
