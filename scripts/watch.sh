#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH='' cd -- "$(dirname "$0")/.." && pwd)"

export COMPOSE_ENV=dev

BACKEND_SERVICES="${BACKEND_SERVICES:-backend worker db redis devhub-registry gitea gitea-runner tempo grafana otel-collector}"
FRONTEND_SERVICES="${FRONTEND_SERVICES:-frontend nginx}"
WATCH_LOG_SERVICES="${WATCH_LOG_SERVICES:-backend worker frontend nginx}"

"${ROOT_DIR}/scripts/docker-build-and-run.sh" up -d --build ${BACKEND_SERVICES} ${FRONTEND_SERVICES}
exec "${ROOT_DIR}/scripts/docker-build-and-run.sh" logs -f ${WATCH_LOG_SERVICES}
