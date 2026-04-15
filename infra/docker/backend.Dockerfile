ARG MODE=dev
ARG ARGOCD_VERSION=v2.14.12

FROM golang:1.23-alpine AS base
ARG ARGOCD_VERSION

WORKDIR /app

RUN apk add --no-cache bash ca-certificates git python3 py3-pip py3-jinja2 wget && \
    wget -O /usr/local/bin/argocd "https://github.com/argoproj/argo-cd/releases/download/${ARGOCD_VERSION}/argocd-linux-amd64" && \
    chmod +x /usr/local/bin/argocd

COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /app/backend
RUN go mod download

COPY backend/ /app/backend/
COPY plugins/ /app/plugins/
COPY templates/ /app/templates/

FROM base AS dev

WORKDIR /app/backend
EXPOSE 8080
CMD ["go", "run", ".", "serve"]

FROM base AS prod-builder

WORKDIR /app/backend
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/devhub .

FROM alpine:3.22 AS prod
ARG ARGOCD_VERSION

WORKDIR /app
RUN apk add --no-cache ca-certificates wget && \
    wget -O /usr/local/bin/argocd "https://github.com/argoproj/argo-cd/releases/download/${ARGOCD_VERSION}/argocd-linux-amd64" && \
    chmod +x /usr/local/bin/argocd
COPY --from=prod-builder /out/devhub /app/devhub

EXPOSE 8080
CMD ["/app/devhub", "serve"]
