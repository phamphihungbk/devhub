#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH='' cd -- "$(dirname "$0")/.." && pwd)"
export COMPOSE_ENV=dev

if [ "${DEV_WITH_FRONTEND:-0}" = "1" ]; then
  export COMPOSE_FRONTEND=1
fi

if [ $# -eq 0 ]; then
  set -- up --build
fi

exec "${ROOT_DIR}/scripts/docker-build-and-run.sh" "$@"
