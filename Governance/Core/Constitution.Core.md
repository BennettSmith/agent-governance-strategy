This document defines non‑negotiable rules for governance and collaboration, shared across profiles.

## Decision-making (core)

- Architectural decisions are captured explicitly when they affect system shape or long-term constraints.
- All work is planned before implementation and documented in markdown.
- Decisions favor clarity, testability, and long-term evolvability over short-term convenience.

## Default preferences (core, not hard rules)

- Explicit boundary types are preferred over leaking internal models.
- Prefer dependency inversion over direct instantiation.
- Prefer constructor-based injection over global state.
- Prefer explicit data flow over implicit framework behavior.
- Prefer small, reviewable changes over large refactors.

## Legacy refactoring protocol (core)

Legacy code may not yet follow the target architecture. Refactoring legacy code is permitted, but must be done in a way that is **orderly, reversible, and separated from behavior change**.

### Definitions

- **Refactor**: An internal restructure that preserves externally observable behavior at defined boundaries (persistence contracts, API interactions, and other profile-defined external boundaries).
- **Behavior change**: Any change that alters externally observable behavior at those boundaries. Behavior changes are features/fixes, not refactors.

### Separation rules

- Refactors and behavior changes must not be mixed in the same branch or pull request.
- Refactor branches should use the `refactor/` prefix; behavior-change branches use `feat/` or `fix/`.

### Required workflow (repeat for each step)

- **Lock behavior**: Add characterization/contract/golden-master tests at the boundary you intend to preserve.
- **Introduce a seam**: Add an adapter/facade/port or boundary that lets you change internals without spreading violations.
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

## Agent rules (core)

- All work is performed in feature branches.
- Work proceeds in small, incremental steps aligned to a written plan.
- Plans must be stored in-repo at `Docs/Plans/<branch-name>.md` and kept up to date during development.
- Test-Driven Development is the default mode of work.
- Quality gates define completion, not code presence.
- Agents must stop after completing each planned task and request manual acceptance.
- Agents may only mark a task complete and create a checkpoint commit after:
  - explicit human approval
  - all quality gates passing
- Agents must never merge changes.
- Humans are responsible for acceptance, direction, and merge decisions.

## Conflict resolution (core)

- `Non-Negotiables.md` overrides all other governing documents.
- `Architecture.md` overrides `Constitution.md` on matters of system shape.
- `Constitution.md` guides behavior when other documents are silent.
- In case of ambiguity, stop work and ask for human direction.

