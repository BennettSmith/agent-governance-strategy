---
branch: "chore/app-agent-gov-release-pipeline"
status: active
---

# Plan: Build and publish `agent-gov` release binaries on tag

## Summary

Add CI automation that builds and publishes `agent-gov` binaries when a tool tag like `agent-gov/vX.Y.Z` is pushed, so target repos can pin and download a deterministic CLI version (Option A).

## Constraints

- Follow repository governance (`Non-Negotiables.md`, `Architecture.md`, `Constitution.md`).
- Preserve the existing quality gates (`make ci` must pass; coverage >= 85%).
- Keep the release process deterministic and auditable.
- Minimize operational complexity for target repos (predictable asset naming).

## Scope

### In scope

- Add CI workflow(s) to:
  - Trigger on tags matching `agent-gov/v*`.
  - Build `agent-gov` for a small OS/arch set (at least `darwin/arm64`, `darwin/amd64`, `linux/amd64`).
  - Publish binaries as release assets with stable names: `agent-gov_<os>_<arch>`.
  - Create (or update) a GitHub Release for the tag.
- Document the release process briefly (either in `README.md` or a dedicated doc).

### Out of scope

- Changing governance content distribution (`gov/vYYYY.MM.DD`) beyond documenting conventions.
- Implementing signing/notarization (can be a follow-up).
- Adding a full packaging ecosystem (Homebrew, apt, etc.).

## Approach

- Prefer a simple tag-triggered workflow (GitHub Actions) with a build matrix.
- Build from `tools/gov/`:
  - `go build -trimpath -ldflags "-s -w" -o dist/<asset> ./cmd/agent-gov`
- Upload artifacts to the GitHub Release corresponding to the tag.
- Keep asset naming aligned with the target repo Makefile snippet documented in `README.md`.

## Checkpoints

- [x] Checkpoint 1 — Add plan + decide exact tag trigger and asset naming.
- [x] Checkpoint 2 — Add CI workflow for building multi-platform binaries on `agent-gov/v*` tags.
- [ ] Checkpoint 3 — Validate with a dry-run tag in a fork/test repo; document the release procedure.

## Progress log

- Implemented tag trigger `agent-gov/v*` and asset naming `agent-gov_<os>_<arch>` via GitHub Actions workflow.
- Added a safe test tag pattern `agent-gov/test/v*` that creates draft prereleases for in-repo validation.
- Added a short maintainer section to `README.md` describing how to cut an `agent-gov` release tag and expected assets.
- Updated the workflow to build all target OS/arch on `ubuntu-22.04` via Go cross-compilation (avoids reliance on hosted `macos-13` runner availability).

## Completion checklist

- [ ] Set frontmatter `status: completed`
- [ ] Check off completed checkpoint(s) above and add PR/commit references as needed
- [ ] Check off all quality gates below (or document any exceptions)

## Quality gates / test plan

- [x] `make ci`
- [ ] Tag-triggered workflow produces release assets for all supported OS/arch

## Notes / open questions

- Should tool tags (`agent-gov/vX.Y.Z`) and governance tags (`gov/vYYYY.MM.DD`) be allowed to point to the same commit, or should we enforce separation?
- Do we want to publish checksums (e.g., `sha256sum.txt`) alongside binaries?
