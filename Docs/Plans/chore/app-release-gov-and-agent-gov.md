---
branch: "chore/app-release-gov-and-agent-gov"
status: completed
---

## Summary

Cut a new release of both (1) the governance content bundle and (2) the `agent-gov` CLI tool. Update the root `README.md` to include a clear “latest version” section near the top so adopters can quickly find the current pins.

## Constraints

- Follow repository governing docs: `Non-Negotiables.md`, `Architecture.md`, `Constitution.md`.
- All changes must be made in a feature branch.
- Keep changes small and reviewable; avoid mixing refactors with behavior changes.
- `make ci` must pass.
- When changing files under `tools/gov/**`, follow embedded governance docs in `tools/gov/{Non-Negotiables,Architecture,Constitution}.md`.

## Scope

### In scope

- Add a “Latest versions” section near the start of `README.md`, showing:
  - the latest governance content tag (`gov/v…`)
  - the latest tool tag (`agent-gov/v…`)
- Ensure any pinned/default references in distributed docs/templates are consistent with the new tool tag:
  - `Governance/Templates/Make/agent-gov.mk`
  - `tools/make/agent-gov.mk`
  - `tools/gov/tools/make/agent-gov.mk`
- Ensure release artifacts embed correct version metadata for `agent-gov` on both GitHub Releases and GitLab (Generic Packages).
- Prepare the release/tagging steps for maintainers.

### Out of scope

- Changing governance semantics, profiles, managed-block format, or CLI subcommand behavior (beyond version metadata injection in release pipelines).
- Adding new release pipelines for governance content (tagging is sufficient for `gov/v*`).

## Approach

- Choose next versions:
  - `agent-gov/v1.2.0` (minor feature since `v1.1.0`: version command + build metadata)
  - `gov/v1.0.0` (first SemVer governance content tag per ADR-0004; keep existing date tags intact)
- Update README near the top with a small “Latest versions” callout.
- Update make include defaults to point at the new `agent-gov` tag.
- Update GitLab release build job to pass `-ldflags -X ...` for `Version`, `Commit`, and `Date` (to match GitHub release behavior).
- Run `make ci`.

## Checkpoints

- [x] Checkpoint 1 — Add branch plan (no behavior changes).
- [x] Checkpoint 2 — Update `README.md` with “Latest versions” section and align tag references.
- [x] Checkpoint 3 — Update `agent-gov.mk` defaults to the new tool tag (template + emitted copies).
- [x] Checkpoint 4 — Fix GitLab release build to embed version/commit/date metadata via ldflags; validate locally.
- [x] Checkpoint 5 — Run `make ci` and verify repo is release-ready.
- [x] Final checkpoint — PR wrap-up (final approval gate)

## Completion checklist

- [x] Set frontmatter `status: completed`
- [x] Check off completed checkpoint(s) above and add PR/commit references as needed
- [x] Check off all quality gates below (or document any exceptions)

## Quality gates / test plan

- [x] `make ci`

## Notes / open questions

- Confirm whether the governance content tag should start at `gov/v1.0.0` or use a different initial SemVer value; adjust README pins accordingly before tagging.
- Release steps (after merge to `main`):
  - `git tag gov/v1.0.0 && git push origin gov/v1.0.0`
  - `git tag agent-gov/v1.2.0 && git push origin agent-gov/v1.2.0`

## References

- Merge request: `!12`
- Checkpoint commits:
  - Checkpoint 1 (plan): `6cd4e45`
  - Checkpoint 2 (README + pins + GitLab release metadata): `4b3a7de`
