#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH='' cd -- "$(dirname "$0")/.." && pwd)"
ACTION="${1:-up}"

export COMPOSE_ENV=dev

exec "${ROOT_DIR}/scripts/docker-build-and-run.sh" run --rm backend go run . migrate --action "${ACTION}"
