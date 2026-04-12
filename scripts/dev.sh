#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH='' cd -- "$(dirname "$0")/.." && pwd)"
export COMPOSE_ENV=dev

if [ $# -eq 0 ]; then
  set -- up --build
fi

exec "${ROOT_DIR}/scripts/docker-build-and-run.sh" "$@"
