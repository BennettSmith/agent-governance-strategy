---
branch: "chore/app-ci-make-ci-gates"
status: active
---

# Plan: Require `make ci` for PRs and releases

## Summary

Ensure GitHub runs `make ci` for every PR update and requires it to pass before publishing `agent-gov` release assets on tag builds.

## Constraints

- Follow repository governance (`Non-Negotiables.md`, `Architecture.md`, `Constitution.md`).
- Keep changes small and reviewable.
- `make ci` must pass.

## Scope

### In scope

- Add a GitHub Actions workflow that runs `make ci` on `pull_request` updates.
- Update the tag-triggered release workflow so it runs `make ci` first and fails the release if `make ci` fails.

### Out of scope

- Changing the `make ci` target behavior.
- Adding additional linters or test suites beyond what `make ci` already runs.
- Enforcing branch protection rules (GitHub settings) in this repo (doc-only note if needed).

## Approach

- Keep CI logic centralized by calling `make ci` from workflows.
- Use separate jobs with `needs:` so the release job cannot run unless the CI job succeeds.

## Checkpoints

- [x] Checkpoint 1 — Add PR CI workflow to run `make ci` on PR updates. (commit: `25554f2`)
- [x] Checkpoint 2 — Gate tag release workflow on a successful `make ci` job. (commit: `859c821`)
- [x] Checkpoint 3 — Run `make ci`, push branch, and open PR. (make ci: ✅, PR: `#11`, CI fix: `b3f9544`)
- [ ] Checkpoint 4 — PR wrap-up (final approval gate): update completion checklist + set `status: completed`, then create the final wrap-up commit.

## Completion checklist

- [ ] Set frontmatter `status: completed`
- [ ] Check off completed checkpoint(s) above and add PR/commit references as needed
- [ ] Check off all quality gates below (or document any exceptions)

## Quality gates / test plan

- [ ] `make ci`

## Notes / open questions

- None.

