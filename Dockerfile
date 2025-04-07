FROM golang:1.22-alpine AS builder

RUN apk add --no-cache make git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

FROM alpine
RUN apk add --no-cache ca-certificates

WORKDIR /auth

ENV CONFIG_PATH=./config/config.yaml

COPY --from=builder /app/bin/auth .
COPY --from=builder /app/config ./config

CMD ["./auth"]