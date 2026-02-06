# Plan: Governance via profiles builder (`agent-gov`)

## Summary

Refactor this repo into a **governance builder** that generates profile-specific governance bundles (v1: `mobile-clean-ios`, `backend-go-hex`) and supports deterministic re-sync of **managed doc blocks** in target repos without clobbering local addenda.

## Constraints

- Follow precedence: `Non-Negotiables.md` > `Architecture.md` > `Constitution.md`.
- Work in small, reviewable steps; no large refactors.
- Provide a root `make ci` target intended to be the canonical CI status check.
- `agent-gov` tooling must have **>= 85% test coverage** for `tools/gov/...` and explicit unit + integration tests.
- Generated docs must preserve **Local Addenda (project-owned)** content.

## Scope

### In scope

- Introduce **governance profiles** and fragment-based document assembly.
- Build `agent-gov` (Go CLI) to `build/init/sync/verify` using remote governance content pinned by tag/release, with local caching.
- Add bootstrap script to vendor CLI into a target repo (one-time), after which all operations run from the vendored CLI.
- Add template selection so profiles can emit templates (use cases, bounded contexts, ADR/MADR, etc.) as desired.

### Out of scope

- Android profile (future).
- A production-ready release/distribution strategy (prebuilt binaries, etc.) beyond the vendored CLI approach.

## Approach

- Split governance into **Core** + **Profile** + **Platform overlay/playbooks**.
- Keep the “key three” generated as single files at target repo root with **managed blocks** and a local addenda section.
- Use `.governance/config.yaml` in the target repo:
  - `source.repo`, `source.ref` (tag/release), `source.profile`
  - `paths.docsRoot: "."`
  - cache default: `os.UserCacheDir()` + `/govbuilder` (overrideable)

## Checkpoints

- [ ] Checkpoint 1 — Add decision record (MADR) for governance-by-profiles + managed block sync
- [ ] Checkpoint 2 — Create `Governance/` fragment layout + migrate existing docs into fragments
- [ ] Checkpoint 3 — Implement `agent-gov` CLI skeleton + managed-block sync logic + tests
- [ ] Checkpoint 4 — Add remote fetch + cache + config parsing + integration tests
- [ ] Checkpoint 5 — Add v1 profiles (`mobile-clean-ios`, `backend-go-hex`) outputs + templates selection
- [ ] Checkpoint 6 — Add root `Makefile` with `make ci` and coverage enforcement
- [ ] Checkpoint 7 — Update `AGENTS.md` and template docs to align with profiles strategy

## Quality gates / test plan

- [ ] `make ci` passes (format, tests, >=85% coverage for `tools/gov/...`, bundle generation/verification)

## Notes / open questions

- None (key decisions already captured for v1).

