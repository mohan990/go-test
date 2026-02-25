# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o server \
    ./cmd/server

# Final stage - minimal runtime image
FROM alpine:latest

# Install CA certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/server .

# Use non-root user
USER appuser

# Expose port (Cloud Run uses PORT env var, defaults to 8080)
EXPOSE 8080

# Health check (optional, Cloud Run will use /healthz endpoint)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/server", "-healthcheck"] || exit 1

# Run the binary
ENTRYPOINT ["/app/server"]
