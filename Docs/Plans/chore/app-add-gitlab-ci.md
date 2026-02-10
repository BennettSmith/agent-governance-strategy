---
branch: "chore/app-add-gitlab-ci"
status: completed
---

## Summary

Add GitLab CI support that mirrors the existing GitHub Actions behavior: run `make ci` on merge requests, and on `agent-gov/*` tags build and publish the same cross-compiled `agent-gov` binaries with persistent download links.

## Constraints

- `make ci` must pass with no failures (coverage >= 85%).
- GitLab project is **private** (internal), so release binaries must be downloadable with authentication (team-safe tokens).
- Prefer small, reviewable changes and checkpoint commits with explicit approval.

## Scope

### In scope

- Add `.gitlab-ci.yml` with MR pipeline parity for `make ci`.
- Add tag pipelines for `agent-gov/v*` and `agent-gov/test/v*` that:
  - run `make ci`
  - build `agent-gov_{linux,darwin}_{amd64,arm64}` (same naming as GitHub)
  - upload binaries to GitLab Generic Package Registry (persistent)
  - create a GitLab Release with asset links to the package files
- Update `README.md` with a GitLab-compatible “pinned binary” Makefile snippet for private projects (token-based download).

### Out of scope

- Running CI on every branch push without an MR.
- Adding Windows builds (not present in existing GitHub release matrix).

## Approach

- Implement in checkpoints:
  - Plan file committed first.
  - MR pipeline job (make ci).
  - Tag build/upload jobs.
  - GitLab Release creation job.
  - README doc update for private GitLab consumption.
- Keep each checkpoint small and ask for explicit approval before committing.

## Checkpoints

- [x] Checkpoint 0 — Add branch plan file (`c27c430`)
- [x] Checkpoint 1 — Add MR GitLab CI job for `make ci` (`5e8bb47`)
- [x] Checkpoint 2 — Add tag pipeline build/upload jobs to Generic Package Registry (`bd7947e`)
- [x] Checkpoint 3 — Add GitLab Release creation job linking uploaded assets (`b35c3a9`)
- [x] Checkpoint 4 — Update `README.md` with private GitLab download snippet + notes (`2503984`)
- [x] Final checkpoint — PR wrap-up (final approval gate)

## Completion checklist

- [x] Set frontmatter `status: completed`
- [x] Check off completed checkpoint(s) above and add PR/commit references as needed
- [x] Check off all quality gates below (or document any exceptions)

## Quality gates / test plan

- [x] `make ci` (local)
- [x] Merge request pipeline runs and reports `make ci` success (GitLab MR `!1`, pipeline `2317291368`)
- [x] Tag pipeline for `agent-gov/test/v*` uploads 3 binaries + creates a GitLab Release with correct links (tag `agent-gov/test/v0.0.0-gitlab-ci-20260210081012`, pipeline `2317313067`)

## Notes / open questions

- MR: `https://gitlab.com/bsmith.quanata/agent-governance-strategy/-/merge_requests/1`
- Verified `CI_JOB_TOKEN` can upload to the GitLab Generic Package Registry on GitLab.com shared runners (uploads returned `201 Created`).

