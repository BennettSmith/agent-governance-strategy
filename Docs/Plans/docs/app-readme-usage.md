---
branch: "docs/app-readme-usage"
status: active
---

# Plan: Update README with local usage instructions

## Summary

Add a root `README.md` that explains how to use this repository as a governance source + toolchain, with team-safe setup guidance (remote `source.repo`, pinned `source.ref`) and a local-dev override workflow.

## Constraints

- Follow repository governance (`Non-Negotiables.md`, `Architecture.md`, `Constitution.md`).
- Keep changes small, reviewable, and docs-focused.
- Avoid team-hostile setup that requires committing machine-local paths.

## Scope

### In scope

- Add a root `README.md` that documents:
  - What this repo is (`Governance/` sources + `tools/gov` CLI).
  - How to apply governance to a target repo (team-safe default).
  - Profiles overview and common commands (`init`, `sync`, `verify`, `build`).
  - Team-safe `.governance/config.yaml` patterns (remote URL, pinned ref).
  - Recommended team workflow for distributing the CLI without copying `tools/gov/` into every repo:
    - pinned binary version (tool tags like `agent-gov/vX.Y.Z`)
    - example Makefile/script snippet for download-on-demand into `tools/bin/agent-gov`
  - Optional local override workflow using `--config` and a gitignored dev config.

### Out of scope

- Changing the CLI behavior, config schema, or profiles.
- Adding new profiles or templates.
- Publishing/releasing binaries.

## Approach

- Use `README.md` as the primary onboarding surface (there is currently no root README).
- Provide copy/paste config examples, with explicit guidance on remote URLs vs local paths.
- Keep instructions aligned to how `agent-gov` actually resolves `source.repo` (git clone into a cache) and how config discovery works.

## Checkpoints

- [x] Checkpoint 1 — Create branch plan (`Docs/Plans/docs/app-readme-usage.md`).
- [x] Checkpoint 2 — Draft `README.md` with team-safe usage instructions and examples.
- [x] Checkpoint 3 — Run `make ci` and polish wording for clarity and correctness.

## Completion checklist

- [ ] Set frontmatter `status: completed`
- [ ] Check off completed checkpoint(s) above and add PR/commit references as needed
- [ ] Check off all quality gates below (or document any exceptions)

## Quality gates / test plan

- [x] `make ci`

## Notes / open questions

- Pinning conventions:
  - CLI tool tag: `agent-gov/vX.Y.Z` (SemVer)
  - Governance content tag: `gov/vYYYY.MM.DD` (date-based)
- Follow-up work: CI/release automation plan in `Docs/Plans/chore/app-agent-gov-release-pipeline.md`
