---
branch: "docs/app-plan-completion-gate"
status: completed
---

## Summary

Reduce missed “plan completion” updates by making plan completion an explicit agent requirement and by adding a completion checklist to the plan template.

## Constraints

- Keep changes small and reviewable.
- `make ci` must pass.
- Update governance sources (under `Governance/`) and sync generated docs in-place.

## Scope

### In scope

- Update `Governance/Core/NonNegotiables.Core.md` to require updating the branch plan to “completed” as part of wrapping up an approved task.
- Update `Governance/Templates/Plans/Plan.Template.md` with an explicit completion checklist section.
- Run governance sync for:
  - root scope (`.governance/config.yaml`)
  - embedded `tools/gov` scope (`tools/gov/.governance/config.yaml`)

### Out of scope

- Enforcing this via tooling/hooks (future enhancement).

## Approach

- Add a concise, testable rule under “Agents (core)” that gates “done”/final push on plan completion updates.
- Add a completion checklist section to the plan template with status/checkpoint/quality gate/PR+commit links.
- Sync generated docs so the key-three in this repo reflects the updated governance immediately.

## Checkpoints

- [x] Checkpoint 1 — Update governance sources, sync generated docs, and open PR.
  - PR: `https://github.com/BennettSmith/agent-governance-strategy/pull/5`
  - Commits: `1053be6`, `d8b0cf8`, `32b2205`

## Quality gates / test plan

- [x] `make ci`

## Notes / open questions

- None.
