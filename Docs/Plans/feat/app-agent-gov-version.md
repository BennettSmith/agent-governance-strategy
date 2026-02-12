---
branch: "feat/app-agent-gov-version"
status: completed
---

## Summary

Add a first-class version display to the `agent-gov` CLI (`version` subcommand and `--version` flag) so target repos and CI can quickly confirm which binary is installed. Update the root `README.md` to include a quick start for adopting `agent-gov` and a complete list of available commands.

## Constraints

- Follow `tools/gov` embedded governance (backend-go-hex) and repo non-negotiables.
- Keep changes small and reviewable; do not mix refactors with behavior changes.
- Tests come before production code for new behavior.
- `make ci` must pass.

## Scope

### In scope

- Add `agent-gov version` command.
- Add `agent-gov --version` (and `-v`) global-style behavior.
- Print build metadata (version, commit, build date) with sane defaults for local dev builds.
- Wire release/build pipeline to inject version metadata at build time (ldflags) so released binaries report the correct version.
- Update root `README.md`:
  - Add a “Quick start” section near the top describing the fastest path to adopt `agent-gov` in a target repo.
  - Ensure all available `agent-gov` commands and primary options are listed.

### Out of scope

- Changing governance content semantics, profiles, or sync behavior.
- Adding new subcommands beyond version-related work.
- Rewriting the CLI framework (it currently uses Go’s `flag`).

## Approach

- Implement version output in the `tools/gov/internal/cli` layer:
  - Recognize `version`, `--version`, and `-v` at the top-level command dispatch.
  - Print a stable, script-friendly value (always includes version; may include commit/date when present).
- Provide build metadata via a small internal package with exported variables set via `-ldflags -X`.
- Add focused unit tests around `cli.Run(...)` for version command/flag behavior and output.
- Update `README.md` to include:
  - Quick start steps (download/install, bootstrap config, init/sync/verify usage).
  - A command reference section that matches the CLI usage (`help`, `preflight`, `bootstrap`, `init`, `sync`, `verify`, `build`, `version`).

## Checkpoints

- [x] Checkpoint 1 — Commit the branch plan (no behavior changes).
- [x] Checkpoint 2 — Add tests specifying `agent-gov version` / `--version` behavior and expected output shape.
- [x] Checkpoint 3 — Implement version handling + build metadata variables; ensure tests pass.
- [x] Checkpoint 4 — Wire build/release to inject version/commit/date; verify locally.
- [x] Checkpoint 5 — Update root `README.md` quick start + command list; ensure docs reflect actual CLI behavior.
- [x] Checkpoint 6 — Vendor Go deps + update markdown tooling ignores (if required) and run `make ci`.
- [x] Checkpoint 7 — Exclude vendored Go files from gofmt checks; re-run `make ci`.
- [x] Final checkpoint — PR wrap-up (final approval gate)

## Completion checklist

- [x] Set frontmatter `status: completed`
- [x] Check off completed checkpoint(s) above and add PR/commit references as needed
- [x] Check off all quality gates below (or document any exceptions)

## Quality gates / test plan

- [x] `go test ./...` (at least under `tools/gov/...`)
- [x] `make ci`

## Notes / open questions

- Decide final output format for `--version` (single line vs including commit/date). Prefer single line with optional suffix to remain human-friendly while staying script-tolerant.

## References

- Merge request: `!11`
- Merge request URL: `https://gitlab.com/bsmith.quanata/agent-governance-strategy/-/merge_requests/11`
- Checkpoint commits:
  - Checkpoint 1 (plan): `e81cb02`
  - Checkpoint 2 (tests): `24e4e54`
  - Checkpoint 3 (implementation): `e3f71d2`
  - Checkpoint 4 (release): `41b879c`
  - Checkpoint 5 (README): `614c4f9`
  - Checkpoint 6 (vendor + ignores): `117104f`
  - Checkpoint 7 (gofmt excludes vendor): `4966200`
