# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download and verify dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Generate code
RUN go generate ./...

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o bin/bian-go \
    ./cmd/server 2>/dev/null || \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o bin/bian-go \
    ./examples/basic

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata wget

# Create non-root user
RUN addgroup -g 1000 bian && \
    adduser -D -s /bin/sh -u 1000 -G bian bian

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/bin/bian-go ./bian-go

# Set ownership and permissions
RUN chown bian:bian ./bian-go && \
    chmod +x ./bian-go

# Switch to non-root user
USER bian

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Environment defaults
ENV PORT=8080
ENV GRAPHQL_PLAYGROUND=false

# Run the binary
CMD ["./bian-go"]