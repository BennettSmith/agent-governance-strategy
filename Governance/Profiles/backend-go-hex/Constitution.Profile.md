## Architectural identity (backend-go-hex)

- Backend Go microservice.
- Hexagonal architecture (ports and adapters).
- API-first integration with external contracts (OpenAPI/gRPC/async messages) depending on the service profile.
- Infrastructure and transport are replaceable details, not the core of the system.

## Fundamental rules (backend-go-hex)

- Inbound adapters (HTTP/gRPC/consumers) must depend on application ports, not on domain internals.
- Outbound adapters (DB, queues, external APIs) must be invoked via outbound ports/interfaces.
- Business rules live in domain/application code, not in handlers/controllers/framework glue.

