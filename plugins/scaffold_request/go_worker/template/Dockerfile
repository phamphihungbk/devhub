FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/worker ./cmd/worker

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/worker /app/worker
CMD ["/app/worker"]
