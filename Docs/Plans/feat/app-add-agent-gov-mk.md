---
branch: feat/app-add-agent-gov-mk
status: completed
---

## Summary

Provide a shared `agent-gov.mk` Makefile include that target repositories can `include` to standardize `agent-gov` download + `gov-*` targets, instead of copy/pasting snippets from the README.

## Constraints

- Keep changes small and reviewable.
- Preserve the existing README guidance while making the shared include the recommended path.
- `make ci` remains the quality gate.

## Scope

### In scope

- Add `Governance/Templates/Make/agent-gov.mk` containing the shared make targets.
- Emit `tools/make/agent-gov.mk` from existing profiles so target repos receive the include on `init/sync`.
- Update `README.md` to recommend `-include tools/make/agent-gov.mk` and show minimal configuration for GitHub/GitLab downloads.

### Out of scope

- Changing release asset naming conventions.
- Adding new profiles solely for Makefile integration.
- Non-Make build tooling (e.g., task runners) beyond brief mention.

## Approach

- Create a single `agent-gov.mk` include that supports both GitHub Releases and GitLab Generic Package Registry via `AGENT_GOV_SOURCE ?= github|gitlab`.
- Keep the include self-contained and POSIX-sh compatible (`/bin/sh`).
- Keep target-facing variables explicit and documented at the top of the file.

## Checkpoints

- [x] Checkpoint 1 — Add `agent-gov.mk` template and emit it from profiles.
- [x] Checkpoint 2 — Update `README.md` to use `agent-gov.mk` via `include` and add copy/paste bootstrap one-liner.
- [x] Checkpoint 3 — Add GitLab CI check: plan status must be `completed` for the branch plan file
- [x] Checkpoint 4 — Fix CI job parser for Alpine awk
- [x] Final checkpoint — PR wrap-up (final approval gate)

## Quality gates / test plan

- [x] `make ci` (re-run for final wrap-up)
- [x] `tools/gov` unit tests cover manifest template emission paths (covered via `make ci`)
- [ ] Manual: confirm example Makefile with `-include tools/make/agent-gov.mk` works for both `AGENT_GOV_SOURCE=github` and `AGENT_GOV_SOURCE=gitlab` (download paths only; no credentials committed)

## Notes / open questions

- None.

## References

- Checkpoint commit (Make include + README): `5d9abee`
- Checkpoint commit (CI plan status gate): `6e611a8`
- Checkpoint commit (Alpine awk fix): `4eceb36`
- PR: `https://gitlab.com/bsmith.quanata/agent-governance-strategy/-/merge_requests/3`
