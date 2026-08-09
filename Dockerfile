# Build the Go binary
FROM golang:1.26-alpine AS build
WORKDIR /src

# Cache module downloads across builds
COPY server/go.mod server/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Build (with go build cache persisted via BuildKit)
COPY server/ .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/appdroid ./cmd/server

# Run the app
FROM scratch
WORKDIR /app
ARG PORT=3000
ENV PORT=$PORT
# Release mode enables static asset caching (dev/debug mode disables it).
ENV GIN_MODE=release
EXPOSE $PORT

COPY --from=build /out/appdroid /app/appdroid

CMD ["/app/appdroid"]
