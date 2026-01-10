# Change: Add 3 Musketeers Build System

## Why

The project requires a consistent, reproducible build system that works across all development environments (Linux, macOS, Windows) and CI/CD platforms without requiring developers to install specific language toolchains or versions. The 3 Musketeers pattern (Make + Docker/Podman + Compose) solves the "works on my machine" problem by containerizing all build, test, and deployment operations.

Without this standardized build system:
- Developers must install and maintain Go toolchain locally
- Version mismatches cause inconsistent behavior
- CI/CD configuration differs from local environment
- Onboarding requires extensive environment setup
- Build reproducibility is not guaranteed

## What Changes

This change introduces a containerized build automation system following the 3 Musketeers pattern.

### Build Infrastructure Components

- **Makefile** - Unified command interface (25+ commands)
  - Development: init, build, generate, test, run, run-dev
  - Validation: fmt, vet, lint, test-coverage, validate
  - CI/CD: ci pipeline (init + validate + build)
  - Docker: docker-build, docker-run
  - Utilities: clean, shell, logs, deps
  - OpenSpec: openspec-validate, openspec-list

- **docker-compose.yml** - Optional service orchestration (for run/run-dev targets only)
  - `app` service: Production application server
  - `dev` service: Development with hot reload (Air)
  - Note: Not used for build/test commands (overkill for single container tasks)

- **Dockerfile** - Multi-stage production build
  - Stage 1: Build (golang:1.22-alpine, code generation, compilation)
  - Stage 2: Runtime (alpine:3.19, non-root user, health check)
  - Security: Non-root user, minimal runtime dependencies
  - Optimization: Static binary, reduced size (~10MB)

- **Development Tools**
  - `.air.toml`: Hot reload configuration
  - `.dockerignore`: Build context optimization
  - `README.3musketeers.md`: Complete documentation

### Container Runtime Support

- Primary: Direct Docker run (via DOCKER_RUN variable)
- Alternative: Podman (configurable via DOCKER_RUN variable)
- Docker Compose: Only for multi-container scenarios (run, run-dev)
- Build/test commands: Plain `docker run` (simpler, no orchestration needed)
- Official Go 1.22 Alpine base image

### Build System Features

- **Consistency**: Same commands work everywhere (local, CI/CD)
- **Zero installation**: No Go toolchain required locally
- **Version locked**: Go 1.22 in official Alpine image
- **Simplicity**: Direct `docker run` for build/test (no compose overhead)
- **Hot reload**: Development mode with automatic rebuild (compose for orchestration)
- **Security**: Non-root containers, minimal attack surface
- **Portability**: CI/CD agnostic (GitHub Actions, GitLab CI, CircleCI)

## Impact

### Affected Specs
- **NEW**: `specs/devops/` - DevOps automation capability (build, test, deploy)

### Affected Code
- **NEW**: `Makefile` - Make command definitions (142 lines)
- **NEW**: `docker-compose.yml` - Service definitions
- **NEW**: `Dockerfile` - Production build configuration
- **NEW**: `.dockerignore` - Build context exclusions
- **NEW**: `.air.toml` - Hot reload configuration
- **NEW**: `README.3musketeers.md` - Documentation
- **MODIFIED**: `.gitignore` - Added Docker/Air artifacts

### Dependencies
- **Runtime**: Make, Docker (or Podman)
- **Optional**: Docker Compose (only for `make run` and `make run-dev`)
- **Container Images**:
  - `golang:1.22-alpine` - Official Go build environment
  - `alpine:3.19` - Minimal runtime environment
  - `golangci/golangci-lint:v1.55-alpine` - Linting

### Migration Impact
None - This is an additive change. Developers can continue using local Go toolchain or adopt containerized workflow.

### Deployment Impact
- CI/CD pipelines simplified to `make ci`
- Docker required in CI/CD environment
- Build caching improves CI/CD performance
- Multi-platform support (Linux, macOS, Windows)

### Developer Experience Impact
- **Onboarding**: Reduced from hours to minutes (only need Make + Docker)
- **Environment setup**: Zero manual configuration
- **Testing**: Identical local and CI environments
- **Debugging**: `make shell` provides interactive container access

### Future Extensions
- Kubernetes deployment configurations (Helm charts, manifests)
- Additional container registries (GitHub, AWS ECR, Google Artifact Registry)
- Multi-architecture builds (ARM64 support)
- Build performance profiling and optimization
- Integration with OpenSpec workflows (automated proposal validation)
