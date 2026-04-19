#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH='' cd -- "$(dirname "$0")/.." && pwd)"

export COMPOSE_ENV=dev
export COMPOSE_PROFILES=ui

exec "${ROOT_DIR}/scripts/docker-build-and-run.sh" logs -f frontend nginx
