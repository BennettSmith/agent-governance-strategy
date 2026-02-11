## Overview (mobile-clean)

This profile follows Clean Architecture with Domain‑Driven Design influences. Application behavior is expressed exclusively through **use cases**, which form a stable contract between domain logic and delivery mechanisms (UI, background tasks, integrations).

Two organizing principles are used:

- **Bounded Contexts** — define coherent domain language, rules, and invariants
- **Feature Modules (Vertical Slices)** — deliver user‑visible outcomes by composing use cases

## System responsibility (mobile-clean)

- Provide a native mobile application that functions fully without network connectivity.
- Persist, query, and mutate domain data locally.
- Coordinate domain behavior through asynchronous use cases.
- Synchronize local state with a backend REST API when connectivity is available.
- Present formatted results through a platform UI layer in accordance with profile playbooks.

## Explicit non-responsibilities (mobile-clean)

- Backend business rules or policy enforcement
- Authentication or identity provider implementation
- Cross-platform UI concerns (the profile governs an app, but UI toolkit mechanics are playbook concerns)
- Backend data modeling beyond API contracts
- UI formatting or presentation logic inside use cases
- Dependency wiring outside the composition root

## Major layers (mobile-clean)

- **Domain layer**
  - Domain entities, value objects, and aggregates
  - Domain services expressing pure business logic
  - No dependencies on application, UI, or infrastructure layers
- **Application layer**
  - Use cases expressing application behavior
  - Explicit input and output boundary types owned by each use case
  - Asynchronous execution model
  - Depends on abstractions, not concrete implementations
  - No exposure of domain entities outside the use case boundary
- **Interface layer**
  - Controllers that invoke use cases via `execute`
  - Presenters that transform use case output into UI-ready formats
  - UI views that render presenter output (toolkit-specific)
  - Navigation/routing components that own user flows (toolkit-specific)
  - No business logic or domain mutation
- **Infrastructure layer**
  - Persistence implementations (local storage)
  - Network implementations using generated API clients
  - System services (clock, connectivity, identifiers)
  - Implements abstractions defined by the application or domain layers
- **Composition root**
  - Single place where concrete implementations are assembled
  - Wires dependencies across layers using dependency injection
  - Owns object graph creation and lifecycle decisions
  - Depends on all layers but is depended on by none

## Navigation & routing (mobile-clean)

Routing is a UI concern. It must not leak into use cases or domain models.

Ownership:

- **Feature modules / interface layer** own navigation state and transitions (toolkit-specific).
- **Composition root** assembles the root navigation graph and deep-link entry points, and wires feature flows together.

Rules:

- **No navigation in bounded contexts or use cases.** Use cases must not reference routes, screens, view types, navigation APIs, or coordinator/router abstractions.
- **No domain types in routing state.** Routes must not carry domain entities/value objects; they may carry stable identifiers and UI boundary/presentation state as needed.
- **Navigation is driven by UI state, not by side effects.** A use case returns an output boundary describing what happened; the interface layer decides whether that implies a navigation transition.

Common flow patterns:

- **User-driven navigation**: user action → view/controller updates navigation state → destination UI renders → controller invokes use case as needed.
- **Outcome-driven navigation**: user action → controller executes use case → presenter updates view-ready state → interface layer updates navigation state.
- **Deep links**: app receives URL/notification → composition root maps it to a feature flow entry point → interface layer drives navigation and triggers required use cases.

## Bounded contexts (mobile-clean)

- The system is decomposed into **bounded contexts**, each representing a distinct area of the domain.
- Each bounded context resides within the **domain layer** and owns its:
  - domain entities
  - value types
  - aggregates
  - domain events (if applicable)
- Bounded contexts are **isolated from one another** and do not share domain models.
- Application-layer use cases act as the **only entry point** into a bounded context.
- Mapping between bounded-context domain models and:
  - use case boundary types
  - persistence representations
  - external API representations  
    occurs outside the domain layer.
- The architecture enforces **one-way dependency direction**:
  - domain → nothing
  - application → domain
  - interface/infrastructure → application abstractions

## Use cases (mobile-clean)

Use cases are the **only place where application behavior is defined**.

Characteristics:

- Executed by controllers (UI, background tasks, system events)
- Accept explicit input types (no domain entities)
- Return explicit output/result types
- Perform explicit mapping between boundary types (DTOs) and domain models at the use case boundary (application layer)
- Enforce domain invariants internally
- Are asynchronous and side‑effect aware

## Feature modules (vertical slices) (mobile-clean)

Feature modules represent user journeys or deliverable functionality and **compose** use cases.

Rules:

- No domain logic or invariants
- No persistence decisions
- Must not introduce new business rules

## External boundaries (mobile-clean)

- Backend interaction occurs exclusively through generated SDK clients.
- API contracts are defined solely by OpenAPI specifications.
- No handwritten networking code bypassing the generated client.

## Target architecture (mobile-clean)

- Fully offline-capable core with eventual consistency
- Strict separation between domain, application, interface, infrastructure, and composition
- Use cases form the stable boundary of the system
- Dependency direction always points inward
- Legacy or shortcut implementations are temporary and isolated
