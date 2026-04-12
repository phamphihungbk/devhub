#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH='' cd -- "$(dirname "$0")/.." && pwd)"
DOMAIN="${1:-${DEVHUB_DOMAIN:-devhub.local}}"
API_DOMAIN="${2:-${DEVHUB_API_DOMAIN:-api.devhub.local}}"
CERT_DIR="${ROOT_DIR}/infra/certs"
CERT_FILE="${CERT_DIR}/${DEVHUB_SSL_CERT_FILE:-devhub.local.crt}"
HOSTS_LINE_PRIMARY="127.0.0.1 ${DOMAIN}"
HOSTS_LINE_API="127.0.0.1 ${API_DOMAIN}"

"${ROOT_DIR}/scripts/generate-dev-cert.sh" "${DOMAIN}" "${API_DOMAIN}"

if ! grep -Eq "^[[:space:]]*127\\.0\\.0\\.1[[:space:]]+${DOMAIN}([[:space:]]|\$)" /etc/hosts; then
  printf '%s\n' "Adding ${DOMAIN} to /etc/hosts (sudo may prompt)..."
  printf '%s\n' "${HOSTS_LINE_PRIMARY}" | sudo tee -a /etc/hosts >/dev/null
else
  printf '%s\n' "/etc/hosts already contains ${DOMAIN}"
fi

if ! grep -Eq "^[[:space:]]*127\\.0\\.0\\.1[[:space:]]+${API_DOMAIN}([[:space:]]|\$)" /etc/hosts; then
  printf '%s\n' "Adding ${API_DOMAIN} to /etc/hosts (sudo may prompt)..."
  printf '%s\n' "${HOSTS_LINE_API}" | sudo tee -a /etc/hosts >/dev/null
else
  printf '%s\n' "/etc/hosts already contains ${API_DOMAIN}"
fi

OS_NAME="$(uname -s)"
if [ "${OS_NAME}" = "Darwin" ]; then
  printf '%s\n' "Trusting ${CERT_FILE} in the macOS System keychain (sudo may prompt)..."
  sudo security add-trusted-cert \
    -d \
    -r trustRoot \
    -k /Library/Keychains/System.keychain \
    "${CERT_FILE}"
  printf '%s\n' "Local HTTPS is ready at https://${DOMAIN} and https://${API_DOMAIN}"
else
  printf '%s\n' "Hosts updated, but certificate trust was not automated for ${OS_NAME}."
  printf '%s\n' "Manually trust this certificate: ${CERT_FILE}"
  printf '%s\n' "Then open https://${DOMAIN} and https://${API_DOMAIN}"
fi
