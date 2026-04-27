#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH='' cd -- "$(dirname "$0")/.." && pwd)"
BACKEND_DIR="${ROOT_DIR}/backend"

AIR_BIN="${GOBIN:-$(go env GOPATH)/bin}/air"
if command -v air >/dev/null 2>&1; then
  AIR_BIN="air"
elif [ ! -x "${AIR_BIN}" ]; then
  echo "air is not installed."
  echo "Install it with: go install github.com/air-verse/air@v1.60.0"
  exit 1
fi

cd "${BACKEND_DIR}"
exec "${AIR_BIN}" -c .air.toml
