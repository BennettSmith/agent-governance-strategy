# Quality gates: Go architecture enforcement examples

This playbook provides **Go-specific** examples for enforcing architecture boundaries as part of repo-owned quality gates.

Read first:

- Core principles: `Docs/Playbooks/Target-Repo-Quality-Gates.md`

## Go tooling patterns

Go projects often enforce architecture through a mix of:

- **module boundaries** (separate packages/modules for hard isolation)
- **import restrictions** (preventing forbidden dependencies between layers)
- **compile-time/test-time checks** that run in `make lint` / `make ci`

### Example 1: forbid importing an infrastructure adapter from a core package (imports rule)

Goal: keep core logic independent from infrastructure concerns.

Conceptual rule:

- Packages under `internal/core/**` must not import packages under `internal/adapters/**`.

One practical approach is to enforce this with `golangci-lint` using `depguard` in your `.golangci.yml`:

```yaml
linters:
  enable:
    - depguard

linters-settings:
  depguard:
    rules:
      core-no-adapters:
        files:
          - "internal/core/**/*.go"
        deny:
          - pkg: "your.module/internal/adapters"
            desc: "core must not depend on adapters; depend on ports/interfaces instead"
```

Notes:

- Prefer scoping by **file glob** so the rule is precise and fast.
- Use a deny list that matches your module import paths.

### Example 2: forbid `net/http` usage in the domain layer (purity boundary)

Goal: prevent accidental coupling to transport concerns.

Conceptual rule:

- Packages under `internal/domain/**` must not import `net/http`.

Example (also `depguard`):

```yaml
linters-settings:
  depguard:
    rules:
      domain-no-http:
        files:
          - "internal/domain/**/*.go"
        deny:
          - pkg: "net/http"
            desc: "domain must not depend on net/http; keep transport concerns at the edge"
```

### Example 3: enforce “no TODOs” or “no fmt.Printf” in production code (optional)

These rules are not architecture boundaries, but they are common “quality gates” that can be enforced mechanically. If you use them, keep them scoped and actionable.

Examples (via `golangci-lint` linters such as `forbidigo`, `godox`, etc.) are viable, but prefer rules that don’t create noise.
