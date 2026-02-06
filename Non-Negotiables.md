# Non‑Negotiables

## Architectural

- Use cases define all application behavior
- Bounded contexts are the primary unit of domain ownership
- Domain models (entities, value types, aggregates, domain events) must belong to exactly one bounded context.
- Domain models must not be shared across bounded contexts.
- Interaction with a bounded context must occur only through application-layer use cases.
- Use cases must not expose bounded-context domain models outside their boundary.
- Mapping between bounded-context domain models and external representations must occur outside the domain layer.
- Cross-context communication must occur via use cases and identifiers (IDs), never by sharing domain entities
- Feature modules may not own domain logic
- Domain models must not depend on SwiftUI, networking, persistence, or DI frameworks
- Aggregates must enforce invariants; invariants may not be enforced in UI, presenters, controllers, or infrastructure
- Business logic must live in domain models or application use cases
- Use cases must not accept or return domain entities or value objects
- Use cases must define explicit input and output boundary types
- Controllers must invoke use cases exclusively via the `execute` method
- Use cases must be asynchronous
- Presenters must handle formatting and display concerns only
- UI code must not contain business rules
- Dependencies must be inverted across all architectural boundaries
- Object creation and dependency wiring must occur only in the composition root
- Offline behavior must not be bypassed
- New features must follow the target architecture
- Cross-layer dependencies that violate Clean Architecture are forbidden

## Packaging

- Each bounded context is a Swift Package
- Feature modules may be Swift Packages for delivery units
- Infrastructure adapters are isolated
- All backend communication must use the generated OpenAPI client
- No direct networking calls outside the infrastructure layer

## Documentation

- Every use case has a markdown spec
- A single authoritative use‑case catalog must exist

## Quality Gates

- Tests must be written before production code
- Code coverage must remain above 85%
- `make ci` must pass with no failures
- Quality gates must pass for a task to be considered complete

## Legacy & refactoring

- A change is a **refactor** only if it makes **no externally observable behavior change** at defined boundaries (use case input/output boundaries, persistence contracts, API interactions, and user-visible UI behavior).
- Refactors must **not** be mixed with behavior changes (features/fixes) in the same branch or pull request. If both are needed, do the refactor first, then do the behavior change as a separate step.
- Before changing legacy internals, preserve current behavior with a **safety net** (characterization tests, contract tests, or golden-master tests) at the relevant boundary.
- Refactoring steps must be **small and reversible**. Each step must:
  - keep the prior behavior intact
  - pass all quality gates
  - be trivially revertible (small diff, or guarded by a flag, or a parallel-change/expand-contract approach where the old path remains available until cutover)
- New code introduced during refactoring must follow the target architecture; legacy code may remain only when **isolated behind explicit boundaries** (seams/adapters/facades/use-case boundaries) so violations do not spread.
- Follow `Docs/Refactoring/Legacy-Refactoring-Playbook.md` for approved legacy refactoring patterns and checklists.

## Agents

- Agents must follow the architecture
- Agents must not bypass use cases
- All changes must be made in feature branches
- Agents must request human acceptance before checkpoint commits
- Agents must not mark plan tasks complete without approval
- Agents must not merge branches under any circumstances
