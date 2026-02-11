# Quality gates: Kotlin/JVM architecture enforcement examples

This playbook provides **Kotlin/JVM-specific** examples for enforcing architecture boundaries as part of repo-owned quality gates.

Read first:

- Core principles: `Docs/Playbooks/Target-Repo-Quality-Gates.md`

## Kotlin/JVM tooling patterns

Kotlin/JVM projects often combine:

- **Gradle module boundaries** (strong isolation)
- **static analysis** (Detekt / Ktlint)
- **architecture tests** (e.g., ArchUnit) for rules that are hard to express as simple lint

### Example 1: forbid Android/UI imports in a domain module (Detekt forbidden imports)

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

### Example 2: enforce “no cross-feature imports” with modules

Goal: ensure features don’t depend on each other directly.

Conceptual rule:

- `:featureA` must not depend on `:featureB`; both may depend on `:core` or `:shared`.

The most reliable enforcement is at the build system level:

- declare dependencies explicitly in Gradle
- avoid “catch-all” modules that every feature depends on
- keep `api` usage narrow; prefer `implementation`

You can add a lightweight check in CI that fails if forbidden project dependencies are present (some teams implement this as a small Gradle task).

### Example 3: architecture tests for package/layer rules (ArchUnit)

For rules like “UI may depend on domain, but domain may not depend on UI” when packages are large and imports are indirect, architecture tests are often clearer and less brittle than regex rules.

Architecture tests can:

- run as part of `test` (and therefore `ci`)
- produce clear failure messages
- scale as the codebase grows
