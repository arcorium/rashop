# Builder
FROM golang:1.22.5 AS builder
LABEL authors="arcorium"

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

ARG SERVICE

COPY services/${SERVICE} services/${SERVICE}
COPY contract ./contract
COPY proto ./proto
COPY shared ./shared

WORKDIR /app/services/${SERVICE}

RUN go mod tidy
RUN go mod download

RUN go build -o build/server "./cmd/query_server/"

# Run tester
FROM builder AS test-runner

RUN go test ./...

# Runner
FROM alpine:latest AS runner

COPY --from=builder /app/build/* /app/

WORKDIR /app

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.13 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

ENTRYPOINT ["./server"]