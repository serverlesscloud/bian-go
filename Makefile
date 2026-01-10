.PHONY: build test generate lint run clean help

# Default target
help:
	@echo "Available targets:"
	@echo "  build     - Build the project"
	@echo "  test      - Run tests"
	@echo "  generate  - Run go generate (GraphQL code generation)"
	@echo "  lint      - Run linter (requires golangci-lint)"
	@echo "  run       - Run example server"
	@echo "  clean     - Clean build artifacts"

# Build the project
build:
	go build ./...

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Generate GraphQL code
generate:
	go generate ./...

# Run linter (requires golangci-lint to be installed)
lint:
	golangci-lint run

# Run example server
run:
	cd examples/basic && go run main.go

# Clean build artifacts
clean:
	go clean ./...
	rm -f examples/basic/basic
	rm -f examples/basic/basic.exe

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Run all checks (format, vet, test)
check: fmt vet test

# Build Docker image
docker-build:
	docker build -t bian-go .

# Run Docker container
docker-run:
	docker run -p 8080:8080 bian-go