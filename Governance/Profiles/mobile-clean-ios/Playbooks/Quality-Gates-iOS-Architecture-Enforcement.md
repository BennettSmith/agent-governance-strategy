# Quality gates: iOS architecture enforcement examples (SwiftLint)

This playbook provides **iOS-specific** examples for enforcing architecture boundaries as part of repo-owned quality gates.

Read first:

- Core principles: `Docs/Playbooks/Target-Repo-Quality-Gates.md`

## SwiftLint custom rules (examples)

SwiftLint can enforce architecture rules using custom rules and file path scoping. The examples below show common patterns; tailor them to your project layout.

### Example 1: forbid importing a UI framework in a “Domain” layer

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

### Example 2: enforce “no cross-feature imports” (feature isolation)

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

### Example 3: forbid `Foundation` in a “pure” layer (extreme purity)

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
