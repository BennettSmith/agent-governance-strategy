---
branch: "feat/app-bootstrap-defaults"
status: active
---

## Summary

Improve `agent-gov bootstrap` interactive UX by supporting sensible defaults and reducing required typing. Add support for environment-variable defaults (and optional inference from git remotes in interactive mode) while preserving script-safe non-interactive behavior.

## Constraints

- Follow repository governance and embedded `tools/gov` scope key-three.
- Treat CLI surface (commands/flags/stdout/stderr/exit codes) as externally observable behavior.
- Keep changes small and reviewable; add/extend tests at the CLI boundary.
- Tests must be written before production code; `make ci` must pass (coverage >= 85%).
- Preserve existing config schema and interpretation (`.governance/config.yaml`, schemaVersion 1).

## Scope

### In scope

- **Env-var defaults** for bootstrap inputs (used in both interactive and non-interactive modes):
  - `AGENT_GOV_SOURCE_REPO` → default for `--source-repo`
  - `AGENT_GOV_SOURCE_REF` → default for `--source-ref` (if unset, default `HEAD` in interactive)
  - `AGENT_GOV_PROFILE` → default for `--profile`
  - `AGENT_GOV_DOCS_ROOT` → default for `--docs-root` (fallback remains `.`)
  - (Optional) `AGENT_GOV_CACHE_DIR` → default for `--cache-dir` (still avoid writing `paths.cacheDir` unless explicitly requested)
- **Interactive-mode defaults**:
  - Show prompts with defaults (e.g. `source ref [HEAD]:`).
  - If `--source-repo` is unset and `AGENT_GOV_SOURCE_REPO` is unset, attempt to infer a default from git remotes (best-effort) and present it as the default choice.
- **Non-interactive behavior**:
  - Continue to support `--non-interactive`.
  - Accept env vars as satisfying required values (i.e. required = flag OR env).
- **Docs**: update `README.md` (and/or a short doc section) describing the new env vars and precedence.

### Out of scope

- Changing `.governance/config.yaml` schema fields or defaults beyond existing behavior.
- Adding network-driven selection UIs (no TUI); keep prompting simple.
- Adding automatic default profile selection without explicit user/env selection.

## Approach

### Precedence rules (single source of truth)

For each bootstrap field, compute an effective value using:

1. **Flag** (explicit CLI arg)
2. **Environment variable** (new defaults)
3. **Interactive inference** (interactive mode only; best-effort)
4. **Interactive prompt default** (e.g. `HEAD` for ref, `.` for docsRoot)
5. **Error** (non-interactive only when still missing)

Keep these rules implemented in one place (helper functions) so tests can pin behavior.

### Git remote inference (interactive only)

- When `sourceRepo` is empty after flags+env:
  - Run `git remote -v` in the detected repo root (best-effort; no hard failure if git missing/not a repo).
  - If any remote URL contains `agent-governance-strategy` (or matches a configurable substring), use the first such URL as a suggested default.
  - Otherwise, do not guess; prompt normally.

### Prompt UX changes

- Print prompts as `label [default]:` and accept empty input to select default.
- For profile selection:
  - If `--profile`/`AGENT_GOV_PROFILE` is set, skip selection and validate exists.
  - Otherwise keep existing “list + choose number” flow.

### Tests (safety net)

Add CLI tests that lock precedence and avoid regressions:

- Env var used in `--non-interactive` mode when flag omitted (repo/ref/profile).
- Flag overrides env var.
- Interactive mode uses `HEAD` default for `source-ref` when empty input.
- (If implemented) git-remote inference is used as the prompt default when env+flag missing.
- Ensure `--cache-dir` behavior remains: used for fetch, but only written to config when explicitly provided (current behavior).

## Checkpoints

- [x] Checkpoint 1 — Plan + tests for precedence rules (flag/env/inferred/prompt) and `--non-interactive` requirements
- [x] Checkpoint 2 — Implement env-var defaults for `bootstrap` (repo/ref/profile/docsRoot/cacheDir)
- [x] Checkpoint 3 — Improve prompts to display defaults and accept empty input (commit `6b06df6`)
- [ ] Checkpoint 4 — (Optional) Add git remote inference for default `source-repo` in interactive mode + tests (skipped; prefer env/flags over heuristics)
- [x] Checkpoint 5 — Update docs (`README.md`) explaining env vars and precedence + examples (commit `e741fb9`)
- [ ] Final checkpoint — PR wrap-up (final approval gate)

## Completion checklist

- [ ] Set frontmatter `status: completed`
- [ ] Check off completed checkpoint(s) above and add PR/commit references as needed
- [ ] Check off all quality gates below (or document any exceptions)

## Quality gates / test plan

- [x] `make ci` (passed, 2026-02-10)
- [x] Unit/CLI tests added for env-var defaults and precedence
- [x] Integration: `bootstrap` (env defaults) generates config and `init` + `verify` succeed (covered by CLI test)

## Notes / open questions

- Do we want to standardize the env var prefix as `AGENT_GOV_*` (proposed) or `GOV_*` for consistency with Makefile variables?
- Git remote inference: should we match only `agent-governance-strategy` substring, or support an explicit env var like `AGENT_GOV_SOURCE_REPO_HINT`?
