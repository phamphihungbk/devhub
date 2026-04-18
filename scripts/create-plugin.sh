#!/usr/bin/env sh
set -eu

usage() {
  cat <<'EOF'
Usage:
  ./scripts/create-plugin.sh --name NAME --type scaffolder|runner [options]

Options:
  --name NAME           Plugin name, used for the folder name
  --type TYPE           Plugin type: scaffolder or runner
  --scope SCOPE         Plugin scope: global, project, or environment (default: global)
  --version VERSION     Plugin version (default: 0.1.0)
  --description TEXT    Plugin description
  --language LANG       Starter language: python or shell (default: python)
  --force               Create files even if the target folder already exists
  --help                Show this help

Examples:
  ./scripts/create-plugin.sh --name node-http-api --type scaffolder --description "Scaffold Node APIs"
  ./scripts/create-plugin.sh --name deployment-sync --type runner --language shell
EOF
}

ROOT_DIR="$(CDPATH='' cd -- "$(dirname "$0")/.." && pwd)"

NAME=""
TYPE=""
SCOPE="global"
VERSION="0.1.0"
DESCRIPTION=""
LANGUAGE="python"
FORCE="0"
ENTRYPOINT_FILE="run.py"

while [ $# -gt 0 ]; do
  case "$1" in
    --name)
      NAME="${2:-}"
      shift 2
      ;;
    --type)
      TYPE="${2:-}"
      shift 2
      ;;
    --scope)
      SCOPE="${2:-}"
      shift 2
      ;;
    --version)
      VERSION="${2:-}"
      shift 2
      ;;
    --description)
      DESCRIPTION="${2:-}"
      shift 2
      ;;
    --language)
      LANGUAGE="${2:-}"
      shift 2
      ;;
    --force)
      FORCE="1"
      shift
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

if [ -z "$NAME" ] || [ -z "$TYPE" ]; then
  echo "--name and --type are required" >&2
  usage >&2
  exit 1
fi

case "$TYPE" in
  scaffolder)
    TYPE_DIR="scaffolders"
    ;;
  runner)
    TYPE_DIR="runners"
    ;;
  *)
    echo "Invalid --type: $TYPE" >&2
    exit 1
    ;;
esac

case "$SCOPE" in
  global|project|environment) ;;
  *)
    echo "Invalid --scope: $SCOPE" >&2
    exit 1
    ;;
esac

SLUG="$(printf '%s' "$NAME" | tr '[:upper:]' '[:lower:]' | tr ' ' '-' )"
PLUGIN_DIR="${ROOT_DIR}/plugins/${TYPE_DIR}/${SLUG}"
ENTRYPOINT_PATH="/app/plugins/${TYPE_DIR}/${SLUG}/${ENTRYPOINT_FILE}"

if [ -e "$PLUGIN_DIR" ] && [ "$FORCE" != "1" ]; then
  echo "Plugin directory already exists: $PLUGIN_DIR" >&2
  echo "Re-run with --force if you want to overwrite starter files." >&2
  exit 1
fi

mkdir -p "$PLUGIN_DIR"

printf '' > "${PLUGIN_DIR}/__init__.py"

cat > "${PLUGIN_DIR}/schema.json" <<EOF
{
  "type": "object",
  "required": ["service_name"],
  "properties": {
    "service_name": {
      "type": "string",
      "description": "Primary service or job name for this plugin."
    }
  },
  "additionalProperties": true
}
EOF

if [ "$LANGUAGE" = "python" ]; then
  cat > "${PLUGIN_DIR}/${ENTRYPOINT_FILE}" <<EOF
import json
import sys


def main() -> None:
    raw = sys.stdin.read().strip()
    payload = json.loads(raw) if raw else {}

    print(json.dumps({
        "status": "ok",
        "output": {
            "plugin": "${SLUG}",
            "received": payload,
        }
    }))


if __name__ == "__main__":
    main()
EOF
else
  cat > "${PLUGIN_DIR}/run.sh" <<EOF
#!/usr/bin/env sh
set -eu

PAYLOAD="\$(cat)"

printf '%s\n' "{\\"status\\":\\"ok\\",\\"output\\":{\\"plugin\\":\\"${SLUG}\\",\\"received\\":\$PAYLOAD}}"
EOF
  chmod +x "${PLUGIN_DIR}/run.sh"
fi

cat > "${PLUGIN_DIR}/README.md" <<EOF
# ${NAME}

Type: ${TYPE}
Scope: ${SCOPE}
Version: ${VERSION}

Entrypoint:
\`${ENTRYPOINT_PATH}\`
EOF

cat > "${PLUGIN_DIR}/plugin.yaml" <<EOF
name: ${SLUG}
type: ${TYPE}
version: ${VERSION}
description: ${DESCRIPTION:-${NAME}}
scope: ${SCOPE}
enabled: true
runtime: ${LANGUAGE}
entrypoint: ${ENTRYPOINT_FILE}
EOF

echo "Created plugin scaffold at:"
echo "  ${PLUGIN_DIR}"
echo
echo "Entrypoint:"
echo "  ${ENTRYPOINT_PATH}"
echo
echo "Register it with the API using:"
cat <<EOF
curl -X POST http://localhost:8080/plugins \\
  -H 'Content-Type: application/json' \\
  -d '{
    "name": "${NAME}",
    "version": "${VERSION}",
    "type": "${TYPE}",
    "entrypoint": "${ENTRYPOINT_PATH}",
    "scope": "${SCOPE}",
    "description": "${DESCRIPTION:-${NAME}}"
  }'
EOF
