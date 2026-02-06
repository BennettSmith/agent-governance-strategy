# iOS packaging playbook (SwiftPM)

This playbook describes how the `mobile-clean` governance concepts map onto **Swift Package Manager** (SwiftPM) on iOS.

## Goals

- Make bounded-context boundaries enforceable at build time.
- Keep dependency direction pointing inward (domain/application are independent of UI/infrastructure).
- Make it easy to test use cases and domain logic in isolation.

## Modules (SwiftPM packages/targets)

### Bounded contexts

**Rule mapping:** “Each bounded context is its own independently buildable module/package.”

Recommended mapping:

- One bounded context = one SwiftPM package (or a package with multiple targets if you need sub-modules).
- The bounded context package exposes only:
  - use case boundary types (commands/results/errors)
  - use case interfaces/`execute` entrypoints (or public use case types)
  - ports/protocols needed by the bounded context
- Keep domain entities/value types **internal** to the package/target unless a type is explicitly a boundary type.

### Feature modules (vertical slices)

Feature modules may be separate SwiftPM packages/targets used to deliver UI flows. They:

- compose use cases
- own routing/navigation and view state
- must not own domain invariants or persistence decisions

### Infrastructure adapters

Infrastructure lives outside bounded contexts and implements ports/protocols defined inward.

Examples:

- persistence adapters (SQLite/CoreData wrappers)
- network adapters (generated OpenAPI client usage)
- system services (clock, connectivity, identifiers)

## Dependency rules (practical)

- Bounded-context packages must not depend on SwiftUI/UIKit, URLSession/network frameworks, persistence frameworks, or DI frameworks.
- Feature/UI packages may depend on UI frameworks and presentation libraries, but must depend on bounded contexts **only through use case boundaries**.
- Infrastructure packages depend on external SDKs/frameworks and implement ports defined by bounded contexts/application.
- Composition root depends on everything and wires the object graph.

## Testing targets

For each bounded-context package:

- Add a `Tests` target for:
  - use-case contract tests (boundary-level)
  - domain invariants/aggregate tests
  - mapping tests (DTO <-> domain mapping at the boundary)

For adapters:

- Add tests that validate mapping and error handling at the adapter boundary (do not re-test domain invariants in adapter tests).

## Notes

- SwiftPM is the default packaging mechanism for iOS in this profile.\n+- If your app uses a different module system, keep the same boundaries and dependency direction; document the mapping in a local addendum.

