.PHONY: help init build test test-coverage test-short generate lint fmt vet run run-dev clean clean-docker deps deps-update tidy validate ci docker-build docker-run openspec-validate openspec-list shell logs

# Set help as default target
.DEFAULT_GOAL := help

# Container runtime configuration (can be overridden)
DOCKER_RUN ?= docker run --rm -v "$(CURDIR):/app" -w /app
GO_IMAGE ?= golang:1.22-alpine
LINT_IMAGE ?= golangci/golangci-lint:v1.55-alpine

# Detect if running in container (for hybrid execution)
IN_CONTAINER ?= $(shell test -f /.dockerenv && echo 1 || echo 0)

# Command prefix - use Docker if not in container, direct if in container
ifeq ($(IN_CONTAINER),1)
	CMD_PREFIX =
else
	CMD_PREFIX = $(DOCKER_RUN) $(GO_IMAGE)
endif

# Help target
help:
	@echo "3 Musketeers Build System - BIAN Go"
	@echo ""
	@echo "Available targets:"
	@echo "  Development:"
	@echo "    init          - Initialize project (download dependencies)"
	@echo "    build         - Build the project"
	@echo "    generate      - Run go generate (GraphQL code generation)"
	@echo "    test          - Run tests"
	@echo "    test-coverage - Run tests with coverage"
	@echo "    test-short    - Run short tests only"
	@echo "    run           - Run example server (production mode)"
	@echo "    run-dev       - Run example server (development mode with hot reload)"
	@echo ""
	@echo "  Validation:"
	@echo "    fmt           - Format code"
	@echo "    vet           - Vet code"
	@echo "    lint          - Run linter"
	@echo "    validate      - Run all validation checks (fmt, vet, lint, test)"
	@echo ""
	@echo "  CI/CD:"
	@echo "    ci            - Run complete CI pipeline (init + validate + build)"
	@echo ""
	@echo "  Docker:"
	@echo "    docker-build  - Build Docker image"
	@echo "    docker-run    - Run Docker container"
	@echo ""
	@echo "  Utilities:"
	@echo "    clean         - Clean build artifacts"
	@echo "    clean-docker  - Clean Docker artifacts"
	@echo "    deps          - Download dependencies"
	@echo "    deps-update   - Update dependencies"
	@echo "    tidy          - Tidy go.mod"
	@echo "    shell         - Open interactive shell in container"
	@echo "    logs          - Show application logs"
	@echo ""
	@echo "  OpenSpec:"
	@echo "    openspec-validate - Validate OpenSpec changes"
	@echo "    openspec-list     - List OpenSpec changes"
	@echo ""
	@echo "Environment:"
	@echo "  DOCKER_RUN=$(DOCKER_RUN)"
	@echo "  GO_IMAGE=$(GO_IMAGE)"
	@echo "  IN_CONTAINER=$(IN_CONTAINER)"

# Initialize project
init:
	$(CMD_PREFIX) go mod download
	$(CMD_PREFIX) go mod verify

# Build the project
build:
	$(CMD_PREFIX) go build -o bin/bian-go ./cmd/server 2>/dev/null || $(CMD_PREFIX) go build ./...

# Run tests
test:
	$(CMD_PREFIX) go test ./...

# Run tests with coverage
test-coverage:
	$(CMD_PREFIX) go test -cover ./...

# Run short tests only
test-short:
	$(CMD_PREFIX) go test -short ./...

# Generate GraphQL code
generate:
	$(CMD_PREFIX) go generate ./...

# Run linter
lint:
ifeq ($(IN_CONTAINER),1)
	golangci-lint run
else
	$(DOCKER_RUN) $(LINT_IMAGE) golangci-lint run
endif

# Run example server (production mode)
run:
	docker compose up app

# Run example server (development mode with hot reload)
run-dev:
	docker compose up dev

# Clean build artifacts
clean:
	$(CMD_PREFIX) go clean ./...
	rm -rf bin/
	rm -f examples/basic/basic
	rm -f examples/basic/basic.exe
	rm -f tmp/
	rm -f build-errors.log

# Clean Docker artifacts
clean-docker:
	docker system prune -f
	docker volume prune -f

# Install dependencies
deps:
	$(CMD_PREFIX) go mod download

# Update dependencies
deps-update:
	$(CMD_PREFIX) go get -u ./...
	$(CMD_PREFIX) go mod tidy

# Tidy go.mod
tidy:
	$(CMD_PREFIX) go mod tidy

# Format code
fmt:
	$(CMD_PREFIX) go fmt ./...

# Vet code
vet:
	$(CMD_PREFIX) go vet ./...

# Run all validation checks
validate: fmt vet lint test

# Complete CI pipeline
ci: init validate build

# Build Docker image
docker-build:
	docker build -t bian-go .

# Run Docker container
docker-run:
	docker run -p 8080:8080 bian-go

# OpenSpec validation
openspec-validate:
	@echo "OpenSpec validation - implement when openspec CLI is available"

# OpenSpec list
openspec-list:
	@echo "OpenSpec list - implement when openspec CLI is available"

# Interactive shell in container
shell:
	$(DOCKER_RUN) -it $(GO_IMAGE) sh

# Show application logs (when running via compose)
logs:
	docker compose logs -f