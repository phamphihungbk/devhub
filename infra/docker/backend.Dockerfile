ARG MODE

FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache python3 py3-pip

COPY backend/ ./
COPY plugins/ ./plugins/

RUN go mod download
# RUN if [ "$MODE" = "prod" ]; then go build -o main ./cmd/main.go; fi

# FROM alpine:latest

# WORKDIR /app

# COPY --from=builder /app/main ./main

# COPY --from=builder /app ./src

EXPOSE 8080

# CMD sh -c "if [ '$MODE' = 'dev' ]; then go run ./src/cmd/main.go; else ./main; fi"
