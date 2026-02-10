---
branch: "docs/app-readme-glab-bootstrap"
status: completed
---

## Summary

Update the `README.md` “Copy/paste bootstrap (GitLab Releases)” snippet to use `glab release download` so it works in the Quanata GitLab environment (including private repos) without manual URL/token handling.

## Constraints

- Follow repository governance and keep changes small and reviewable.

## Scope

### In scope

- Update the README bootstrap one-liner to use `glab`.
- Pin the example tag to an existing `agent-gov/v*` tag.

### Out of scope

- Changing the `agent-gov` release / package publishing process or adding new tags.
- Refactoring the Makefile include implementation.

## Approach

- Replace the existing bootstrap one-liner with a `glab release download`-based one-liner.
- Update nearby README references to the pinned tag value used in examples.

## Checkpoints

- [x] Checkpoint 1 — Update README bootstrap one-liner to use `glab`
- [x] Final checkpoint — PR wrap-up (final approval gate)

## Completion checklist

- [x] Set frontmatter `status: completed`
- [x] Check off completed checkpoint(s) above and add PR/commit references as needed
- [x] Check off all quality gates below (or document any exceptions)

## Quality gates / test plan

- [x] `make ci` (passed)
- [ ] Manual: copy/paste the README one-liner and confirm it downloads and runs `agent-gov --version` when authenticated via `glab`

## Notes / open questions

- Merge request: `https://gitlab.com/bsmith.quanata/agent-governance-strategy/-/merge_requests/4`
- Checkpoint commit: `f9575c2`
- Merge conflict resolution: merged `main` and resolved `README.md` conflict (commit `31b4fa5`); re-ran `make ci` (passed)
