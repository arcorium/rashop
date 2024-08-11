FROM golang:1.22.5 AS builder
LABEL authors="arcorium"

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY services/customer/ ./services/customer
COPY shared ./shared
COPY proto ./proto
COPY contract ./contract

WORKDIR /app/services/customer

RUN go mod tidy
RUN go mod download

# Hot reload package
RUN go install github.com/air-verse/air@latest

CMD ["air", "-c", ".air.toml"]