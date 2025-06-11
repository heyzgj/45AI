# Multi-stage Dockerfile for 45AI Backend
# Stage 1: Build stage
FROM golang:1.21-alpine AS builder

# Set build arguments
ARG BUILD_VERSION=dev
ARG BUILD_TIME
ARG CGO_ENABLED=0

# Install git and ca-certificates for Go modules and HTTPS
RUN apk add --no-cache git ca-certificates tzdata

# Create non-root user for building
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=${BUILD_VERSION} -X main.buildTime=${BUILD_TIME}" \
    -a -installsuffix cgo \
    -o main ./cmd/api

# Verify the binary
RUN ./main --version || true

# Stage 2: Production stage
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy the binary from builder stage
COPY --from=builder /build/main /app/main

# Create necessary directories
COPY --from=builder --chown=appuser:appuser /tmp /tmp

# Use non-root user
USER appuser

# Set working directory
WORKDIR /app

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/main", "--health-check"] || exit 1

# Set environment variables
ENV GIN_MODE=release
ENV PORT=8080

# Labels for metadata
LABEL maintainer="45AI Team" \
      version="${BUILD_VERSION}" \
      description="45AI Image Generation Backend API" \
      org.opencontainers.image.title="45AI Backend" \
      org.opencontainers.image.description="AI-powered image generation backend service" \
      org.opencontainers.image.version="${BUILD_VERSION}" \
      org.opencontainers.image.created="${BUILD_TIME}" \
      org.opencontainers.image.vendor="45AI" \
      org.opencontainers.image.licenses="MIT"

# Run the application
ENTRYPOINT ["/app/main"]
CMD [] 