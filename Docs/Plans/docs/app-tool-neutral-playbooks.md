---
branch: "docs/app-tool-neutral-playbooks"
status: active
---

## Summary

Move local governance addenda into upstream governance sources, keep core rules tool-neutral, and move Git-host tooling specifics into GitHub/GitLab CLI playbooks.

## Constraints

- `Non-Negotiables.md` overrides all other governing documents.
- Core fragments under `Governance/Core/` must remain tool-neutral/portable across Git hosts.
- Tooling-specific operational instructions belong in playbooks (not core fragments).
- `make ci` must pass.

## Scope

### In scope

- Promote local addenda currently in `Non-Negotiables.md` into `Governance/Core/NonNegotiables.Core.md`.
- Replace tool-specific core guidance with tool-neutral requirements (PR/MR description outcome).
- Add CLI-only playbooks for GitLab (`glab`) and GitHub (`gh`) and emit them via profiles.
- Add ADR documenting “tool-neutral core; tooling in playbooks”.

### Out of scope

- Adding web UI fallback steps (CLI-only playbooks for now).
- Changing any runtime/product architecture content.
- Changing CI enforcement behavior beyond what is needed to emit new docs/playbooks.

## Approach

- Update core governance fragment(s) to include the promoted rules while removing any tool names.
- Add two new playbooks under `Governance/Core/Playbooks/` with operational steps for `glab` and `gh`.
- Update profile manifests to emit the playbooks into `Docs/Playbooks/` in target repos.
- Sync this repo using `.governance/config.yaml` and clean up now-redundant local addenda.
- Verify with `make ci`.

## Checkpoints

- [x] Checkpoint 1 — Add ADR + update core fragment(s) (commit: `b93c30d`)
- [x] Checkpoint 2 — Add playbooks + emit via profiles (commit: `e560baa`)
- [ ] Checkpoint 3 — Sync outputs + clean up local addenda
- [ ] Final checkpoint — PR wrap-up (final approval gate)

## Completion checklist

- [ ] Set frontmatter `status: completed`
- [ ] Check off completed checkpoint(s) above and add PR/commit references as needed
- [ ] Check off all quality gates below (or document any exceptions)

## Quality gates / test plan

- [x] `make ci`

## Notes / open questions

- None.
