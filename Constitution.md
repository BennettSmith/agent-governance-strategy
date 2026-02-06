# Constitution

This document defines non‑negotiable rules for system design and agent collaboration.

## Architectural Identity

- Offline-first, client-centric iOS application
- Clean Architecture with Domain-Driven Design as the primary modeling approach
- Use-case–driven application core
- API-first integration with backend services defined by OpenAPI
- Dependency-injected system with an explicit composition root
- UI and infrastructure are replaceable details, not the core of the system

## Decision-Making

- Architectural decisions are captured explicitly when they affect system shape or long-term constraints
- All work is planned before implementation and documented in markdown
- Decisions favor clarity, testability, and long-term evolvability over short-term convenience
- When uncertain, choose the option that preserves domain independence, offline-first behavior, explicit boundaries, and composability

## Default Preferences (Not Hard Rules)

- Domain logic is expressed in ubiquitous language
- Use cases are the primary unit of business behavior
- Explicit boundary types are preferred over leaking internal models
- Prefer dependency inversion over direct instantiation
- Prefer constructor-based injection over global state
- Prefer explicit data flow over implicit framework behavior
- Prefer small, reviewable changes over large refactors

## Legacy refactoring protocol

Legacy code may not yet follow the target architecture. Refactoring legacy code is permitted, but must be done in a way that is **orderly, reversible, and separated from behavior change**.

### Definitions

- **Refactor**: An internal restructure that preserves externally observable behavior at defined boundaries (use case input/output boundaries, persistence contracts, API interactions, and user-visible UI behavior).
- **Behavior change**: Any change that alters externally observable behavior at those boundaries. Behavior changes are features/fixes, not refactors.

### Separation rules

- Refactors and behavior changes must not be mixed in the same branch or pull request.
- Refactor branches should use the `refactor/` prefix; behavior-change branches use `feat/` or `fix/`.

### Required workflow (repeat for each step)

- **Lock behavior**: Add characterization/contract/golden-master tests at the boundary you intend to preserve.
- **Introduce a seam**: Add an adapter/facade/port or use-case boundary that lets you change internals without spreading violations.
- **Refactor behind the seam**: Make a small internal change that keeps prior behavior intact.
- **Run quality gates**: All required checks must pass for each step.
- **Checkpoint**: Record progress in the branch plan; keep steps small enough that a revert is safe and cheap.

### Approved patterns

See `Docs/Refactoring/Legacy-Refactoring-Playbook.md` for the approved patterns and checklists, including:

- Strangler Fig
- Branch by Abstraction
- Parallel Change / Expand–Contract
- “Seam first” refactoring
- Sprout method/class

## Fundamental Rules

### 1. Use Cases Are the Source of Truth

All application behavior **must** be represented by a documented use case within a bounded context.

### 2. Bounded Context Ownership

Every use case belongs to **exactly one bounded context**.

### 3. Explicit Boundaries

Use cases accept explicit input/output boundary types and never expose domain entities.

Domain models are **internal to a bounded context** and must not cross context boundaries.

Mapping between boundary types and domain models is **explicit**, performed at the boundary, and owned by the bounded context.

### 4. Dependency Direction

Dependencies must always point inward.

### 5. No Domain Logic Outside Use Cases

Domain rules and invariants live inside use cases or their owned entities.

### 6. Offline‑First by Default

Use cases must function without network connectivity unless explicitly documented.

## Agent Rules

- All work is performed in feature branches
- Work proceeds in small, incremental steps aligned to a written plan
- Plans are executed one task (todo) at a time
- Plans must be stored in-repo at `Docs/Plans/<branch-name>.md` and kept up to date during development
- Feature branch names must follow: `<type>/<area>-<short-slug>` (lowercase, hyphenated)
  - Allowed `<type>` values: `feat`, `fix`, `docs`, `chore`, `refactor`, `spike`
  - `<area>` is usually a bounded context slug (e.g., `identity`)
  - If the change is not related to a bounded context (e.g., composition root / app wiring / routing), use `app` as the `<area>` segment
  - Other `<area>` values may be used when appropriate, but must be clear, stable, and documented in the plan
- Test-Driven Development is the default mode of work
- Quality gates define completion, not code presence
- Agents must stop after completing each planned task and request manual acceptance
- Agents may only mark a task complete and create a checkpoint commit after:
  - Explicit human approval
  - All quality gates passing
- If an agent is unsure how to proceed, it must stop and ask the human
- Agents may propose options and ask for preferences
- Agents must not make arbitrary or irreversible decisions independently
- Agents may create branches, commits, and pull/merge requests
- Agents must never merge changes
- Humans are responsible for acceptance, direction, and merge decisions
- Architectural boundaries and dependency direction must be respected at all times

## Bounded Contexts

- The system is organized around **bounded contexts** as defined by Domain-Driven Design.
- A bounded context defines a coherent domain model with:
  - a consistent ubiquitous language
  - well-defined responsibilities and invariants
- Domain entities, value types, aggregates, and domain events are **owned by a single bounded context**.
- Domain models are **internal to their bounded context** and must not be exposed outside of it.
- Interaction with a bounded context occurs only through application use cases.
- Communication across bounded contexts must be explicit and must not rely on shared domain models.

## Conflict Resolution

- `Non-Negotiables.md` overrides all other documents
- `Architecture.md` overrides `Constitution.md` on matters of system shape
- `Constitution.md` guides behavior when other documents are silent
- In case of ambiguity, stop work and ask for human direction
