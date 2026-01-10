# Implementation Tasks

## 1. Makefile Development
- [x] 1.1 Create Makefile with help target
- [x] 1.2 Define DOCKER_RUN variable (direct docker/podman run, not compose)
  - Rationale: Compose is overkill for single-container build/test tasks
  - Use plain `docker run` or `podman run` for simplicity
- [x] 1.3 Implement init target (go mod download/verify)
- [x] 1.4 Implement build target (compile binary)
- [x] 1.5 Implement generate target (code generation)
- [x] 1.6 Implement test targets (test, test-coverage, test-short)
- [x] 1.7 Implement lint, fmt, vet targets
- [x] 1.8 Implement run and run-dev targets
- [x] 1.9 Implement clean targets (clean, clean-docker)
- [x] 1.10 Implement dependency management targets (deps, deps-update, tidy)
- [x] 1.11 Implement validation and CI targets (validate, ci)
- [x] 1.12 Implement Docker targets (docker-build, docker-run)
- [x] 1.13 Implement OpenSpec targets (openspec-validate, openspec-list)
- [x] 1.14 Implement utility targets (shell, logs)
- [x] 1.15 Set help as default target (.DEFAULT_GOAL)

## 2. Docker Compose Configuration (Minimal - for run targets only)
- [x] 2.1 Create docker-compose.yml (optional, only for multi-container scenarios)
- [x] 2.2 Define app service (production application)
  - Build from Dockerfile
  - Port mapping 8080:8080
  - Environment configuration
  - Used by: `make run`
- [x] 2.3 Define dev service (hot reload development)
  - golang:1.22-alpine with Air
  - Port mapping 8080:8080
  - Volume mounts for live editing
  - Used by: `make run-dev`
- [x] 2.4 Note: Build/test commands use direct `docker run` (via DOCKER_RUN variable)
  - No compose orchestration needed for single-container tasks
  - Simpler, faster, less overhead

## 3. Production Dockerfile
- [x] 3.1 Create multi-stage Dockerfile
- [x] 3.2 Stage 1: Build stage (golang:1.22-alpine)
  - Install build dependencies (git, ca-certificates, tzdata)
  - Copy go.mod and go.sum
  - Download and verify dependencies
  - Copy source code
  - Run code generation (go generate)
  - Build static binary (CGO_ENABLED=0, -ldflags="-w -s")
- [x] 3.3 Stage 2: Runtime stage (alpine:3.19)
  - Install runtime dependencies (ca-certificates, tzdata)
  - Create non-root user (bian:bian, uid/gid 1000)
  - Copy binary from build stage
  - Set ownership and permissions
  - Switch to non-root user
  - Expose port 8080
  - Configure health check (wget /health)
  - Set environment defaults (PORT, GRAPHQL_PLAYGROUND)
  - Define CMD to run binary

## 4. Development Tools
- [x] 4.1 Create .air.toml for hot reload
  - Configure build command and output
  - Set watched file extensions (go, graphql)
  - Exclude directories (vendor, testdata, openspec, .git)
  - Configure rebuild delay and kill behavior
- [x] 4.2 Create .dockerignore
  - Exclude .git, CI/CD configs, documentation
  - Exclude editor/IDE files
  - Exclude openspec/ (not needed in production)
  - Exclude build artifacts and test outputs
  - Exclude temporary files

## 5. Documentation
- [x] 5.1 Create README.3musketeers.md
  - Document prerequisites (Make, Docker, Docker Compose)
  - Explain 3 Musketeers principles (Consistency, Control, Confidence)
  - List all available commands with descriptions
  - Provide quick start guide
  - Document CI/CD integration (GitHub Actions, GitLab CI)
  - Add development workflow examples
  - Include troubleshooting section
  - Document architecture and benefits

## 6. Git Configuration
- [x] 6.1 Update .gitignore
  - Add tmp/ (Air hot reload)
  - Add build-errors.log (Air logs)
  - Add docker-compose.override.yml (local overrides)

## 7. Testing and Validation
- [x] 7.1 Test make help command (display all targets)
- [x] 7.2 Verify Makefile syntax (no syntax errors)
- [x] 7.3 Test docker-compose.yml syntax (docker compose config)
- [x] 7.4 Test Dockerfile build (when code exists)
- [x] 7.5 Verify .dockerignore excludes correct files
- [x] 7.6 Test Podman compatibility (DOCKER_RUN variable)
- [x] 7.7 Validate documentation completeness

## 8. Integration
- [x] 8.1 Ensure Makefile aligns with OpenSpec tasks
  - Phase 1 (Setup): make init
  - Phases 2-10 (Implementation): make test, make generate, make build
  - Phase 11 (Testing): make test-coverage
  - Phase 13 (Build): make docker-build
  - Phase 14 (Validation): make validate
- [x] 8.2 Verify all commands work without local Go installation
- [x] 8.3 Test volume caching improves build performance
- [x] 8.4 Verify hot reload works in dev mode

## 9. Final Verification
- [x] 9.1 Run make help successfully
- [x] 9.2 Verify all .PHONY targets defined
- [x] 9.3 Confirm Podman support via DOCKER_RUN variable
- [x] 9.4 Validate docker-compose.yml with docker compose config
- [x] 9.5 Check documentation is complete and accurate
- [x] 9.6 Ensure .gitignore excludes all generated artifacts
