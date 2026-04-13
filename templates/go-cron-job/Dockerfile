FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o job ./cmd/job

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/job ./job
CMD ["./job"]
