---
branch: "docs/app-playbooks-core-profile-split"
status: active
---

## Summary

Refactor governance playbooks so target repositories receive the right operational guidance in the right place: **core playbooks** for universally applicable procedures, and **profile-specific playbooks** for platform/stack-specific guidance. This reduces ambiguity and keeps “governance as code” distribution consistent.

## Constraints

- Follow repository governing docs: `Non-Negotiables.md`, `Architecture.md`, `Constitution.md`.
- Keep changes **small and reviewable** and proceed checkpoint-by-checkpoint with explicit approval.
- Preserve builder repo architecture: governance source-of-truth in `Governance/`, emitted to target repos via profiles.
- Avoid duplication/drift: each playbook must have **one source of truth**.
- Keep core governance **tool-neutral**; tool-specific operational steps belong in playbooks (per `Docs/Decisions/ADR-0003-Tool-Neutral-Core-And-Tooling-Playbooks.md`).
- `make ci` must pass before completion.

## Scope

### In scope

- Produce an inventory/mapping table of existing playbooks to their destination:
  - core vs profile-specific vs builder-only
  - emitted target path (e.g., `Docs/Playbooks/<name>.md`)
- Define and document criteria for:
  - **Core playbooks** (universal target-repo operations)
  - **Profile playbooks** (platform/stack-specific operations)
- Move selected playbooks from `Docs/Playbooks/` into governance sources:
  - `Governance/Core/Playbooks/` (core)
  - `Governance/Profiles/<id>/Playbooks/` (profile-specific)
- Update profile manifests (`Governance/Profiles/<id>/profile.yaml`) so target repos receive the right playbooks.
- Update cross-links so readers can navigate between README ↔ ADRs ↔ playbooks regardless of where the playbook sources live.

### Out of scope

- Changing `.governance/config.yaml` schema or `agent-gov` CLI behavior.
- Introducing new playbook topics unrelated to the refactor.
- Rewriting playbook content beyond small edits required to:
  - remove duplication
  - clarify scope (core vs profile)
  - fix references/paths after the move

## Approach

- Inventory existing playbooks and classify each as:
  - **Core**: applies across repo types (adoption, upgrades, exceptions discipline, governance verification integration)
  - **Profile-specific**: depends on stack/platform conventions (e.g., iOS packaging, Go CLI structure)
  - **Builder-only**: instructions for maintaining this governance repo/toolchain (keep in `Docs/Playbooks/` only if truly builder-specific)
- Treat emitted playbook paths as a compatibility contract:
  - Prefer keeping stable emitted target paths (avoid moving/renaming files in target repos without a breaking-change justification).
  - If a path must change, document the migration and evaluate whether it is a governance-content breaking change (`gov/vX.Y.Z`).
- Establish a single source-of-truth rule:
  - If target repos should receive it, its source must live under `Governance/**`.
  - If it’s builder-only, keep it under `Docs/Playbooks/` and ensure it is not emitted by profiles.
- Target-repo playbook paths (decision):
  - Profiles will emit governance playbooks to `Docs/Playbooks/`.
  - `Docs/Playbooks/Local/` is reserved for target-repo-owned playbooks and notes.
  - Governance will ensure both folders exist by emitting:
    - `Docs/Playbooks/README.md` (explains ownership; target repos generally should not add new files here)
    - `Docs/Playbooks/Local/README.md` (explains it is project-owned and safe to edit/add)
- Run `make ci` after each checkpoint.

## Checkpoints

- [x] Checkpoint 1 — Inventory + classification
  - List each playbook and decide: core vs profile-specific vs builder-only
  - Identify any playbooks that should be split (core procedure + profile examples)
  - Produce an inventory table mapping current files → destination + emitted target path

### Checkpoint 1 inventory (current playbooks → destination + emitted target path)

Notes:

- “**Destination**” refers to the **governance source-of-truth** location:
  - **Core** → `Governance/Core/Playbooks/`
  - **Profile-specific** → `Governance/Profiles/<id>/Playbooks/`
  - **Builder-only** → keep outside `Governance/**` (not emitted to target repos)
- “**Emitted target path**” is the **compatibility contract** path in a target repo (typically under `Docs/Playbooks/`).
- Some rows below are “emitted copies” already (this builder repo is itself a target for the `docs-only` profile; `tools/gov` is an embedded target for `backend-go-hex`).

| Current playbook (in this repo)                            | Destination                         | Governance source-of-truth (expected)                                                                                                           | Emitted target path (contract)                                                                                                   | Split? |
| ---------------------------------------------------------- | ----------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------- | ------ |
| `Docs/Playbooks/GitHub-PR-Workflow.md`                     | Core                                | `Governance/Core/Playbooks/GitHub-PR-Workflow.md`                                                                                               | `Docs/Playbooks/GitHub-PR-Workflow.md`                                                                                           | No     |
| `Docs/Playbooks/GitLab-MR-Workflow.md`                     | Core                                | `Governance/Core/Playbooks/GitLab-MR-Workflow.md`                                                                                               | `Docs/Playbooks/GitLab-MR-Workflow.md`                                                                                           | No     |
| `Docs/Playbooks/Governance-Upgrades.md`                    | Core                                | `Governance/Core/Playbooks/Governance-Upgrades.md`                                                                                              | `Docs/Playbooks/Governance-Upgrades.md`                                                                                          | No     |
| `Docs/Playbooks/Governance-Exceptions.md`                  | Core                                | `Governance/Core/Playbooks/Governance-Exceptions.md`                                                                                            | `Docs/Playbooks/Governance-Exceptions.md`                                                                                        | No     |
| `Docs/Playbooks/Target-Repo-Quality-Gates.md`              | Core + profile-specific examples    | Core principles: `Governance/Core/Playbooks/Target-Repo-Quality-Gates.md` + profile example playbooks under `Governance/Profiles/**/Playbooks/` | Core principles remain at `Docs/Playbooks/Target-Repo-Quality-Gates.md`; profile example outputs **TBD** (new emitted playbooks) | Yes    |
| `tools/gov/Docs/Playbooks/Go-Packaging.md`                 | Profile-specific (`backend-go-hex`) | `Governance/Profiles/backend-go-hex/Playbooks/Go-Packaging.md`                                                                                  | `Docs/Playbooks/Go-Packaging.md`                                                                                                 | No     |
| `tools/gov/Docs/Playbooks/Go-CLI-Structure.md`             | Profile-specific (`backend-go-hex`) | `Governance/Profiles/backend-go-hex/Playbooks/Go-CLI-Structure.md`                                                                              | `Docs/Playbooks/Go-CLI-Structure.md`                                                                                             | No     |
| `tools/gov/Docs/Playbooks/Hexagonal-Ports-And-Adapters.md` | Profile-specific (`backend-go-hex`) | `Governance/Profiles/backend-go-hex/Playbooks/Hexagonal-Ports-And-Adapters.md`                                                                  | `Docs/Playbooks/Hexagonal-Ports-And-Adapters.md`                                                                                 | No     |

Initial split notes for `Target-Repo-Quality-Gates.md` (to execute in later checkpoints):

- **Core principles to keep in core**: the “Principles”, “Golden commands”, and general guidance about wiring governance verification alongside product checks.
- **Profile-specific examples to move out of core**:
  - Go (`depguard` / `golangci-lint`) → `backend-go-hex`
  - iOS (SwiftLint examples) → `mobile-clean-ios`
  - Kotlin/JVM examples → **requires a JVM/Android-oriented profile** (none exists today); treat as “future profile playbook” and avoid leaving Kotlin-only examples in core.
- [x] Checkpoint 2 — Move core playbooks into `Governance/Core/Playbooks/`
  - Ensure no duplication remains in `Docs/Playbooks/`
  - Update links that pointed at old locations
  - Pending approval notes:
    - Moved:
      - `Docs/Playbooks/Governance-Upgrades.md` → `Governance/Core/Playbooks/Governance-Upgrades.md`
      - `Docs/Playbooks/Governance-Exceptions.md` → `Governance/Core/Playbooks/Governance-Exceptions.md`
      - `Docs/Playbooks/Target-Repo-Quality-Gates.md` → `Governance/Core/Playbooks/Target-Repo-Quality-Gates.md`
    - Updated `README.md` playbook links to point at the new governance source paths (while noting the emitted target paths).
- [ ] Checkpoint 3 — Move profile-specific playbooks into `Governance/Profiles/<id>/Playbooks/`
  - Confirm ownership boundaries and avoid cross-profile drift
- [ ] Checkpoint 4 — Update `profile.yaml` manifests to emit correct playbooks into target repos
  - Verify the emitted target paths and naming are consistent across profiles
  - Ensure `Docs/Playbooks/README.md` and `Docs/Playbooks/Local/README.md` are emitted (so folders are created)
- [ ] Checkpoint 5 — Dogfood both agent-gov usages in this repo (sync + verify)
  - Root scope (docs-only profile):
    - Run `agent-gov sync` and `agent-gov verify` against the root `.governance/config.yaml` (or the repo’s established wrapper/targets).
    - Confirm emitted playbooks/docs update as expected and there is no duplication/drift.
  - Embedded scope (`tools/gov`, backend-go-hex profile):
    - Run `agent-gov sync` and `agent-gov verify` against `tools/gov/.governance/config.yaml`.
    - Confirm embedded generated docs/playbooks update as expected.
- [ ] Checkpoint 6 — Consistency pass (README ↔ ADRs ↔ playbooks)
  - Ensure the “happy path” for a target repo adopter is clear
- [ ] Checkpoint 7 — Run quality gates (`make ci`) and ensure no docs/link regressions
- [ ] Final checkpoint — PR wrap-up (final approval gate)

## Completion checklist

- [ ] Set frontmatter `status: completed`
- [ ] Check off completed checkpoint(s) above and add PR/commit references as needed
- [ ] Check off all quality gates below (or document any exceptions)

## Quality gates / test plan

- [ ] `make ci`
- [ ] Manual spot-check: each playbook has a single source-of-truth and is emitted only where intended
- [ ] Manual spot-check: core governance remains tool-neutral; tool-specific steps live only in playbooks
- [ ] Links check (spot-check moved playbooks + README references)
- [ ] Migration sanity: a target repo running `sync` receives the expected playbook set and `Docs/Playbooks/Local/` remains project-owned

## Notes / open questions

- Decision: for “Quality Gates” guidance, keep **principles** in core and move **language/platform-specific examples** into profile playbooks.
