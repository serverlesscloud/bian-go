# DevOps Capability

## ADDED Requirements

### Requirement: Make Command Interface
The system SHALL provide a Makefile with standardized commands that execute all development, testing, and deployment tasks in containers. All commands use the 3 Musketeers pattern (Make delegates to Docker/Podman).

#### Scenario: Display Available Commands
- **GIVEN** the Makefile exists
- **WHEN** user runs `make help`
- **THEN** the system displays all available commands with descriptions
- **AND** commands are sorted alphabetically
- **AND** descriptions are formatted with color coding

#### Scenario: Build Application
- **GIVEN** the source code exists
- **WHEN** user runs `make build`
- **THEN** the system builds the application in a golang:1.22-alpine container
- **AND** the binary is output to `bin/bian-go`
- **AND** the build uses static linking (CGO_ENABLED=0)

#### Scenario: Run Tests
- **GIVEN** tests exist in the codebase
- **WHEN** user runs `make test`
- **THEN** the system executes all tests in a container
- **AND** test output is verbose (-v flag)
- **AND** tests run without local Go installation

#### Scenario: Run Tests with Coverage
- **GIVEN** tests exist in the codebase
- **WHEN** user runs `make test-coverage`
- **THEN** the system generates coverage.out profile
- **AND** the system generates coverage.html report
- **AND** coverage report is viewable in browser

### Requirement: Container Runtime Abstraction
The system SHALL support both Docker and Podman container runtimes through a configurable DOCKER_RUN variable. Build and test commands use direct container execution (not Docker Compose) to avoid orchestration overhead for single-container tasks.

#### Scenario: Direct Docker Execution
- **GIVEN** Docker is available
- **WHEN** DOCKER_RUN uses `docker run --rm golang:1.22-alpine`
- **THEN** build/test commands execute in standalone Docker containers
- **AND** containers are removed after execution (--rm)
- **AND** no compose orchestration overhead

#### Scenario: Podman Execution
- **GIVEN** Podman is configured
- **WHEN** DOCKER_RUN uses `podman run --rm golang:1.22-alpine`
- **THEN** build/test commands execute in Podman containers
- **AND** behavior is identical to Docker
- **AND** no compose dependency required

#### Scenario: Container Isolation
- **GIVEN** a make command is executed
- **WHEN** the command runs in a container via DOCKER_RUN
- **THEN** the container has isolated filesystem
- **AND** the container is removed after execution
- **AND** no persistent state except mounted volumes

#### Scenario: Design Rationale - No Compose for Build/Test
- **GIVEN** build and test commands need single container execution
- **WHEN** evaluating Docker Compose vs direct docker run
- **THEN** direct docker run is chosen for simplicity
- **AND** Docker Compose is overkill when no orchestration needed
- **AND** Compose only used for multi-container scenarios (app + dev services)

### Requirement: Development Workflow
The system SHALL provide commands for common development workflows including initialization, code generation, formatting, linting, and running the application.

#### Scenario: Project Initialization
- **GIVEN** go.mod and go.sum exist
- **WHEN** user runs `make init`
- **THEN** the system downloads all Go modules
- **AND** the system verifies module checksums
- **AND** modules are cached in Docker volume

#### Scenario: Code Generation
- **GIVEN** source files with go:generate directives exist
- **WHEN** user runs `make generate`
- **THEN** the system executes `go generate ./...` in container
- **AND** generated code is written to host filesystem
- **AND** graphql/generated/ is created if using gqlgen

#### Scenario: Code Formatting
- **GIVEN** Go source files exist
- **WHEN** user runs `make fmt`
- **THEN** the system formats all .go files with gofmt
- **AND** formatting simplifies code (-s flag)
- **AND** formatted files are written to host

#### Scenario: Code Linting
- **GIVEN** Go source files exist
- **WHEN** user runs `make lint`
- **THEN** the system runs golangci-lint in dedicated container
- **AND** linting uses golangci/golangci-lint:v1.55-alpine image
- **AND** lint results are cached for performance

### Requirement: Hot Reload Development Mode
The system SHALL provide a development mode with automatic code reloading when source files change. Hot reload uses the Air tool running in a container.

#### Scenario: Start Development Server
- **GIVEN** application code exists
- **WHEN** user runs `make run-dev`
- **THEN** the system starts dev service with Air
- **AND** application runs on port 8080
- **AND** Air watches for file changes (*.go, *.graphql)

#### Scenario: Auto-Rebuild on Change
- **GIVEN** development server is running
- **WHEN** a .go file is modified
- **THEN** Air automatically rebuilds the application
- **AND** the application restarts with new code
- **AND** changes are visible within 1-2 seconds

#### Scenario: Excluded Directories
- **GIVEN** development server is running with Air
- **WHEN** files change in excluded directories
- **THEN** Air does not trigger rebuild for vendor/, testdata/, openspec/, .git/
- **AND** only relevant changes trigger rebuild

### Requirement: Docker Compose Service Orchestration (Minimal)
The system SHALL define Docker Compose services ONLY for multi-container scenarios (running application with dependencies). Build and test commands use direct docker run for simplicity.

#### Scenario: App Service Definition
- **GIVEN** docker-compose.yml exists
- **WHEN** inspecting app service
- **THEN** service builds from Dockerfile
- **AND** service exposes port 8080
- **AND** service is used by `make run` command
- **AND** service orchestrates production runtime

#### Scenario: Dev Service Definition
- **GIVEN** docker-compose.yml exists
- **WHEN** inspecting dev service
- **THEN** service uses golang:1.22-alpine with Air
- **AND** service provides hot reload capability
- **AND** service is used by `make run-dev` command
- **AND** service mounts source for live editing

#### Scenario: Build Commands Use Direct Docker Run
- **GIVEN** build/test commands need single container
- **WHEN** `make build` or `make test` is executed
- **THEN** commands use DOCKER_RUN variable (direct docker/podman run)
- **AND** no compose overhead or orchestration
- **AND** simpler, faster execution
- **AND** no volume or network complexity

### Requirement: Multi-Stage Production Build
The system SHALL provide a Dockerfile with multi-stage build pattern producing minimal, secure runtime images. Build stage compiles code, runtime stage executes binary.

#### Scenario: Build Stage Execution
- **GIVEN** Dockerfile with build stage
- **WHEN** docker build is executed
- **THEN** build stage uses golang:1.22-alpine
- **AND** build installs git, ca-certificates, tzdata
- **AND** build downloads and verifies Go modules
- **AND** build generates code (go generate)
- **AND** build compiles static binary (CGO_ENABLED=0)

#### Scenario: Runtime Stage Execution
- **GIVEN** Dockerfile with runtime stage
- **WHEN** runtime stage builds from build stage
- **THEN** runtime uses alpine:3.19 base
- **AND** runtime installs only ca-certificates and tzdata
- **AND** runtime creates non-root user (bian:1000)
- **AND** runtime copies binary from build stage
- **AND** runtime switches to non-root user
- **AND** final image size is ~10MB

#### Scenario: Security Hardening
- **GIVEN** production Docker image
- **WHEN** inspecting image configuration
- **THEN** container runs as non-root user (bian)
- **AND** binary is statically linked (no libc dependency)
- **AND** minimal packages installed (reduced attack surface)
- **AND** health check endpoint configured (/health)

### Requirement: CI/CD Integration
The system SHALL provide a unified CI command that executes the complete CI pipeline (init, validate, build) in containers. Pipeline is portable across all CI/CD platforms.

#### Scenario: CI Pipeline Execution
- **GIVEN** source code and tests exist
- **WHEN** user runs `make ci`
- **THEN** the system executes `make init` (download dependencies)
- **AND** the system executes `make validate` (fmt + vet + lint + test)
- **AND** the system executes `make build` (compile binary)
- **AND** pipeline fails fast on first error

#### Scenario: Validation Pipeline
- **GIVEN** source code exists
- **WHEN** user runs `make validate`
- **THEN** the system formats code (make fmt)
- **AND** the system runs go vet (make vet)
- **AND** the system runs linter (make lint)
- **AND** the system runs tests (make test)
- **AND** all steps must pass for success

#### Scenario: GitHub Actions Integration
- **GIVEN** GitHub Actions workflow
- **WHEN** workflow runs `make ci`
- **THEN** workflow behavior matches local execution exactly
- **AND** Docker layer caching improves performance
- **AND** no Go installation required in workflow

#### Scenario: GitLab CI Integration
- **GIVEN** GitLab CI pipeline with docker:dind service
- **WHEN** pipeline runs `make ci`
- **THEN** pipeline behavior matches local execution exactly
- **AND** same Makefile commands work without modification

### Requirement: Build Artifact Management
The system SHALL provide commands to clean build artifacts, Docker resources, and caches. Clean operations are destructive and require explicit invocation.

#### Scenario: Clean Build Artifacts
- **GIVEN** build artifacts exist (bin/, coverage files, generated code)
- **WHEN** user runs `make clean`
- **THEN** the system removes bin/ directory
- **AND** the system removes coverage.out and coverage.html
- **AND** the system removes graphql/generated/
- **AND** the system clears Go caches (cache, testcache, modcache)

#### Scenario: Clean Docker Resources
- **GIVEN** Docker containers and volumes exist
- **WHEN** user runs `make clean-docker`
- **THEN** the system stops all containers
- **AND** the system removes all containers
- **AND** the system removes all volumes
- **AND** the system removes orphan containers

#### Scenario: Selective Cleaning
- **GIVEN** build artifacts and Docker resources exist
- **WHEN** user runs `make clean` (not clean-docker)
- **THEN** only build artifacts are removed
- **AND** Docker volumes persist (faster subsequent builds)

### Requirement: Dependency Management
The system SHALL provide commands to inspect, update, and tidy Go module dependencies. All dependency operations execute in containers.

#### Scenario: Show Dependency Tree
- **GIVEN** go.mod exists with dependencies
- **WHEN** user runs `make deps`
- **THEN** the system displays full module dependency graph
- **AND** output shows direct and indirect dependencies

#### Scenario: Update Dependencies
- **GIVEN** go.mod exists with dependencies
- **WHEN** user runs `make deps-update`
- **THEN** the system updates all dependencies to latest versions
- **AND** the system runs go mod tidy after update
- **AND** go.sum is updated with new checksums

#### Scenario: Tidy Dependencies
- **GIVEN** go.mod and go.sum exist
- **WHEN** user runs `make tidy`
- **THEN** the system removes unused dependencies from go.mod
- **AND** the system adds missing dependencies
- **AND** go.sum is synchronized with go.mod

### Requirement: Development Utilities
The system SHALL provide utility commands for debugging, log viewing, and interactive container access. Utilities assist with troubleshooting and exploration.

#### Scenario: Interactive Shell Access
- **GIVEN** Docker/Podman is running
- **WHEN** user runs `make shell`
- **THEN** the system opens interactive sh shell in golang:1.22-alpine container
- **AND** shell has access to /app working directory
- **AND** shell has access to all Go tools
- **AND** shell is removed on exit

#### Scenario: Application Logs
- **GIVEN** application is running via `make run`
- **WHEN** user runs `make logs`
- **THEN** the system displays application container logs
- **AND** logs stream in real-time (-f flag)
- **AND** logs show stdout and stderr

#### Scenario: Docker Image Building
- **GIVEN** Dockerfile exists
- **WHEN** user runs `make docker-build`
- **THEN** the system builds production Docker image
- **AND** image is tagged as bian-go:latest
- **AND** multi-stage build is used (optimized size)

### Requirement: OpenSpec Integration
The system SHALL provide commands to validate OpenSpec proposals and list changes. OpenSpec commands execute in containers if OpenSpec CLI is available.

#### Scenario: Validate OpenSpec Proposal
- **GIVEN** openspec/ directory with proposals
- **WHEN** user runs `make openspec-validate`
- **THEN** the system validates bootstrap-mvp-foundation proposal
- **AND** validation uses --strict mode
- **AND** validation errors are displayed

#### Scenario: List OpenSpec Changes
- **GIVEN** openspec/ directory with proposals
- **WHEN** user runs `make openspec-list`
- **THEN** the system lists all active changes
- **AND** output shows change IDs and titles

### Requirement: Build System Documentation
The system SHALL provide comprehensive documentation covering prerequisites, usage, workflows, CI/CD integration, and troubleshooting. Documentation is versioned with code.

#### Scenario: Documentation Completeness
- **GIVEN** README.3musketeers.md exists
- **WHEN** reviewing documentation
- **THEN** prerequisites are clearly listed (Make, Docker, Docker Compose)
- **AND** 3 Musketeers principles are explained
- **AND** all make commands are documented
- **AND** CI/CD integration examples are provided
- **AND** troubleshooting section covers common issues

#### Scenario: Quick Start Guide
- **GIVEN** new developer joins project
- **WHEN** following quick start in README.3musketeers.md
- **THEN** developer can build and run application within 5 minutes
- **AND** only Make and Docker installation required
- **AND** no Go toolchain installation needed

### Requirement: Consistency Across Environments
The system SHALL ensure identical build behavior across all environments (Linux, macOS, Windows, CI/CD). Container execution eliminates environment-specific differences.

#### Scenario: Cross-Platform Consistency
- **GIVEN** same codebase on Linux, macOS, and Windows
- **WHEN** running `make build` on each platform
- **THEN** build produces identical binary (same checksum)
- **AND** build uses same Go version (1.22)
- **AND** build process is identical

#### Scenario: Local and CI Parity
- **GIVEN** same commit hash in local and CI environment
- **WHEN** running `make ci` in both environments
- **THEN** both environments produce identical results
- **AND** test outcomes are identical
- **AND** build artifacts are identical
