# Multi-stage Dockerfile for BIAN-Go
# Stage 1: Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Generate code (GraphQL, etc.)
RUN go generate ./...

# Build the application
# CGO_ENABLED=0 for static binary
# -ldflags="-w -s" to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /build/bin/bian-go \
    ./cmd/server

# Stage 2: Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 bian && \
    adduser -D -u 1000 -G bian bian

# Set working directory
WORKDIR /app

# Copy binary from build stage
COPY --from=builder /build/bin/bian-go /app/bian-go

# Copy any static assets if needed
# COPY --from=builder /build/static /app/static

# Change ownership
RUN chown -R bian:bian /app

# Switch to non-root user
USER bian

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set default environment variables
ENV PORT=8080 \
    GRAPHQL_PLAYGROUND=false

# Run the application
CMD ["/app/bian-go"]
