# 3 Musketeers Makefile for BIAN-Go
# Requires: Make, Docker, Docker Compose
# All commands run in Docker containers for consistency

.PHONY: help
help: ## Show this help message
	@echo "BIAN-Go - 3 Musketeers Build System"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Docker Compose command (all targets use this)
DOCKER_RUN = podman run --rm go

.PHONY: init
init: ## Initialize project (download dependencies)
	$(DOCKER_RUN) go mod download
	$(DOCKER_RUN) go mod verify

.PHONY: build
build: ## Build the application binary
	@echo "Building BIAN-Go..."
	$(DOCKER_RUN) go build -o bin/bian-go ./cmd/server
	@echo "Build completed: bin/bian-go"

.PHONY: generate
generate: ## Generate code (GraphQL, mocks, etc.)
	@echo "Generating code..."
	$(DOCKER_RUN) go generate ./...
	@echo "Code generation completed."

.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	$(DOCKER_RUN) go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	$(DOCKER_RUN) go test -v -coverprofile=coverage.out -covermode=atomic ./...
	$(DOCKER_RUN) go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: test-short
test-short: ## Run short tests only
	$(DOCKER_RUN) go test -short ./...

.PHONY: lint
lint: ## Run linter (golangci-lint)
	@echo "Running linter..."
	docker compose run --rm golangci-lint golangci-lint run ./...

.PHONY: fmt
fmt: ## Format code with gofmt
	@echo "Formatting code..."
	$(DOCKER_RUN) gofmt -s -w .

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(DOCKER_RUN) go vet ./...

.PHONY: tidy
tidy: ## Tidy go.mod and go.sum
	@echo "Tidying dependencies..."
	$(DOCKER_RUN) go mod tidy

.PHONY: run
run: ## Run the example application
	@echo "Starting BIAN-Go server..."
	docker compose up --build app

.PHONY: run-dev
run-dev: ## Run with hot reload (if air is configured)
	@echo "Starting development server with hot reload..."
	docker compose up --build dev

.PHONY: clean
clean: ## Clean build artifacts and caches
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf coverage.out coverage.html
	rm -rf graphql/generated/
	$(DOCKER_RUN) go clean -cache -testcache -modcache
	@echo "Clean completed."

.PHONY: clean-docker
clean-docker: ## Remove all Docker containers and volumes
	@echo "Cleaning Docker resources..."
	docker compose down -v --remove-orphans
	@echo "Docker cleanup completed."

.PHONY: deps
deps: ## Show dependency tree
	$(DOCKER_RUN) go mod graph

.PHONY: deps-update
deps-update: ## Update all dependencies
	@echo "Updating dependencies..."
	$(DOCKER_RUN) go get -u ./...
	$(DOCKER_RUN) go mod tidy

.PHONY: validate
validate: fmt vet lint test ## Run all validation checks (fmt, vet, lint, test)
	@echo "All validation checks passed!"

.PHONY: ci
ci: init validate build ## Run CI pipeline (init, validate, build)
	@echo "CI pipeline completed successfully!"

.PHONY: docker-build
docker-build: ## Build production Docker image
	@echo "Building production Docker image..."
	docker build -t bian-go:latest -f Dockerfile .
	@echo "Docker image built: bian-go:latest"

.PHONY: docker-run
docker-run: docker-build ## Run production Docker image
	@echo "Running production Docker image..."
	docker run --rm -p 8080:8080 bian-go:latest

.PHONY: openspec-validate
openspec-validate: ## Validate OpenSpec proposal
	@echo "Validating OpenSpec proposal..."
	$(DOCKER_RUN) openspec validate bootstrap-mvp-foundation --strict

.PHONY: openspec-list
openspec-list: ## List OpenSpec changes
	$(DOCKER_RUN) openspec list

.PHONY: shell
shell: ## Open a shell in the Go container
	podman run --rm go sh

.PHONY: logs
logs: ## Show application logs
	podman logs -f app

# Default target
.DEFAULT_GOAL := help
