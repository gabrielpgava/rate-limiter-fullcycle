FROM golang:1.25.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/rate-limiter ./cmd

FROM alpine:3.20

RUN apk add --no-cache ca-certificates curl
RUN adduser -D app

WORKDIR /app
COPY --from=builder /app/bin/rate-limiter /app/rate-limiter

EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=2s --retries=5 CMD curl -fsS http://localhost:8080/ || exit 1

USER app
ENTRYPOINT ["/app/rate-limiter"]
