## Architectural (backend-go-hex)

- Service behavior must be defined at explicit boundaries (ports/handlers) and tested at those boundaries.
- Domain and application logic must not depend on transport frameworks (HTTP, gRPC) or persistence implementations.
- Ports (interfaces) are defined inward (domain/application) and implemented outward (adapters).
- Adapters must be isolated; infrastructure/framework code must not leak into domain logic.
- Dependency wiring must occur only in the composition root.

## Documentation (backend-go-hex)

- Architectural decisions affecting system shape or long-term constraints must be captured as MADRs/ADRs.

