ARG OTEL_COLLECTOR_VERSION=latest

FROM otel/opentelemetry-collector-contrib:${OTEL_COLLECTOR_VERSION}

COPY infra/otel/otel-collector-config.yaml /etc/otelcol/config.yaml

CMD ["--config=/etc/otelcol/config.yaml"]
