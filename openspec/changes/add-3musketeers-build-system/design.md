# Design: 3 Musketeers Build System

## Context

The project requires a reproducible build system that works consistently across all development and CI/CD environments. The traditional 3 Musketeers pattern uses Make + Docker + Docker Compose, but we need to evaluate whether Docker Compose is necessary for all operations or if it adds unnecessary complexity.

**Key Constraints:**
- Must work on Linux, macOS, Windows
- Must work in CI/CD without special configuration
- Should minimize dependencies
- Should be simple to understand and maintain
- No local Go installation required

**Stakeholders:**
- Developers working locally
- CI/CD pipelines
- New contributors onboarding to project

## Goals / Non-Goals

### Goals
- Consistent build behavior across all environments
- Minimal required dependencies (Make + Docker/Podman)
- Simple, understandable build process
- Fast build/test execution
- No Go toolchain installation required

### Non-Goals
- Complex multi-container orchestration (not needed for this project)
- Advanced networking between build containers
- Persistent data volumes for build artifacts (ephemeral is fine)
- Supporting non-containerized workflows (this is container-first)

## Decisions

### Decision 1: Direct Docker Run vs Docker Compose for Build/Test

**Choice:** Use direct `docker run` (via DOCKER_RUN variable) for build/test commands. Use Docker Compose only for multi-container scenarios (run, run-dev).

**Rationale:**

**Problem:**
Traditional 3 Musketeers pattern always uses `docker compose run --rm service-name command`. This requires:
1. Docker Compose installation
2. docker-compose.yml with service definitions
3. Named volumes for caching
4. Network configuration
5. More moving parts to understand

For single-container tasks (build, test, lint), this orchestration is **overkill**.

**Analysis:**

| Aspect | Docker Compose | Direct Docker Run | Winner |
|--------|---------------|-------------------|---------|
| Dependencies | Docker + Compose | Docker only | Direct run |
| Complexity | High (YAML config, services, networks) | Low (just run command) | Direct run |
| Performance | Compose overhead | No overhead | Direct run |
| Caching | Named volumes | Host mount | Equal |
| Multi-container | Native support | Manual coordination | Compose |
| Single container | Overkill | Perfect fit | Direct run |

**When Docker Compose IS valuable:**
- `make run`: Running production app container (may add DB, Redis later)
- `make run-dev`: Running dev container with hot reload

**When Docker Compose is OVERKILL:**
- `make build`: Single container compiles code
- `make test`: Single container runs tests
- `make lint`: Single container runs linter
- `make fmt`: Single container formats code

**Implementation:**
```makefile
# Simple, direct execution
DOCKER_RUN = docker run --rm -v $(PWD):/app -w /app golang:1.22-alpine

build:
    $(DOCKER_RUN) go build -o bin/bian-go ./cmd/server

test:
    $(DOCKER_RUN) go test ./...

# Multi-container scenarios still use compose
run:
    docker compose up app

run-dev:
    docker compose up dev
```

**Alternatives Considered:**

1. **Use Docker Compose for Everything**
   - ✅ Consistent pattern (always compose)
   - ❌ Requires Compose installation
   - ❌ More complex (docker-compose.yml must define all services)
   - ❌ Slower (compose startup overhead)
   - ❌ Harder to understand for newcomers

2. **Use Direct Docker Run for Everything**
   - ✅ Minimal dependencies
   - ✅ Simple, easy to understand
   - ❌ Multi-container scenarios become complex
   - ❌ Need manual networking for app + db scenarios
   - ❌ No easy way to run multiple services together

3. **Hybrid Approach (Chosen)**
   - ✅ Direct run for simple tasks (build, test, lint)
   - ✅ Compose for multi-container (run, run-dev)
   - ✅ Minimal dependencies (Compose optional)
   - ✅ Simple where it can be, powerful where needed
   - ➖ Two patterns to learn (acceptable trade-off)

**Trade-offs:**
- ✅ Simpler: Direct `docker run` is easier to understand than Compose
- ✅ Faster: No compose startup overhead for build/test
- ✅ Fewer deps: Compose only needed for `make run` targets
- ✅ More portable: Works with Podman seamlessly (podman run)
- ➖ Two patterns: Direct run AND compose (but clear separation)
- ✅ Flexibility: Can easily switch between Docker and Podman

### Decision 2: DOCKER_RUN Variable Abstraction

**Choice:** Use a DOCKER_RUN Make variable that can be configured for Docker or Podman.

**Rationale:**
- Single point of configuration for container runtime
- Easy to switch between Docker and Podman
- No changes to individual make targets
- CI/CD can override for different container runtimes

**Implementation:**
```makefile
# Default to docker, can be overridden
DOCKER_RUN = docker run --rm -v $(PWD):/app -w /app golang:1.22-alpine

# Or use podman
# DOCKER_RUN = podman run --rm -v $(PWD):/app -w /app golang:1.22-alpine

# All commands use this variable
build:
    $(DOCKER_RUN) go build ...

test:
    $(DOCKER_RUN) go test ...
```

**Trade-offs:**
- ✅ Single source of truth for container runtime
- ✅ Easy runtime switching
- ✅ Clear what's containerized (uses DOCKER_RUN)
- ➖ Variable must be maintained consistently

### Decision 3: Minimal docker-compose.yml

**Choice:** docker-compose.yml defines only `app` and `dev` services. No `go`, `lint`, or build-related services.

**Rationale:**
- Build services in compose add complexity without value
- Direct `docker run` handles build/test better
- Compose focused on runtime orchestration only
- Easier to understand (app runtime vs build tools)

**docker-compose.yml structure:**
```yaml
services:
  # Production application runtime
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - GRAPHQL_PLAYGROUND=false

  # Development with hot reload
  dev:
    image: golang:1.22-alpine
    command: air
    ports:
      - "8080:8080"
    volumes:
      - .:/app
    working_dir: /app
```

No need for:
- `go` service (handled by DOCKER_RUN)
- `lint` service (handled by DOCKER_RUN)
- Complex networks (not needed)
- Named volumes (host mounts suffice)

**Trade-offs:**
- ✅ Simpler compose file
- ✅ Clear purpose (runtime only)
- ✅ Less to maintain
- ❌ Can't use `docker compose run go ...` (but we don't want to)

## Risks / Trade-offs

### Risk 1: Learning Curve - Two Patterns
**Risk:** Developers must learn both direct docker run and compose

**Mitigation:**
- Clear documentation explaining when to use each
- Pattern is simple: build/test = direct, run = compose
- README.3musketeers.md documents both approaches
- Makefile hides complexity (developers just run `make build`)

**Impact:** Low - patterns are intuitive

### Risk 2: Portability Between Docker and Podman
**Risk:** Direct docker run commands may not be fully compatible with Podman

**Mitigation:**
- Use DOCKER_RUN variable for abstraction
- Test with both Docker and Podman
- Keep command flags minimal and compatible
- Document any known incompatibilities

**Impact:** Low - Podman aims for Docker CLI compatibility

### Risk 3: Compose Optional Dependency
**Risk:** Developers confused about when Compose is needed

**Mitigation:**
- README clearly states: "Docker Compose optional (only for make run)"
- Makefile help shows which commands need Compose
- CI/CD examples show working without Compose for build/test

**Impact:** Low - clear documentation resolves this

## Migration Plan

**N/A** - This is the initial implementation.

**For Future Changes:**
- If multi-container scenarios increase (add DB, Redis, etc.), expand docker-compose.yml
- If build complexity increases, consider adding build services to compose
- Monitor developer feedback - if hybrid approach is confusing, consider consolidating

## Open Questions

1. **Volume Caching Strategy**
   - Direct docker run doesn't use named volumes by default
   - Should we add volume caching for go-modules?
   - Decision: Start without, add if build speed becomes issue

2. **Network Isolation**
   - Direct docker run uses default bridge network
   - Should we create custom network for isolation?
   - Decision: Not needed for single containers, only for multi-container scenarios

3. **Windows Path Handling**
   - $(PWD) may need adjustment for Windows
   - Should we use Docker's / Compose's path handling?
   - Decision: Document in README, provide Windows-specific examples if needed

4. **CI/CD Optimization**
   - Should CI/CD use different caching strategy?
   - Docker layer caching vs volume caching?
   - Decision: Start simple, optimize later based on CI times
