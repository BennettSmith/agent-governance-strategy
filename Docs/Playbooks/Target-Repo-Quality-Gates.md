# Target repo quality gates

This playbook describes how teams should implement **repeatable, tool-enforced quality gates** in target repositories, beyond the governance rules themselves.

Governance provides shared expectations (and verification of managed blocks). Target repos still need **repo-owned** enforcement for product code quality, architecture boundaries, and delivery checks.

## Principles

- **Make it repeatable**: the same command should behave the same on every machine and in CI.
- **Make it fast**: prefer checks that finish quickly and can run on every change.
- **Make it deterministic**: avoid checks that depend on network, time, or nondeterministic ordering.
- **Fail loudly**: errors should be actionable (clear messages, links, and next steps).
- **Prefer static enforcement**: compile-time and static analysis checks are usually more reliable than heuristics.

## Golden commands (repo-owned)

Most teams standardize on:

- `make fmt` — formatting (may write files)
- `make lint` — lint/static analysis (no writes)
- `make test` — tests (no writes)
- `make ci` — the CI gate that runs the full required set (no writes)

Your `ci` target is the contract: if `make ci` is green, the change is acceptable.

## Governance verification fits alongside product checks

Governance verification should be explicit and CI-safe:

- In CI, run **verify only** (no sync / no working tree writes).
- Keep governance checks **additive** to product checks, not a replacement.

Recommended wiring in a target repo:

- `make gov-ci` → runs `agent-gov verify ...` (or `make gov-verify`)
- `make ci` includes `gov-ci` explicitly (e.g., `ci: gov-ci`)

This avoids “attach to ci” magic inside shared includes and keeps product stacks free to define their own `ci` contract.

## What to put in `make ci`

As a starting point, `make ci` often includes:

- `make fmt-check` (or a formatter run in “check” mode)
- `make lint`
- `make test`
- `make gov-ci` (governance verify-only)

Keep `make fmt` available for developers (may write), but keep `make ci` write-free.

## Architecture enforcement (examples)

Architecture enforcement should be:

- **explicit** (document what is forbidden and why)
- **mechanical** (enforced by a tool, not by review memory)
- **scoped** (enforce at boundaries you can describe and test)

### SwiftLint-style architecture enforcement examples (iOS)

SwiftLint can enforce architecture rules using custom rules and file path scoping. The examples below show common patterns; tailor them to your project layout.

#### Example 1: forbid importing a UI framework in a “Domain” layer

Goal: keep domain logic independent from UI frameworks.

Conceptual rule:

- Files under `App/Domain/**` must not `import SwiftUI` or `import UIKit`.

Example configuration sketch:

```yaml
custom_rules:
  no_swiftui_in_domain:
    name: "No SwiftUI in Domain"
    included: "App/Domain"
    regex: "^\\s*import\\s+SwiftUI\\b"
    message: "Domain must not depend on SwiftUI. Move UI-facing code to the UI layer."
    severity: error

  no_uikit_in_domain:
    name: "No UIKit in Domain"
    included: "App/Domain"
    regex: "^\\s*import\\s+UIKit\\b"
    message: "Domain must not depend on UIKit. Move UI-facing code to the UI layer."
    severity: error
```

Notes:

- Keep rules **path-scoped** (`included`) to avoid false positives in other layers.
- Prefer clear messages that tell developers where to move the code.

#### Example 2: enforce “no cross-feature imports” (feature isolation)

Goal: prevent a feature module/layer from importing another feature directly.

Conceptual rule:

- `App/Features/<FeatureA>/**` must not import types from `FeatureB`.

A coarse but effective approach uses import pattern matching:

```yaml
custom_rules:
  no_feature_to_feature_imports:
    name: "No cross-feature imports"
    included: "App/Features"
    regex: "^\\s*import\\s+AppFeature(?!Common)\\w+\\b"
    message: "Features must not import other features directly. Depend on shared interfaces or a common module."
    severity: error
```

Notes:

- Regex-based rules are blunt instruments; keep them **simple** and supported by conventions (naming, folder layout).
- For stronger guarantees, prefer build-system/module boundaries (e.g., SPM packages, Xcode targets) and allow SwiftLint rules to serve as an early, fast signal.

#### Example 3: forbid `Foundation` in a “pure” layer (extreme purity)

Some teams define a layer that should not depend on `Foundation` for testability/portability.

```yaml
custom_rules:
  no_foundation_in_pure_layer:
    name: "No Foundation in Pure layer"
    included: "App/Pure"
    regex: "^\\s*import\\s+Foundation\\b"
    message: "Pure layer must not depend on Foundation."
    severity: warning
```

Use warnings sparingly; if the rule matters, prefer errors and provide a migration plan.

### Go tooling architecture enforcement examples

Go projects often enforce architecture through a mix of:

- **module boundaries** (separate packages/modules for hard isolation),
- **import restrictions** (preventing forbidden dependencies between layers), and
- **compile-time/test-time checks** that run in `make lint` / `make ci`.

#### Example 1: forbid importing an infrastructure adapter from a core package (imports rule)

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

#### Example 2: forbid `net/http` usage in the domain layer (purity boundary)

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

#### Example 3: enforce “no TODOs” or “no fmt.Printf” in production code (optional)

These rules are not architecture boundaries, but they are common “quality gates” that can be enforced mechanically. If you use them, keep them scoped and actionable.

Examples (via `golangci-lint` linters such as `forbidigo`, `godox`, etc.) are viable, but prefer rules that don’t create noise.

### Kotlin tooling architecture enforcement examples

Kotlin/JVM projects often combine:

- **Gradle module boundaries** (strong isolation),
- **static analysis** (Detekt / Ktlint), and
- **architecture tests** (e.g., ArchUnit) for rules that are hard to express as simple lint.

#### Example 1: forbid Android/UI imports in a domain module (Detekt forbidden imports)

Goal: keep domain code free of Android/UI dependencies.

Conceptual rule:

- In `:domain`, forbid `android.*`, `androidx.*`, and UI frameworks.

Detekt’s `ForbiddenImport` can enforce this via your `detekt.yml`:

```yaml
style:
  ForbiddenImport:
    active: true
    imports:
      - "android.*"
      - "androidx.*"
      - "kotlinx.coroutines.flow.FlowPreview" # example of intentionally forbidden API
```

Notes:

- Combine with **module boundaries** so the rule is defense-in-depth rather than your only guardrail.
- Prefer patterns that are stable and low-noise.

#### Example 2: enforce “no cross-feature imports” with modules

Goal: ensure features don’t depend on each other directly.

Conceptual rule:

- `:featureA` must not depend on `:featureB`; both may depend on `:core` or `:shared`.

The most reliable enforcement is at the build system level:

- declare dependencies explicitly in Gradle
- avoid “catch-all” modules that every feature depends on
- keep `api` usage narrow; prefer `implementation`

You can add a lightweight check in CI that fails if forbidden project dependencies are present (some teams implement this as a small Gradle task).

#### Example 3: architecture tests for package/layer rules (ArchUnit)

For rules like “UI may depend on domain, but domain may not depend on UI” when packages are large and imports are indirect, architecture tests are often clearer and less brittle than regex rules.

Architecture tests can:

- run as part of `test` (and therefore `ci`)
- produce clear failure messages
- scale as the codebase grows

### Reliability guidance for architecture rules

- Prefer **module boundaries** (targets/packages) for hard isolation; use lint rules as fast feedback.
- Avoid overly clever regexes; choose conventions that are easy to match.
- Minimize false positives; developers will disable noisy rules.
- Provide a clear escape hatch:
  - allow inline disables only with justification
  - record long-lived exceptions via a MADR (see `Docs/Playbooks/Governance-Exceptions.md`)

## Operational guidance

- Start with a small, enforceable set of gates and tighten over time.
- Keep a single source of truth: `make ci` plus a short doc describing what it runs.
- Treat flaky checks as production incidents: fix or remove quickly.
- When introducing a new required gate, provide:
  - a migration plan
  - a clear failure message
  - an owner responsible for keeping it healthy
