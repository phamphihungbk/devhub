DC="$(which docker) compose -f ./docker-compose.dev.yml"

HOST_UID="$(id -u)"
HOST_GID="$(id -g)"

SERVICE="api"

export HOST_UID
export HOST_GID
export DOCKER_BUILDKIT=1

RUN="${DC} run --rm"

set -x

case $1 in

"start")
  ${DC} up -d
  if [ "$2" == '-f' ]; then
    ${DC} logs -f
  fi
  ;;

"restart")
  ${DC} down || true
  exec "$0" start "${@:2}"
  ;;

"build")
  ${DC} build
  ;;

"run")
  ${RUN} "${@:2}"
  ;;

"shell")
  ${DC} run -u 0:0 --rm ${SERVICE} /bin/bash
  ;;

*)
  ${DC} "${@:1}"
  ;;

esac