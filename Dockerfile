# Multi-stage build using Chainguard images for security
# Build stage
FROM cgr.dev/chainguard/go:latest-dev AS builder

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o latency-exporter \
    ./cmd/latency-exporter

# Final stage - use minimal distroless image
FROM cgr.dev/chainguard/static:latest

# Add fping for ICMP measurements
# Note: We'll need to use a slightly larger base image that includes fping
FROM cgr.dev/chainguard/wolfi-base:latest

# Install fping for ICMP functionality
RUN apk add --no-cache fping

# Create a non-root user
RUN addgroup -g 10001 -S latency && \
    adduser -u 10001 -S latency -G latency

# Create necessary directories
RUN mkdir -p /var/latency-parser && \
    chown -R latency:latency /var/latency-parser

# Copy the binary from builder stage
COPY --from=builder /app/latency-exporter /usr/local/bin/latency-exporter

# Make binary executable and owned by latency user
RUN chmod +x /usr/local/bin/latency-exporter

# Copy default config (optional)
COPY --chown=latency:latency config/example.yml /var/latency-parser/config.yml

# Switch to non-root user
USER latency

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set default environment variables
ENV LATENCY_PARSER_CONFIG_PATH=/var/latency-parser/config.yml

# Run the application
ENTRYPOINT ["/usr/local/bin/latency-exporter"]
