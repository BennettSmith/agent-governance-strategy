## Overview (backend-go-hex)

This profile follows hexagonal architecture (ports and adapters). The application core defines ports (interfaces) and orchestrates domain behavior. Inbound and outbound adapters are replaceable details.

## Layers & boundaries (backend-go-hex)

- **Domain**
  - Business rules and invariants
  - No dependencies on infrastructure, transports, or frameworks
- **Application**
  - Use cases / services orchestrating domain behavior
  - Defines inbound and outbound ports
- **Inbound adapters**
  - HTTP/gRPC handlers, message consumers, schedulers
  - Translate transport concepts into application port calls
- **Outbound adapters**
  - Databases, queues, external APIs, clocks, system services
  - Implement outbound ports defined by the application core
- **Composition root**
  - Wires concrete adapters into ports
  - Owns configuration and lifecycle

## Testing strategy (backend-go-hex)

- Prefer tests at port boundaries (application services / handlers calling ports).
- Keep adapters thin; test the mapping and error handling at adapter boundaries.

