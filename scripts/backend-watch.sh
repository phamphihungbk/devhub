#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH='' cd -- "$(dirname "$0")/.." && pwd)"
BACKEND_DIR="${ROOT_DIR}/backend"

if ! command -v air >/dev/null 2>&1; then
  echo "air is not installed."
  echo "Install it with: go install github.com/air-verse/air@latest"
  exit 1
fi

cd "${BACKEND_DIR}"
exec air -c .air.toml
