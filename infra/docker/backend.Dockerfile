ARG MODE

FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./

RUN go mod download

COPY backend/ ./

RUN if [ "$MODE" = "prod" ]; then go build -o main ./cmd/main.go; fi

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main ./main

COPY --from=builder /app ./src

EXPOSE 8080

CMD sh -c "if [ '$MODE' = 'dev' ]; then go run ./src/cmd/main.go; else ./main; fi"
