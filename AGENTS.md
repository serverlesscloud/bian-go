<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

# Agents Configuration

## Dev Agent

**Primary References:**
- `@/agents/ARCHITECTURE.md` - Technology stack, plugin architecture, design decisions
- `@/agents/DEV.md` - Development guidelines
- `@/openspec/AGENTS.md` - Only for spec-driven development (proposals, changes, archiving)

## Test Agent

**Primary References:**
- `@/agents/TESTING.md` - Complete testing guide, patterns, and best practices
- `@/openspec/AGENTS.md` - Only for spec-driven test development

## Document Agent

**Primary References:**
- `@/agents/ARCHITECTURE.md` - System architecture and design decisions
- `@/agents/TESTING.md` - Testing documentation standards
- `@/openspec/AGENTS.md` - Only for spec documentation

**Documentation Standards:**
- Update `README.md` for major feature changes
- Maintain OpenSpec specs as source of truth for requirements
- Document API changes in GraphQL schema
- Keep deployment guides in their relevate spec or change folder