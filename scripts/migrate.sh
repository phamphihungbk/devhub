#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH='' cd -- "$(dirname "$0")/.." && pwd)"
ACTION="${1:-up}"
FORCE_VERSION="${FORCE_VERSION:-}"

export COMPOSE_ENV=dev

if [ -n "${FORCE_VERSION}" ]; then
  exec "${ROOT_DIR}/scripts/docker-build-and-run.sh" run --rm backend go run . migrate --force-version "${FORCE_VERSION}"
fi

exec "${ROOT_DIR}/scripts/docker-build-and-run.sh" run --rm backend go run . migrate --action "${ACTION}"
