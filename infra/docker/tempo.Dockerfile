ARG TEMPO_VERSION=latest

FROM grafana/tempo:${TEMPO_VERSION}

USER 0:0

COPY infra/otel/tempo.yaml /etc/tempo.yaml

CMD ["-config.file=/etc/tempo.yaml"]
