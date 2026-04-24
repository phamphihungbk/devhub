ARG GRAFANA_VERSION=latest

FROM grafana/grafana:${GRAFANA_VERSION}

COPY infra/grafana/provisioning /etc/grafana/provisioning
