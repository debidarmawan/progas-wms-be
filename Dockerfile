# syntax=docker/dockerfile:1

# --- Build stage ---
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

# Cache module downloads (invalidates only when go.mod/go.sum change)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /app/bin/main \
    .

# --- Runtime stage ---
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -g 1000 app \
    && adduser -u 1000 -G app -D app

WORKDIR /app

COPY --from=builder /app/bin/main ./main

USER app

# Default app port for local/self-hosted deployments.
EXPOSE 3131

ENV GO_ENV=production

ENTRYPOINT ["./main"]
