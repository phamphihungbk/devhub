#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH='' cd -- "$(dirname "$0")/.." && pwd)"
CERT_DIR="${ROOT_DIR}/infra/certs"
DOMAIN="${1:-${DEVHUB_DOMAIN:-devhub.local}}"
API_DOMAIN="${2:-${DEVHUB_API_DOMAIN:-api.devhub.local}}"
CERT_FILE="${DEVHUB_SSL_CERT_FILE:-devhub.local.crt}"
KEY_FILE="${DEVHUB_SSL_KEY_FILE:-devhub.local.key}"

mkdir -p "${CERT_DIR}"

openssl req \
  -x509 \
  -nodes \
  -newkey rsa:2048 \
  -keyout "${CERT_DIR}/${KEY_FILE}" \
  -out "${CERT_DIR}/${CERT_FILE}" \
  -days 825 \
  -subj "/CN=${DOMAIN}" \
  -addext "subjectAltName=DNS:${DOMAIN},DNS:${API_DOMAIN},DNS:localhost,IP:127.0.0.1"

printf 'Created local TLS files:\n- %s\n- %s\n' "${CERT_DIR}/${CERT_FILE}" "${CERT_DIR}/${KEY_FILE}"
