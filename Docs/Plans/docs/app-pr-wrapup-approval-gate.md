---
branch: "docs/app-pr-wrapup-approval-gate"
status: completed
---

# Plan: Make PR wrap-up + final approval gate explicit

## Summary

Refine governance guidance so that “PR wrap-up” is an explicit, repeatable final step: the agent updates the branch plan to `status: completed`, adds PR/commit references, re-runs quality gates, then **stops and waits** for explicit human approval before creating the final wrap-up commit.

## Constraints

- Follow repository governance (`Non-Negotiables.md`, `Architecture.md`, `Constitution.md`).
- Keep changes small and reviewable.
- Preserve the existing rules:
  - human acceptance required before checkpoint commits
  - agents must not merge
- `make ci` must pass.

## Scope

### In scope

- Update plan guidance/templates so PR wrap-up is explicit:
  - Add a standard last checkpoint (e.g., “PR wrap-up”) to `Governance/Templates/Plans/Plan.Template.md`.
  - Clarify what “wrap-up” includes: set plan frontmatter to `status: completed`, check off checkpoints and quality gates, add PR/commit references.
- Update agent behavior guidance to reduce ambiguity at the end of a PR:
  - Add an explicit “final approval request” expectation (including a recommended fixed phrase) before the wrap-up commit.
- Sync generated docs into:
  - root scope (`.governance/config.yaml` → `Non-Negotiables.md`, `Constitution.md`, templates)
  - embedded `tools/gov` scope (`tools/gov/.governance/config.yaml`) if applicable.

### Out of scope

- Implementing automation/hooks that enforce plan completion.
- Changing the `agent-gov` CLI behavior.

## Approach

- Prefer a minimal update to the plan template and core agent rules that:
  - makes the PR wrap-up step unmissable
  - keeps the existing “approval before checkpoint commits” rule as the enforcement mechanism
- Keep wording concrete and testable (what to update, when to stop, and what to ask the human for).

## Checkpoints

- [x] Checkpoint 1 — Update governance sources (template + core rules) to describe PR wrap-up and final approval gate. (commit: `f902736`)
- [x] Checkpoint 2 — Run governance sync for root + embedded scopes so generated docs reflect the changes. (commit: `1c44cf5`)
- [x] Checkpoint 3 — Run `make ci`, update this plan with references, and prepare PR. (make ci: ✅, PR: `#10`)

## Completion checklist

- [x] Set frontmatter `status: completed`
- [x] Check off completed checkpoint(s) above and add PR/commit references as needed
- [x] Check off all quality gates below (or document any exceptions)

## Quality gates / test plan

- [x] `make ci`

## Notes / open questions

- What fixed approval phrase should we standardize on for the final wrap-up commit (if any)?
- Should we recommend putting PR URL into the plan as part of wrap-up, or treat it as optional?
