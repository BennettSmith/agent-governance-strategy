# Legacy refactoring playbook

This playbook defines how we refactor legacy code in a way that is **orderly**, **small-step**, and **reversible**, while converging on the target architecture.

This document is referenced by:

- `Non‑Negotiables.md` (enforceable constraints)
- `Constitution.md` (required protocol and separation rules)
- `Architecture.md` (legacy modernization guidance)

## Goals

- Improve internal design without changing externally observable behavior (unless explicitly doing a behavior-change task).
- Move legacy code toward Clean Architecture boundaries and the use-case–driven core.
- Reduce risk by ensuring every step has a safety net and a rollback path.

## Core concepts

### What counts as “externally observable behavior”

At minimum, the following are “externally observable” boundaries for the purposes of refactoring:

- **Use case boundaries**: input/output boundary types, error mapping, idempotency guarantees, side effects
- **Persistence boundaries**: schema/serialization formats, migration behavior, query semantics, ordering, uniqueness rules
- **API boundaries**: requests/responses via the generated OpenAPI client, retry/backoff behavior, status-code handling
- **User-visible UI behavior**: visible states, copy, navigation outcomes, timing constraints users can observe

If any of the above changes, the work is a **behavior change** (feature/fix), not a refactor.

### Safety net

A safety net is the mechanism that makes refactoring safe. Examples:

- **Characterization tests**: capture today’s behavior (even if it is “weird”) before changing structure.
- **Contract tests**: assert behavior at an explicit boundary (e.g., a port protocol, a use case output mapping).
- **Golden-master tests**: snapshot outputs for a variety of inputs and compare for regressions.

Guideline: put the safety net at the **highest stable boundary** you can (often use-case boundaries or ports), not deep inside unstable internals.

### Seam

A seam is a place you can change behavior/structure without editing callers.

Common seams in this architecture:

- **Use case boundary** (controllers call `execute`; use cases map boundary types to/from domain)
- **Ports/adapters** (protocols for persistence, networking, clocks, connectivity)
- **Facades** in the interface/infrastructure layer that wrap legacy APIs

“Seam first” is often the difference between a safe refactor and a sprawling rewrite.

## Required workflow (the loop)

Use this loop for each refactor step:

1. **Define the boundary** you promise not to change.
2. **Add the safety net** that will fail if behavior changes at that boundary.
3. **Introduce or strengthen a seam** so changes remain localized.
4. **Refactor behind the seam** in the smallest step that provides value.
5. **Run quality gates** and ensure the diff is trivially revertible.
6. **Checkpoint**: record what changed, what remained stable, and what’s next.

If you cannot name your rollback plan for the current step, the step is too large.

## Approved migration/refactoring patterns

### Strangler Fig (recommended for large legacy surfaces)

Use when: legacy component is large and risky to rewrite; you can route calls through a boundary.

How:

- Wrap the legacy behavior behind a seam (facade/adapter/use case boundary).
- Add a “new path” in parallel for a small slice.
- Route a small subset of calls to the new path.
- Expand slice-by-slice until legacy path is unused, then remove.

Rollback: route back to legacy path.

### Branch by Abstraction

Use when: you need to replace an implementation but can’t change all callers at once.

How:

- Introduce an abstraction (protocol/interface) at the boundary.
- Provide two implementations: legacy and new.
- Move callers to the abstraction (no behavior change).
- Switch implementations (guarded by a flag if needed).

Rollback: switch back to legacy implementation.

### Parallel Change / Expand–Contract

Use when: you must change a contract (data format, API shape, persistence schema) safely.

How:

- **Expand**: introduce the new format/field/path alongside the old.
- Run both in parallel (write both; read old with fallback to new).
- Migrate consumers.
- **Contract**: remove the old format/field/path after stabilization.

Rollback: keep the old path functioning; delay contract.

### Sprout method / sprout class

Use when: adding new structure inside messy code without untangling everything.

How:

- Identify a small, cohesive behavior.
- Add a new method/class that implements it with good design.
- Call it from legacy code, keeping the integration point small.

Rollback: revert the call site; remove sprout code.

### “Seam first” wrapper

Use when: legacy code is hard to test or tightly coupled.

How:

- Add a wrapper that exposes a stable, testable boundary.
- Add characterization tests against the wrapper.
- Refactor the internals behind it.

Rollback: wrapper continues to delegate to legacy behavior.

## How to keep steps small and reversible

Use one of these tactics when a refactor is “too big”:

- **Split by boundary**: do use-case boundary cleanup first, internals later.
- **Split by direction**: first remove dependencies (introduce ports), then reorganize code.
- **Parallelize**: add a new implementation without removing the old one yet.
- **Flags** (sparingly): use only when parallel change is otherwise impractical; remove flags as part of the plan.

Preferred rollback strategy order:

1. Simple revert (small diff)
2. Switch routing back to legacy path (strangler/branch-by-abstraction)
3. Flag off the new path (temporary)

## Refactor PR checklist (required)

- [ ] The PR is **refactor-only** (no behavior changes), or explicitly labeled as behavior change and split accordingly.
- [ ] The preserved boundary is stated (use case I/O, port contract, persistence contract, UI behavior).
- [ ] A safety net exists and demonstrates behavior preservation (characterization/contract/golden-master tests).
- [ ] The step is reversible (revert/route-back/flag-off) and the rollback is documented.
- [ ] New code introduced follows the target architecture; legacy violations do not spread.
- [ ] Quality gates pass.

## Examples (what “good” looks like)

### Example: wrapping a legacy persistence API

- Boundary: use case output mapping and persistence semantics must not change.
- Safety net: characterization tests around the use case outputs for representative inputs.
- Seam: introduce a persistence port protocol and an adapter that calls the legacy API.
- Refactor: move persistence calls behind the port; keep behavior the same.
- Next: gradually replace the adapter internals with a new persistence implementation.

### Example: migrating a legacy service to a bounded context

- Boundary: use case I/O and offline behavior must not change.
- Safety net: contract tests for the use case and offline behavior.
- Seam: create a new bounded-context package and move only the use case boundary first.
- Refactor: move domain logic behind the use case boundary in small slices.
- Next: isolate remaining legacy calls behind adapters and strangle them out.
