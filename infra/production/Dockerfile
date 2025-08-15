# Multi-stage Dockerfile for APIWeaver (Go + React + MongoDB integration)
# Optimized for security, performance, and minimal image size

# ===============================
# Stage 1: Frontend Build (React + TypeScript + Vite)
# ===============================
FROM node:20-alpine AS frontend-builder

# Install build dependencies and security patches
RUN apk add --no-cache \
    git \
    python3 \
    make \
    g++ \
    && apk upgrade --no-cache

# Set working directory
WORKDIR /app/web

# Copy package files first for better caching
COPY web/package*.json ./
COPY web/tsconfig*.json ./
COPY web/vite.config.ts ./
COPY web/tailwind.config.js ./
COPY web/postcss.config.js ./

# Install dependencies with clean install for reproducible builds
RUN npm ci \
    && npm cache clean --force

# Copy source code
COPY web/src ./src
COPY web/public ./public
COPY web/index.html ./

# Build frontend for production
ENV NODE_ENV=production
ENV VITE_API_BASE_URL=/api
ENV VITE_APP_VERSION=1.0.0

RUN npm run build:ci \
    && ls -la dist/ \
    && echo "Frontend build completed successfully"

# ===============================
# Stage 2: Go Build Stage
# ===============================
FROM golang:1.24.4-alpine AS go-builder

# Install build dependencies and security patches
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    && apk upgrade --no-cache

# Set build environment
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create non-root user for build process
RUN adduser -D -s /bin/sh apiweaver

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download \
    && go mod verify

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY testutil/ ./testutil/

# Copy frontend build from previous stage
COPY --from=frontend-builder /app/web/dist ./web/dist

# Build version information
ARG VERSION=dev
ARG COMMIT_SHA=unknown
ARG BUILD_TIME=unknown

# Build the Go binary with embedded frontend assets and optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.commitSHA=${COMMIT_SHA} -X main.buildTime=${BUILD_TIME} -extldflags '-static'" \
    -a -installsuffix cgo \
    -trimpath \
    -tags netgo \
    -o apiweaver \
    ./cmd/apiweaver \
    && chmod +x apiweaver \
    && echo "Go build completed successfully"

# Verify the binary
RUN ./apiweaver --version

# ===============================
# Stage 3: Production Runtime
# ===============================
FROM alpine:3.19 AS production

# Install runtime dependencies and security patches
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    jq \
    && apk upgrade --no-cache \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1001 -S apiweaver \
    && adduser -u 1001 -S apiweaver -G apiweaver

# Create necessary directories
RUN mkdir -p /app/data /app/logs /app/config \
    && chown -R apiweaver:apiweaver /app

# Set working directory
WORKDIR /app

# Copy binary from build stage
COPY --from=go-builder --chown=apiweaver:apiweaver /app/apiweaver ./apiweaver

# Copy configuration templates
COPY --chown=apiweaver:apiweaver config.yaml.sample ./config/config.yaml.template

# Create runtime configuration
RUN echo '# APIWeaver Runtime Configuration' > ./config/config.yaml \
    && echo 'server:' >> ./config/config.yaml \
    && echo '  host: "0.0.0.0"' >> ./config/config.yaml \
    && echo '  port: 8080' >> ./config/config.yaml \
    && echo '  read_timeout: "30s"' >> ./config/config.yaml \
    && echo '  write_timeout: "30s"' >> ./config/config.yaml \
    && echo '  idle_timeout: "60s"' >> ./config/config.yaml \
    && echo '' >> ./config/config.yaml \
    && echo 'database:' >> ./config/config.yaml \
    && echo '  mongodb:' >> ./config/config.yaml \
    && echo '    uri: "mongodb://mongodb:27017"' >> ./config/config.yaml \
    && echo '    database: "apiweaver"' >> ./config/config.yaml \
    && echo '    timeout: "10s"' >> ./config/config.yaml \
    && echo '    max_pool_size: 100' >> ./config/config.yaml \
    && echo '' >> ./config/config.yaml \
    && echo 'api:' >> ./config/config.yaml \
    && echo '  rate_limit:' >> ./config/config.yaml \
    && echo '    requests_per_minute: 100' >> ./config/config.yaml \
    && echo '    burst: 10' >> ./config/config.yaml \
    && echo '  cors:' >> ./config/config.yaml \
    && echo '    enabled: true' >> ./config/config.yaml \
    && echo '    allowed_origins: ["*"]' >> ./config/config.yaml \
    && echo '' >> ./config/config.yaml \
    && echo 'logging:' >> ./config/config.yaml \
    && echo '  level: "info"' >> ./config/config.yaml \
    && echo '  format: "json"' >> ./config/config.yaml \
    && echo '  file: "/app/logs/apiweaver.log"' >> ./config/config.yaml \
    && echo '' >> ./config/config.yaml \
    && echo 'security:' >> ./config/config.yaml \
    && echo '  api_keys:' >> ./config/config.yaml \
    && echo '    enabled: true' >> ./config/config.yaml \
    && echo '    header: "X-API-Key"' >> ./config/config.yaml \
    && echo '  jwt:' >> ./config/config.yaml \
    && echo '    secret: "${JWT_SECRET}"' >> ./config/config.yaml \
    && echo '    expiry: "24h"' >> ./config/config.yaml

# Create health check script
RUN echo '#!/bin/sh' > /app/healthcheck.sh \
    && echo 'set -e' >> /app/healthcheck.sh \
    && echo 'if [ -f /app/apiweaver ]; then' >> /app/healthcheck.sh \
    && echo '  # Check if the binary is responsive' >> /app/healthcheck.sh \
    && echo '  timeout 5s /app/apiweaver --version >/dev/null 2>&1' >> /app/healthcheck.sh \
    && echo '  # Check if the server is responding (when running as server)' >> /app/healthcheck.sh \
    && echo '  if [ "${SERVER_MODE:-false}" = "true" ]; then' >> /app/healthcheck.sh \
    && echo '    curl -f http://localhost:${PORT:-8080}/health >/dev/null 2>&1 || exit 1' >> /app/healthcheck.sh \
    && echo '  fi' >> /app/healthcheck.sh \
    && echo 'else' >> /app/healthcheck.sh \
    && echo '  exit 1' >> /app/healthcheck.sh \
    && echo 'fi' >> /app/healthcheck.sh \
    && chmod +x /app/healthcheck.sh

# Create entrypoint script
RUN echo '#!/bin/sh' > /app/entrypoint.sh \
    && echo 'set -e' >> /app/entrypoint.sh \
    && echo '' >> /app/entrypoint.sh \
    && echo '# Set default values' >> /app/entrypoint.sh \
    && echo 'export PORT=${PORT:-8080}' >> /app/entrypoint.sh \
    && echo 'export HOST=${HOST:-0.0.0.0}' >> /app/entrypoint.sh \
    && echo 'export LOG_LEVEL=${LOG_LEVEL:-info}' >> /app/entrypoint.sh \
    && echo 'export SERVER_MODE=${SERVER_MODE:-true}' >> /app/entrypoint.sh \
    && echo '' >> /app/entrypoint.sh \
    && echo '# MongoDB connection configuration' >> /app/entrypoint.sh \
    && echo 'export MONGODB_URI=${MONGODB_URI:-mongodb://mongodb:27017}' >> /app/entrypoint.sh \
    && echo 'export MONGODB_DATABASE=${MONGODB_DATABASE:-apiweaver}' >> /app/entrypoint.sh \
    && echo 'export MONGODB_TIMEOUT=${MONGODB_TIMEOUT:-10s}' >> /app/entrypoint.sh \
    && echo '' >> /app/entrypoint.sh \
    && echo '# Security configuration' >> /app/entrypoint.sh \
    && echo 'if [ -z "$JWT_SECRET" ]; then' >> /app/entrypoint.sh \
    && echo '  echo "Warning: JWT_SECRET not set, generating random secret"' >> /app/entrypoint.sh \
    && echo '  export JWT_SECRET=$(head -c 32 /dev/urandom | base64)' >> /app/entrypoint.sh \
    && echo 'fi' >> /app/entrypoint.sh \
    && echo '' >> /app/entrypoint.sh \
    && echo '# Wait for MongoDB to be ready if in server mode' >> /app/entrypoint.sh \
    && echo 'if [ "$SERVER_MODE" = "true" ] && [ "$WAIT_FOR_MONGODB" = "true" ]; then' >> /app/entrypoint.sh \
    && echo '  echo "Waiting for MongoDB to be ready..."' >> /app/entrypoint.sh \
    && echo '  until nc -z mongodb 27017; do' >> /app/entrypoint.sh \
    && echo '    echo "MongoDB is unavailable - sleeping"' >> /app/entrypoint.sh \
    && echo '    sleep 2' >> /app/entrypoint.sh \
    && echo '  done' >> /app/entrypoint.sh \
    && echo '  echo "MongoDB is ready!"' >> /app/entrypoint.sh \
    && echo 'fi' >> /app/entrypoint.sh \
    && echo '' >> /app/entrypoint.sh \
    && echo '# Execute the main command' >> /app/entrypoint.sh \
    && echo 'if [ "$SERVER_MODE" = "true" ]; then' >> /app/entrypoint.sh \
    && echo '  echo "Starting APIWeaver server on ${HOST}:${PORT}"' >> /app/entrypoint.sh \
    && echo '  exec /app/apiweaver serve --host="${HOST}" --port="${PORT}" --config=/app/config/config.yaml "$@"' >> /app/entrypoint.sh \
    && echo 'else' >> /app/entrypoint.sh \
    && echo '  echo "Running APIWeaver CLI"' >> /app/entrypoint.sh \
    && echo '  exec /app/apiweaver "$@"' >> /app/entrypoint.sh \
    && echo 'fi' >> /app/entrypoint.sh \
    && chmod +x /app/entrypoint.sh

# Switch to non-root user
USER apiweaver

# Set environment variables
ENV PORT=8080 \
    HOST=0.0.0.0 \
    LOG_LEVEL=info \
    SERVER_MODE=true \
    WAIT_FOR_MONGODB=true \
    GIN_MODE=release

# Expose port
EXPOSE 8080

# Add health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD /app/healthcheck.sh

# Set entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]

# Default command (can be overridden)
CMD []

# Add labels for better maintainability
LABEL maintainer="APIWeaver Team" \
      description="APIWeaver - Markdown to API Specification Parser with Web UI" \
      version="${VERSION:-dev}" \
      build-date="${BUILD_TIME:-unknown}" \
      vcs-url="https://github.com/sukhera/APIWeaver" \
      vcs-ref="${COMMIT_SHA:-unknown}" \
      schema-version="1.0"