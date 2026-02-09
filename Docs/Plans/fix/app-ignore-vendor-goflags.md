---
branch: "fix/app-ignore-vendor-goflags"
status: active
---

## Summary

Ensure `make ci` succeeds even when a developer machine has `GOFLAGS=-mod=vendor` set globally, by overriding `GOFLAGS` for the `tools/gov` module’s Go commands.

## Constraints

- Keep the change small and reviewable.
- `make ci` must pass with no failures.

## Scope

### In scope

- Update the root `Makefile` to explicitly set `GOFLAGS` for `tools/gov` `go` commands.
- Add a plan document for this branch.

### Out of scope

- Changing developer machine configuration or dotfiles.
- Introducing or committing a `vendor/` directory for `tools/gov`.

## Approach

- Add a `Makefile` variable (defaulting to `-mod=mod`) used only for `tools/gov` subprocesses.
- Prefix `go run`, `go test`, and `go build` invocations under `tools/gov` with `GOFLAGS="$(...)"`.
- Verify by running `make ci` with `GOFLAGS=-mod=vendor` set in the environment.

## Checkpoints

- [ ] Checkpoint 1 — Add plan + Makefile override and verify `make ci` under `GOFLAGS=-mod=vendor`.

## Quality gates / test plan

- [ ] `GOFLAGS=-mod=vendor make ci`

## Notes / open questions

- None.

