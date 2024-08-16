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

RUN go build -o ../../build/migrate "./cmd/migrate/"

# Runner
FROM alpine:latest AS runner

COPY --from=builder /app/build/* /app/

WORKDIR /app

ENTRYPOINT ["./migrate"]