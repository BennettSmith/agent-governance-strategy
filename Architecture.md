# Architecture

## Overview

This system follows Clean Architecture with Domain‑Driven Design influences. Application behavior is expressed exclusively through **use cases**, which form the stable contract between domain logic and delivery mechanisms (UI, background tasks, integrations).

Two organizing principles are used:

- **Bounded Contexts** — define coherent domain language, rules, and invariants
- **Feature Modules (Vertical Slices)** — deliver user‑visible outcomes by composing use cases

## System Responsibility

- Provide a native iOS application that functions fully without network connectivity
- Persist, query, and mutate domain data locally
- Coordinate domain behavior through asynchronous use cases
- Synchronize local state with a backend REST API when connectivity is available
- Present formatted results through a SwiftUI-based user interface

## Explicit Non-Responsibilities

- Backend business rules or policy enforcement
- Authentication or identity provider implementation
- Cross-platform UI concerns
- Backend data modeling beyond API contracts
- UI formatting or presentation logic inside use cases
- Dependency wiring outside the composition root

## Major Layers

- **Domain Layer**
  - Domain entities, value objects, and aggregates
  - Domain services expressing pure business logic
  - No dependencies on application, UI, or infrastructure layers
- **Application Layer**
  - Use cases expressing application behavior
  - Explicit input and output boundary types owned by each use case
  - Asynchronous execution model
  - Depends on abstractions, not concrete implementations
  - No exposure of domain entities outside the use case boundary
- **Interface Layer**
  - Controllers that invoke use cases via `execute`
  - Presenters that transform use case output into UI-ready formats
  - SwiftUI views that render presenter output
  - Navigation/routing components (coordinators/routers) that own user flows
  - No business logic or domain mutation
- **Infrastructure Layer**
  - Persistence implementations (local storage)
  - Network implementations using generated API clients
  - System services (clock, connectivity, identifiers)
  - Implements abstractions defined by the application or domain layers
- **Composition Root**
  - Single place where concrete implementations are assembled
  - Wires dependencies across layers using dependency injection
  - Owns object graph creation and lifecycle decisions
  - Depends on all layers but is depended on by none

## Navigation & Routing

Routing is a UI concern. It must not leak into use cases or domain models.

Ownership:

- **Feature modules / Interface layer** own navigation state and transitions (e.g., SwiftUI
  `NavigationStack`, coordinators/routers, flow state machines).
- **Composition root** assembles the root navigation graph and deep-link entry points, and
  wires feature flows together.

Rules:

- **No navigation in bounded contexts or use cases.** Use cases must not reference routes, screens, view types, SwiftUI navigation APIs, or coordinator/router abstractions.
- **No domain types in routing state.** Routes must not carry domain entities/value objects; they may carry stable identifiers and UI boundary/presentation state as needed.
- **Navigation is driven by UI state, not by side effects.** A use case returns an output boundary describing what happened; the interface layer decides whether that implies a navigation transition.

Common flow patterns:

- **User-driven navigation**: user action → view/controller updates navigation state → destination UI renders → controller invokes use case as needed.
- **Outcome-driven navigation**: user action → controller executes use case → presenter updates view-ready state → interface layer updates navigation state (e.g., show detail on successful creation).
- **Deep links**: app receives URL/notification → composition root maps it to a feature flow entry point → interface layer drives navigation to the appropriate screen and triggers required use cases.

## Bounded Contexts

A bounded context represents a cohesive domain model with a shared ubiquitous language. In this system, **each bounded context is implemented as its own Swift Package**.

- The system is decomposed into **bounded contexts**, each representing a distinct area of the domain.
- Each bounded context resides within the **Domain layer** and owns its:
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
- Each bounded context is documented in a dedicated domain document (e.g. `docs/domain/<context>.md`) and referenced by use cases that operate within it.

A bounded‑context package:

- Owns domain entities, value types, and aggregates
- May define domain events (optional), scoped to the bounded context
- Defines all use cases for that context
- Defines explicit input/output boundary types (DTOs)
- Exposes ports (protocols) for persistence and external services

Rules:

- No UI or navigation
- No platform or framework dependencies
- No direct dependencies on other bounded contexts
- Domain events (if used) are internal to the bounded context and intended to be raised/handled within the context (typically within an aggregate boundary)
- Domain events must not be used as cross-bounded-context integration contracts; cross-context communication happens via use cases and explicit boundary/application-level messages

## Use Cases

Use cases are the **only place where application behavior is defined**.

Characteristics:

- Executed by controllers (UI, background tasks, system events)
- Accept explicit input types (no domain entities)
- Return explicit output/result types
- Perform explicit mapping between boundary types (DTOs) and domain models at the use case boundary (application layer)
- Enforce domain invariants internally
- Are asynchronous and side‑effect aware

## Feature Modules (Vertical Slices)

Feature modules represent user journeys or deliverable functionality.

Rules:

- No domain logic or invariants
- No persistence decisions
- Must not introduce new business rules

## Infrastructure & Adapters

Infrastructure packages implement bounded‑context ports using platform or external systems.

## Composition Root

The application target wires feature modules, bounded contexts, and adapters.

## Dependency Direction

UI / Features → Use Cases → Domain
Adapters → Ports → Use Cases

## Data and Control Flow

- UI events are handled by controllers
- Controllers invoke asynchronous use cases via `execute`
- Use cases accept input boundary objects and return output boundary objects
- Domain entities remain internal to the use case boundary
- Results are passed to presenters for formatting
- Presenters expose view-ready state to SwiftUI
- Local persistence occurs before or independently of network synchronization
- All concrete dependencies are provided by the composition root

## External Boundaries

- Backend interaction occurs exclusively through generated SDK clients
- API contracts are defined solely by OpenAPI specifications
- No handwritten networking code bypassing the generated client

## Target Architecture

- Fully offline-capable core with eventual consistency
- Strict separation between domain, application, interface, infrastructure, and composition
- Use cases form the stable boundary of the system
- Dependency direction always points inward
- Legacy or shortcut implementations are temporary and isolated

## Legacy modernization

When working in legacy areas that do not yet follow the target architecture:

- Keep new code aligned to the target architecture; do not spread boundary violations.
- Prefer seams and incremental migration patterns (e.g., strangler/branch-by-abstraction/expand–contract) over rewrites.
- Follow `Docs/Refactoring/Legacy-Refactoring-Playbook.md` and the refactoring protocol in `Constitution.md`.

