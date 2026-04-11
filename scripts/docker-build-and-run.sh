#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH='' cd -- "$(dirname "$0")/.." && pwd)"
ENVIRONMENT="${COMPOSE_ENV:-dev}"

case "$ENVIRONMENT" in
  dev)
    COMPOSE_FILES="-f ${ROOT_DIR}/docker-compose.yml -f ${ROOT_DIR}/docker-compose.dev.yml"
    ;;
  prod)
    COMPOSE_FILES="-f ${ROOT_DIR}/docker-compose.yml -f ${ROOT_DIR}/docker-compose.prod.yml"
    ;;
  *)
    echo "Unsupported COMPOSE_ENV: ${ENVIRONMENT}" >&2
    exit 1
    ;;
esac

PROFILE_ARGS=""
if [ "${COMPOSE_FRONTEND:-0}" = "1" ]; then
  PROFILE_ARGS="--profile frontend"
fi

set -- docker compose ${COMPOSE_FILES} ${PROFILE_ARGS} "$@"
exec "$@"
