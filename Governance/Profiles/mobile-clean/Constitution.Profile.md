## Architectural identity (mobile-clean)

- Offline-first, client-centric mobile application.
- Clean Architecture with Domain-Driven Design as the primary modeling approach.
- Use-case–driven application core.
- API-first integration with backend services defined by OpenAPI.
- Dependency-injected system with an explicit composition root.
- UI and infrastructure are replaceable details, not the core of the system.

## Fundamental rules (mobile-clean)

### 1. Use cases are the source of truth

All application behavior **must** be represented by a documented use case within a bounded context.

### 2. Bounded context ownership

Every use case belongs to **exactly one bounded context**.

### 3. Explicit boundaries

Use cases accept explicit input/output boundary types and never expose domain entities.

Domain models are **internal to a bounded context** and must not cross context boundaries.

Mapping between boundary types and domain models is **explicit**, performed at the boundary, and owned by the bounded context.

### 4. Dependency direction

Dependencies must always point inward.

### 5. No domain logic outside use cases

Domain rules and invariants live inside use cases or their owned entities.

### 6. Offline‑first by default

Use cases must function without network connectivity unless explicitly documented.
