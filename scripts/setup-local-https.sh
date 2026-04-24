#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH='' cd -- "$(dirname "$0")/.." && pwd)"
DOMAIN="${1:-${DEVHUB_DOMAIN:-devhub.local}}"
API_DOMAIN="${2:-${DEVHUB_API_DOMAIN:-api.devhub.local}}"
GITEA_DOMAIN="${3:-${GITEA_DOMAIN:-gitea.devhub.local}}"
GRAFANA_DOMAIN="${4:-${GRAFANA_DOMAIN:-grafana.devhub.local}}"
ARGOCD_DOMAIN="${DEVHUB_ARGOCD_DOMAIN:-argocd.devhub.local}"
CERT_DIR="${ROOT_DIR}/infra/certs"
CERT_FILE="${CERT_DIR}/${DEVHUB_SSL_CERT_FILE:-devhub.local.crt}"

add_host_entry() {
  host_ip="$1"
  host_name="$2"

  if grep -Eq "^[[:space:]]*${host_ip//./\\.}[[:space:]]+${host_name}([[:space:]]|\$)" /etc/hosts; then
    printf '%s\n' "/etc/hosts already contains ${host_name} -> ${host_ip}"
    return 0
  fi

  printf '%s\n' "Adding ${host_name} -> ${host_ip} to /etc/hosts (sudo may prompt)..."
  printf '%s\n' "${host_ip} ${host_name}" | sudo tee -a /etc/hosts >/dev/null
}

"${ROOT_DIR}/scripts/generate-dev-cert.sh" "${DOMAIN}" "${API_DOMAIN}" "${GITEA_DOMAIN}" "${GRAFANA_DOMAIN}"

add_host_entry "127.0.0.1" "${DOMAIN}"
add_host_entry "127.0.0.1" "${API_DOMAIN}"
add_host_entry "127.0.0.1" "${GITEA_DOMAIN}"
add_host_entry "127.0.0.1" "${GRAFANA_DOMAIN}"

ARGOCD_IP="${DEVHUB_ARGOCD_IP:-}"
if [ -z "${ARGOCD_IP}" ] && command -v minikube >/dev/null 2>&1; then
  ARGOCD_IP="$(minikube ip 2>/dev/null || true)"
fi

if [ -n "${ARGOCD_IP}" ]; then
  add_host_entry "${ARGOCD_IP}" "${ARGOCD_DOMAIN}"
else
  printf '%s\n' "Skipping ${ARGOCD_DOMAIN}: set DEVHUB_ARGOCD_IP or install/start Minikube to auto-detect its IP."
fi

OS_NAME="$(uname -s)"
if [ "${OS_NAME}" = "Darwin" ]; then
  printf '%s\n' "Trusting ${CERT_FILE} in the macOS System keychain (sudo may prompt)..."
  sudo security add-trusted-cert \
    -d \
    -r trustRoot \
    -k /Library/Keychains/System.keychain \
    "${CERT_FILE}"
  printf '%s\n' "Local HTTPS is ready at https://${DOMAIN}, https://${API_DOMAIN}, https://${GITEA_DOMAIN}, and https://${GRAFANA_DOMAIN}"
  if [ -n "${ARGOCD_IP}" ]; then
    printf '%s\n' "Argo CD UI host is mapped at https://${ARGOCD_DOMAIN}"
  fi
else
  printf '%s\n' "Hosts updated, but certificate trust was not automated for ${OS_NAME}."
  printf '%s\n' "Manually trust this certificate: ${CERT_FILE}"
  printf '%s\n' "Then open https://${DOMAIN}, https://${API_DOMAIN}, https://${GITEA_DOMAIN}, and https://${GRAFANA_DOMAIN}"
  if [ -n "${ARGOCD_IP}" ]; then
    printf '%s\n' "Argo CD UI host is mapped at https://${ARGOCD_DOMAIN}"
  fi
fi
