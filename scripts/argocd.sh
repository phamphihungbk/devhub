#!/usr/bin/env sh
set -eu

ARGOCD_NAMESPACE="${ARGOCD_NAMESPACE:-argocd}"
DEVHUB_NAMESPACE="${DEVHUB_NAMESPACE:-devhub}"
ARGOCD_VERSION="${ARGOCD_VERSION:-stable}"
UI_PORT="${ARGOCD_UI_PORT:-8081}"
ARGOCD_UI_HOST="${ARGOCD_UI_HOST:-argocd.devhub.local}"
GITOPS_REPO_URL="${GITOPS_REPO_URL:-http://host.minikube.internal:3000/phamphihungbk/gitops-repo.git}"
GITOPS_REPO_BRANCH="${GITOPS_REPO_BRANCH:-main}"
GITOPS_ENV_GLOB="${GITOPS_ENV_GLOB:-envs/dev/*.yaml}"
ROOT_DIR="$(CDPATH='' cd -- "$(dirname "$0")/.." && pwd)"
APP_MANIFEST="${ROOT_DIR}/infra/kubernetes/argocd/devhub.yaml"
INGRESS_MANIFEST="${ROOT_DIR}/infra/kubernetes/argocd/argocd-ui-ingress.yaml"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

ensure_namespace() {
  if ! kubectl get namespace "$1" >/dev/null 2>&1; then
    kubectl create namespace "$1" >/dev/null
  fi
}

configure_admin_access() {
  kubectl -n "${ARGOCD_NAMESPACE}" patch configmap argocd-cm \
    --type merge \
    -p '{"data":{"accounts.admin":"apiKey, login"}}' >/dev/null
  kubectl -n "${ARGOCD_NAMESPACE}" patch configmap argocd-rbac-cm \
    --type merge \
    -p '{"data":{"policy.default":"role:admin"}}' >/dev/null
  kubectl -n "${ARGOCD_NAMESPACE}" rollout restart deployment/argocd-server >/dev/null
  kubectl -n "${ARGOCD_NAMESPACE}" rollout status deployment/argocd-server --timeout=180s
}

install_argocd() {
  ensure_namespace "${ARGOCD_NAMESPACE}"
  kubectl apply --server-side -n "${ARGOCD_NAMESPACE}" \
    -f "https://raw.githubusercontent.com/argoproj/argo-cd/${ARGOCD_VERSION}/manifests/install.yaml"
  configure_admin_access
  kubectl rollout status deployment/argocd-server -n "${ARGOCD_NAMESPACE}" --timeout=180s
}

wait_for_server() {
  kubectl rollout status deployment/argocd-server -n "${ARGOCD_NAMESPACE}" --timeout=180s
}

print_password() {
  kubectl -n "${ARGOCD_NAMESPACE}" get secret argocd-initial-admin-secret \
    -o jsonpath='{.data.password}' | base64 --decode
  printf '\n'
}

apply_devhub_app() {
  ensure_namespace "${DEVHUB_NAMESPACE}"
  require_cmd envsubst
  export GITOPS_REPO_URL GITOPS_REPO_BRANCH GITOPS_ENV_GLOB
  envsubst '${GITOPS_REPO_URL} ${GITOPS_REPO_BRANCH} ${GITOPS_ENV_GLOB}' < "${APP_MANIFEST}" | kubectl apply -f -
}

enable_minikube_ingress() {
  require_cmd minikube
  minikube addons enable ingress
}

apply_ui_ingress() {
  enable_minikube_ingress
  kubectl apply -f "${INGRESS_MANIFEST}"
}

print_ui_domain_instructions() {
  require_cmd minikube
  minikube_ip="$(minikube ip)"
  echo "Argo CD UI domain: https://${ARGOCD_UI_HOST}"
  echo "Add this hosts entry on your machine:"
  echo "${minikube_ip} ${ARGOCD_UI_HOST}"
  echo "Then open https://${ARGOCD_UI_HOST}"
}

port_forward_ui() {
  wait_for_server
  echo "Argo CD UI: http://127.0.0.1:${UI_PORT}"
  echo "Username: admin"
  printf 'Password: '
  print_password
  echo "Starting port-forward. Press Ctrl+C to stop."
  exec kubectl -n "${ARGOCD_NAMESPACE}" port-forward svc/argocd-server "${UI_PORT}:443"
}

login_cli() {
  require_cmd argocd
  wait_for_server
  password="$(kubectl -n "${ARGOCD_NAMESPACE}" get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 --decode)"
  argocd login "127.0.0.1:${UI_PORT}" --username admin --password "${password}" --insecure
}

print_auth_token() {
  require_cmd argocd
  login_cli >/dev/null
  token="$(argocd account generate-token)"
  printf 'export ARGOCD_AUTH_TOKEN=%s\n' "${token}"
}

usage() {
  cat <<EOF
Usage: ./scripts/argocd.sh <command>

Commands:
  install   Install Argo CD into the current Kubernetes cluster
  configure Configure Argo CD admin for login, apiKey, and admin RBAC
  app       Apply the DevHub Argo CD ApplicationSet manifest
  ingress   Enable Minikube ingress and apply the Argo CD UI hostname
  domain    Print the Argo CD UI domain and the required /etc/hosts entry
  password  Print the Argo CD admin password
  token     Log in and print an ARGOCD_AUTH_TOKEN export line
  ui        Start a local port-forward for the Argo CD UI
  login     Log in with the Argo CD CLI after the UI port-forward is running
  all       Install Argo CD, apply the DevHub app, configure ingress, then start the UI

Environment:
  ARGOCD_NAMESPACE   Default: argocd
  DEVHUB_NAMESPACE   Default: devhub
  ARGOCD_VERSION     Default: stable
  ARGOCD_UI_PORT     Default: 8081
  ARGOCD_UI_HOST     Default: argocd.devhub.local
  GITOPS_REPO_URL    Default: https://gitea.devhub.local/phamphihungbk/gitops-repo.git
  GITOPS_REPO_BRANCH Default: main
  GITOPS_ENV_GLOB    Default: envs/dev/*.yaml
EOF
}

require_cmd kubectl

cmd="${1:-ui}"

case "${cmd}" in
  install)
    install_argocd
    ;;
  configure)
    configure_admin_access
    ;;
  app)
    apply_devhub_app
    ;;
  ingress)
    apply_ui_ingress
    ;;
  domain)
    print_ui_domain_instructions
    ;;
  password)
    print_password
    ;;
  token)
    print_auth_token
    ;;
  ui)
    port_forward_ui
    ;;
  login)
    login_cli
    ;;
  all)
    install_argocd
    apply_devhub_app
    apply_ui_ingress
    print_ui_domain_instructions
    port_forward_ui
    ;;
  -h|--help|help)
    usage
    ;;
  *)
    echo "Unknown command: ${cmd}" >&2
    usage >&2
    exit 1
    ;;
esac
