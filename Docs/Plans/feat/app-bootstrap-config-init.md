---
branch: "feat/app-bootstrap-config-init"
status: active
---

## Summary

Add a new `agent-gov bootstrap` command that helps users generate an initial `.governance/config.yaml` in a target repo during bootstrapping. The command should support a non-interactive mode for CI/scripts, and an interactive mode that can list available profiles (and optionally tags/refs) to guide selection.

## Constraints

- Follow governing docs for this repo and embedded `tools/gov` scope.
- All changes must be small and reviewable.
- Treat CLI surface, config schema/interpretation, and output shape as externally observable behavior.
- Tests must be written before production code; `make ci` must pass.
- Avoid changing existing `init/sync/verify/build` behavior except to add the new command.

## Scope

### In scope

- Add `agent-gov bootstrap` CLI command.
- Generate `.governance/config.yaml` (schemaVersion 1) at the correct location (default: `.governance/config.yaml` under repo root).
- Non-interactive mode with required flags (`--source-repo`, `--source-ref`, `--profile`) and safe defaults for other fields.
- Profile discovery by fetching the governance source and enumerating `Governance/Profiles/*/profile.yaml` (show `id` + `description`).
- Optional: best-effort tag listing (“governance versions”) from the source repo to help choose `source.ref`.
- Optional: `--run-init` to write config and immediately run `init`.
- Documentation update in `README.md` to recommend bootstrap flow.

### Out of scope

- Any config schema version changes (schemaVersion remains `1`).
- New governance profiles or changes to existing profile behavior/content.
- Release/publishing pipeline changes.
- Breaking changes to existing commands/flags or config auto-discovery behavior.

## Approach

- Introduce a small internal package (e.g. `tools/gov/internal/bootstrap`) that:
  - fetches a source repo/ref using existing `internal/source.Fetch`
  - lists available profiles by reading manifests under `Governance/Profiles/`
  - generates a `config.Config`-compatible structure and writes YAML deterministically
- Wire a new `bootstrap` command in `tools/gov/internal/cli/run.go` using `flag.FlagSet`, mirroring existing command patterns and exit-code semantics.
- Keep interactive prompting optional and minimal:
  - enable only when stdin is a TTY and `--non-interactive` is not set
  - otherwise require explicit flags (script-safe)
- Add end-to-end CLI tests in `tools/gov/internal/cli` similar to existing `init/sync/verify/build` tests (using a temporary local git repo as the governance source).

## Checkpoints

- [x] Checkpoint 1 — Define bootstrap behavior + tests (non-interactive config generation/write + overwrite rules) (commit `a68baab`)
- [x] Checkpoint 2 — Implement `bootstrap` command (non-interactive) and wire into CLI usage text (commit `a68baab`)
- [x] Checkpoint 3 — Add profile discovery/listing (fetch + enumerate manifests) and validate selected profile exists (commit `a68baab`)
- [x] Checkpoint 4 — Add optional `--run-init` path and end-to-end test covering config write + init output (commit `a68baab`)
- [x] Checkpoint 5 — (Optional) Add interactive selection + (best-effort) tag listing for “governance versions” (commit `a68baab`)
- [x] Checkpoint 6 — Update `README.md` bootstrap docs and add examples (`bootstrap` and `bootstrap --run-init`) (commit `a68baab`)
- [ ] Final checkpoint — PR wrap-up (final approval gate)

## Completion checklist

- [ ] Set frontmatter `status: completed`
- [ ] Check off completed checkpoint(s) above and add PR/commit references as needed
- [ ] Check off all quality gates below (or document any exceptions)

## Quality gates / test plan

- [x] `make ci` (passed, 2026-02-10)
- [x] Unit/CLI tests for `bootstrap` (non-interactive) passing
- [ ] Manual: run `agent-gov bootstrap --run-init` against a sample target repo and confirm:
  - `.governance/config.yaml` is created in the expected location
  - `init` emits docs/templates for the chosen profile
  - `verify` succeeds using the generated config

## Notes / open questions

- Should `bootstrap` default to writing `.governance/config.yaml` relative to current working directory, or should it attempt to detect git repo root? (Plan assumes “repo root” is current working directory unless `--repo-root` is provided.)
- Tag listing: do we filter by a prefix convention (e.g. `gov/`), or show all tags and let the user choose?
- Interactive mode: if we add it, prefer minimal prompts and keep a strict `--non-interactive` mode for CI.
