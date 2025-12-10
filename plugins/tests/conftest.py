from collections.abc import AsyncGenerator
from typing import Any
from unittest.mock import Mock

import pytest
from aiokafka import AIOKafkaProducer
from fastapi import FastAPI
from httpx import ASGITransport, AsyncClient

from devhub_plugin.services.kafka.dependencies import get_kafka_producer
from devhub_plugin.services.kafka.lifespan import init_kafka, shutdown_kafka
from devhub_plugin.web.application import get_app


@pytest.fixture(scope="session")
def anyio_backend() -> str:
    """
    Backend for anyio pytest plugin.

    :return: backend name.
    """
    return "asyncio"


@pytest.fixture
async def test_kafka_producer() -> AsyncGenerator[AIOKafkaProducer]:
    """
    Creates kafka's producer.

    :yields: kafka's producer.
    """
    app_mock = Mock()
    await init_kafka(app_mock)
    yield app_mock.state.kafka_producer
    await shutdown_kafka(app_mock)


@pytest.fixture
def fastapi_app(
    test_kafka_producer: AIOKafkaProducer,
) -> FastAPI:
    """
    Fixture for creating FastAPI app.

    :return: fastapi app with mocked dependencies.
    """
    application = get_app()
    application.dependency_overrides[get_kafka_producer] = lambda: test_kafka_producer
    return application


@pytest.fixture
async def client(
    fastapi_app: FastAPI, anyio_backend: Any
) -> AsyncGenerator[AsyncClient]:
    """
    Fixture that creates client for requesting server.

    :param fastapi_app: the application.
    :yield: client for the app.
    """
    async with AsyncClient(
        transport=ASGITransport(fastapi_app), base_url="http://test", timeout=2.0
    ) as ac:
        yield ac
