ARG MODE=dev

FROM golang:1.23-alpine AS base

WORKDIR /app

RUN apk add --no-cache bash ca-certificates git python3 py3-pip

COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /app/backend
RUN go mod download

COPY backend/ /app/backend/
COPY plugins/ /app/plugins/

FROM base AS dev

WORKDIR /app/backend
EXPOSE 8080
CMD ["go", "run", ".", "serve"]

FROM base AS prod-builder

WORKDIR /app/backend
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/devhub .

FROM alpine:3.22 AS prod

WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=prod-builder /out/devhub /app/devhub

EXPOSE 8080
CMD ["/app/devhub", "serve"]
