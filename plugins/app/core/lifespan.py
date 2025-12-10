from collections.abc import AsyncGenerator
from contextlib import asynccontextmanager

from fastapi import FastAPI

from devhub_plugin.services.kafka.lifespan import init_kafka, shutdown_kafka


@asynccontextmanager
async def lifespan_setup(
    app: FastAPI,
) -> AsyncGenerator[None]:  # pragma: no cover
    """
    Actions to run on application startup.

    This function uses fastAPI app to store data
    in the state, such as db_engine.

    :param app: the fastAPI application.
    :return: function that actually performs actions.
    """

    app.middleware_stack = None
    await init_kafka(app)
    app.middleware_stack = app.build_middleware_stack()

    yield
    await shutdown_kafka(app)
