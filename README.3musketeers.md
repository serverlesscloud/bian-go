# 3 Musketeers Build System

This project uses the [3 Musketeers](https://3musketeers.pages.dev/) pattern for consistent, reproducible builds across all environments.

## Prerequisites

You only need these three tools installed:
- **Make** - Task runner
- **Docker** - Container runtime
- **Docker Compose** - Container orchestration

**No Go installation required!** All Go commands run in Docker containers.

## Quick Start

```bash
# Show all available commands
make help

# Initialize project (download dependencies)
make init

# Build the application
make build

# Run tests
make test

# Run the application
make run

# Run with hot reload for development
make run-dev
```

## Core Principles

### 1. Consistency
Run the same commands everywhere:
- ✅ Linux, macOS, Windows
- ✅ Local development
- ✅ CI/CD (GitHub Actions, GitLab CI, CircleCI)

### 2. Control
- Go version locked to `golang:1.22-alpine`
- All dependencies defined in `go.mod`
- Configuration in version control

### 3. Confidence
- Test locally with exact CI environment
- If it works locally, it works in CI

## Available Commands

### Development
```bash
make init          # Download dependencies
make build         # Build application binary
make generate      # Generate code (GraphQL, mocks)
make test          # Run all tests
make test-coverage # Run tests with coverage report
make lint          # Run linter (golangci-lint)
make fmt           # Format code
make vet           # Run go vet
make run           # Run application
make run-dev       # Run with hot reload
```

### Validation
```bash
make validate      # Run fmt + vet + lint + test
make ci            # Full CI pipeline (init + validate + build)
```

### Dependencies
```bash
make deps          # Show dependency tree
make deps-update   # Update all dependencies
make tidy          # Tidy go.mod and go.sum
```

### Docker
```bash
make docker-build  # Build production Docker image
make docker-run    # Run production Docker image
```

### Utilities
```bash
make clean         # Clean build artifacts
make clean-docker  # Remove Docker containers and volumes
make shell         # Open shell in Go container
make logs          # Show application logs
```

### OpenSpec
```bash
make openspec-validate  # Validate OpenSpec proposal
make openspec-list      # List OpenSpec changes
```

## How It Works

### Makefile (Interface)
The Makefile provides a consistent interface for all tasks. Every command delegates to Docker:

```makefile
COMPOSE_RUN = docker compose run --rm go

build:
	$(COMPOSE_RUN) go build -o bin/bian-go ./cmd/server
```

### Docker Compose (Orchestration)
`docker-compose.yml` defines services:

- **go** - Primary service for build/test tasks
- **app** - Application server
- **dev** - Development server with hot reload
- **golangci-lint** - Linter service

All services use the official `golang:1.22-alpine` image.

### Docker (Environment)
Every command runs in a Docker container with:
- Consistent Go version (1.22)
- Isolated environment
- Cached dependencies (Docker volumes)
- Same behavior everywhere

## Caching

Docker volumes cache expensive operations:

```yaml
volumes:
  - go-modules:/go/pkg/mod           # Go module cache
  - go-build-cache:/root/.cache      # Build cache
```

This makes subsequent builds much faster.

## CI/CD Integration

### GitHub Actions Example
```yaml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run CI pipeline
        run: make ci
```

### GitLab CI Example
```yaml
image: docker:latest
services:
  - docker:dind

test:
  script:
    - make ci
```

The same `make ci` command works everywhere!

## Development Workflow

### First Time Setup
```bash
# Clone repository
git clone https://github.com/serverlesscloud/bian-go.git
cd bian-go

# Initialize project (downloads dependencies)
make init

# Verify setup
make validate
```

### Daily Development
```bash
# Start development server with hot reload
make run-dev

# In another terminal, run tests on save
make test

# Format and lint code
make fmt
make lint

# Build production binary
make build
```

### Before Commit
```bash
# Run all checks
make validate

# If checks pass, commit
git add .
git commit -m "Your message"
```

## Troubleshooting

### Command not found: make
- **macOS**: `brew install make`
- **Ubuntu/Debian**: `sudo apt-get install make`
- **Windows**: Use WSL2 or install via Chocolatey

### Docker permission denied
```bash
# Linux: Add user to docker group
sudo usermod -aG docker $USER
newgrp docker
```

### Cache issues
```bash
# Clear Docker caches
make clean-docker

# Rebuild from scratch
docker compose build --no-cache
```

### Slow builds on Windows
Enable WSL2 backend in Docker Desktop for better performance.

## Benefits

### For Developers
- ✅ No Go installation needed
- ✅ No version conflicts
- ✅ Reproducible builds
- ✅ Fast onboarding

### For Teams
- ✅ Same environment for everyone
- ✅ No "works on my machine" issues
- ✅ Easy CI/CD setup
- ✅ Version-controlled tooling

### For CI/CD
- ✅ Identical to local environment
- ✅ No special configuration
- ✅ Fast with Docker layer caching
- ✅ Portable across CI providers

## Architecture

```
┌─────────────────────────────────────────┐
│           Makefile (Interface)          │
│  Developer runs: make build, make test  │
└─────────────────┬───────────────────────┘
                  │
                  ↓
┌─────────────────────────────────────────┐
│      Docker Compose (Orchestration)     │
│  Defines services: go, app, dev, lint   │
└─────────────────┬───────────────────────┘
                  │
                  ↓
┌─────────────────────────────────────────┐
│       Docker Containers (Execution)     │
│  golang:1.22-alpine - All commands run  │
│  in isolated, consistent environment    │
└─────────────────────────────────────────┘
```

## Resources

- [3 Musketeers Documentation](https://3musketeers.pages.dev/)
- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Reference](https://docs.docker.com/compose/)
- [Official Go Docker Images](https://hub.docker.com/_/golang)
