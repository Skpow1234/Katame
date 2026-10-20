# Multi-arch (amd64, arm64): build with
#   docker buildx build --platform linux/amd64,linux/arm64 -t katame:tag --push .
# Or: make docker-buildx
# ---- Build stage ----
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src

# Cache dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w -X github.com/Skpow1234/Katame/internal/cli.Version=$(git describe --tags --always --dirty 2>/dev/null || echo dev)" \
    -o /bin/katame ./cmd/katame

# ---- Runtime stage ----
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /bin/katame /usr/local/bin/katame

USER nonroot:nonroot
WORKDIR /work

ENTRYPOINT ["katame"]
