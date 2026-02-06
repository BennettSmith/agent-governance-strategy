---
branch: feat/app-docs-only-embedded-governance
status: active
---

# Summary

This change introduces a `docs-only` governance profile so this repository’s root key-three docs are generated from a neutral, builder-focused profile. It also treats `tools/gov` as an embedded target repo governed by a Go profile, and adds guardrails (`preflight`) to prevent starting new work on the wrong branch/baseline.

## Constraints

- `Non-Negotiables.md` overrides all other governing documents.
- `Architecture.md` overrides `Constitution.md` on matters of system shape.
- Generated docs must preserve local addenda and only update managed blocks.
- Changes should be small and reviewable; `make ci` is the quality gate.

## Scope

### In scope

- Add `Governance/Profiles/docs-only/` and use it to generate the builder repo root key-three.
- Treat `tools/gov` as an embedded target repo with its own `.governance/config.yaml` and generated key-three.
- Add agent scoping guidance so agents follow the nearest-scope governing docs when working in subtrees.
- Improve `agent-gov` config selection so it can be run “from wherever you are” (upward config auto-discovery).
- Add `agent-gov preflight` + `make preflight` to reduce plan/branch mistakes, using `Docs/Plans/**` as a branch registry (frontmatter preferred; path-derived fallback).

### Out of scope

- Reworking existing profile content beyond what is needed for docs-only and embedded CLI governance.
- Strong “main baseline” checks using commit SHAs/tags (we’ll keep it rough via expected-file presence for now).

## Approach

- **docs-only profile**:
  - Generate root `Non-Negotiables.md` from core.
  - Generate root `Constitution.md` from core + a small docs-only identity overlay.
  - Generate root `Architecture.md` as a builder-focused document (profiles/fragments/templates, managed blocks + addenda, `agent-gov` workflow).
- **Embedded `tools/gov` governance**:
  - Add `tools/gov/.governance/config.yaml` using profile `backend-go-hex` (docsRoot `.`) to generate key-three under `tools/gov/`.
  - Add `tools/gov/AGENTS.md` clarifying that when working under `tools/gov/**`, the embedded key-three governs.
- **Agent scoping rules**:
  - Update root `AGENTS.md` and `.cursor/rules/always-read-governing-docs.mdc` to prefer nearest-scope `AGENTS.md` + key-three when present.
- **Config auto-discovery**:
  - When `--config` is omitted, search upward for `.governance/config.yaml` and use the nearest one; print the selected config path for auditability.
- **Preflight**:
  - Implement `agent-gov preflight` that checks:
    - Not on `main`
    - Not on a different known plan branch from `Docs/Plans/**` (excluding the active plan)
    - Rough baseline artifacts exist
  - Add `make preflight` wrapper.
  - Update `Docs/Plans/Plan.Template.md` to include frontmatter (`branch:`) so future plans participate.

## Checkpoints

- [ ] Checkpoint 1 — Add in-repo plan; scaffold docs-only + embedded-governance configs and docs layout
- [ ] Checkpoint 2 — Implement config auto-discovery + `preflight` command with tests; add `make preflight` and CI wiring
- [ ] Checkpoint 3 — Generate root and embedded key-three; add agent scoping docs/rules; add Go CLI playbook and update governance README

## Quality gates / test plan

- [ ] `make ci` passes
- [ ] `agent-gov verify` passes for root (`docs-only`) and embedded `tools/gov` configs
- [ ] New CLI behavior is covered by unit tests (including error paths)

## Notes / open questions

- We merged `feat/app-governance-profiles-builder` into `main` before starting this work, so we are not using preflight as a gate for *starting* this plan; we are implementing preflight to prevent future mistakes.

