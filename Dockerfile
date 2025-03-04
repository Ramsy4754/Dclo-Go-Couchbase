FROM golang:1.21 AS builder

WORKDIR /app

COPY . .

RUN GOOS=linux GOARCH=arm64 go build -o couchbase_bridge .

FROM debian:latest

WORKDIR /app

COPY --from=builder /app/couchbase_bridge .

CMD ["./couchbase_bridge"]