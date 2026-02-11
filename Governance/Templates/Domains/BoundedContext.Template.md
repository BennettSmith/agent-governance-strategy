# <Context Name> — Bounded Context

## Purpose

<!-- Concise description of why this bounded context exists and the business capability it owns -->

---

## Scope

### In Scope

<!-- Responsibilities and concepts this bounded context owns -->

### Out of Scope

<!-- Explicitly list responsibilities and concepts this bounded context does NOT own -->

---

## Ubiquitous Language

<!-- Shared terms with precise meanings used consistently within this context -->

- **Term** — Definition
- **Term** — Definition

---

## Domain Model (High Level)

### Entities

<!-- Domain objects with identity and lifecycle -->

- **Entity Name**
  - Identity: <!-- What uniquely identifies this entity -->
  - Responsibility: <!-- What this entity is responsible for -->
  - Lifecycle Notes: <!-- Optional -->

### Value Types

<!-- Immutable domain concepts defined by value -->

- **Value Type Name**
  - Meaning: <!-- What this value represents -->
  - Constraints: <!-- Invariants enforced by this value -->

### Aggregates

<!-- Consistency boundaries -->

- **Aggregate Root**
  - Members: <!-- Entities/value types within the aggregate -->
  - Invariants Enforced:
    - <!-- Declarative invariant -->
    - <!-- Declarative invariant -->

### Domain Events (Optional)

<!-- Only include if meaningful reactions or auditing are required -->

- **Event Name**
  - Occurs When: <!-- Past-tense fact -->
  - Key Data: <!-- IDs / minimal values -->

---

## Invariants

<!-- Rules that must always hold true within this bounded context -->

- <!-- Invariant statement -->
- <!-- Invariant statement -->

---

## Boundary Rules

<!-- How this bounded context interacts with the rest of the system -->

- Domain models are internal to this bounded context.
- Domain entities and value types must not cross use-case boundaries.
- Interaction occurs only via application-layer use cases.
- Cross-context interaction uses identifiers, boundary types, or events—not shared models.

---

## Mapping Notes

<!-- High-level guidance only; no implementation details -->

- Use case boundary types are mapped to/from domain models in the application layer.
- Persistence and external API representations are mapped to/from domain models in the infrastructure layer.
- Mapping logic must be explicit.

---

## Offline-First Considerations

<!-- Constraints imposed by offline-first behavior -->

- Domain state must be representable and enforceable locally.
- Invariants must not require server round-trips to validate.
- Eventual consistency assumptions (if any).

---

## Open Questions

<!-- Ambiguities or decisions that must be resolved -->

- <!-- Question -->
- <!-- Question -->

---

## Related Use Cases

<!-- Use cases that operate within this bounded context -->

- <!-- Use Case Name -->
- <!-- Use Case Name -->

---

## Change Log

<!-- Track intentional evolution of the bounded context -->

- YYYY-MM-DD — Initial version
- YYYY-MM-DD — <!-- Description of change -->
