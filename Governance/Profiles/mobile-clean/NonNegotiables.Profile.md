## Architectural (mobile-clean)

- Use cases define all application behavior.
- Bounded contexts are the primary unit of domain ownership.
- Domain models (entities, value types, aggregates, domain events) must belong to exactly one bounded context.
- Domain models must not be shared across bounded contexts.
- Interaction with a bounded context must occur only through application-layer use cases.
- Use cases must not expose bounded-context domain models outside their boundary.
- Mapping between bounded-context domain models and external representations must occur outside the domain layer.
- Cross-context communication must occur via use cases and identifiers (IDs), never by sharing domain entities.
- Feature modules may not own domain logic.
- Domain models must not depend on UI frameworks, networking, persistence, or DI frameworks.
- Aggregates must enforce invariants; invariants may not be enforced in UI, presenters, controllers, or infrastructure.
- Business logic must live in domain models or application use cases.
- Use cases must not accept or return domain entities or value objects.
- Use cases must define explicit input and output boundary types.
- Controllers must invoke use cases exclusively via the `execute` method.
- Use cases must be asynchronous.
- Presenters must handle formatting and display concerns only.
- UI code must not contain business rules.
- Dependencies must be inverted across all architectural boundaries.
- Object creation and dependency wiring must occur only in the composition root.
- Offline behavior must not be bypassed.
- New features must follow the target architecture.
- Cross-layer dependencies that violate Clean Architecture are forbidden.

## Packaging (mobile-clean)

- Each bounded context is its own independently buildable module/package with explicit dependencies.
- Feature modules may be separate delivery modules/packages, but may not own domain logic.
- Infrastructure adapters are isolated behind ports/interfaces defined inward.
- All backend communication must use the generated OpenAPI client.
- No direct networking calls outside the infrastructure layer.

## Documentation (mobile-clean)

- Every use case has a markdown spec.
- A single authoritative use‑case catalog must exist.

