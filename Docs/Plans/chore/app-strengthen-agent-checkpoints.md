---
branch: "chore/app-strengthen-agent-checkpoints"
status: completed
---

## Summary

Strengthen core governance so agents must pause for explicit human approval at each checkpoint **even when no commit is being created**, closing the gap where “approval before checkpoint commits” is insufficient. Provide a governance-owned playbook that makes the checkpoint/approval ritual concrete and repeatable across all profiles.

## Constraints

- Follow repository governance (root key-three): `Non-Negotiables.md` overrides all other docs.
- Prefer small, reviewable checkpoints.
- Do not mix refactors with behavior changes.
- Quality gates must pass: `make ci`.
- Do not create checkpoint commits without explicit human approval.

## Scope

### In scope

- Update core agent rules to make approval gates apply to **checkpoint progression**, not only commits:
  - `Governance/Core/NonNegotiables.Core.md`
  - `Governance/Core/Constitution.Core.md`
- Add a core playbook documenting the checkpoint approval workflow:
  - `Governance/Core/Playbooks/Agent-Checkpoint-Workflow.md`
- Ensure the new playbook is emitted by all profiles:
  - `Governance/Profiles/docs-only/profile.yaml`
  - `Governance/Profiles/backend-go-hex/profile.yaml`
  - `Governance/Profiles/mobile-clean/profile.yaml`

### Out of scope

- Changing profile semantics unrelated to agent checkpoint workflow.
- Tooling changes (e.g., adding new CI checks) beyond what’s needed to keep `make ci` passing.
- Editing generated target-repo docs in-place (we update governance sources and profile wiring here).

## Approach

- Tighten the binding language in core governance to eliminate ambiguity:
  - Explicitly require “stop and request approval” before proceeding past a checkpoint, regardless of whether a commit is created.
  - Keep “approval before checkpoint commits” as an additional safeguard (not the only one).
- Add a governance-owned playbook with:
  - Definitions (checkpoint vs planned task; what “stop” means).
  - A required “APPROVAL REQUEST (checkpoint X)” template that agents must use.
  - Minimal human response expectations (approve / request changes / re-scope).
- Wire the playbook into all profiles so it is consistently emitted into target repos.

## Checkpoints

- [x] Checkpoint 1 — Draft proposed text changes for core rules (NonNegotiables + Constitution) and present for review (no commits yet).
- [x] Checkpoint 2 — Implement core rule changes in `Governance/Core/*.md` and update/format as needed.
- [x] Checkpoint 3 — Add `Governance/Core/Playbooks/Agent-Checkpoint-Workflow.md`.
- [x] Checkpoint 4 — Wire playbook into all profiles (`Governance/Profiles/**/profile.yaml`).
- [x] Checkpoint 5 — Run quality gates (`make ci`) and fix any lint/format fallout.
- [ ] Final checkpoint — PR wrap-up (final approval gate)

## Completion checklist

- [x] Set frontmatter `status: completed`
- [ ] Check off completed checkpoint(s) above and add PR/commit references as needed
- [x] Check off all quality gates below (or document any exceptions)

## Quality gates / test plan

- [x] `make ci`

## Notes / open questions

- Ensure the new wording doesn’t conflict with existing “stop after completing each planned task” language; clarify definitions in the playbook to keep expectations crisp.

## References

- Merge request: `!10`
- Merge request URL: `https://gitlab.com/bsmith.quanata/agent-governance-strategy/-/merge_requests/10`
- Checkpoint commits:
  - Checkpoint 1 (plan created + marked): `9d521c7`
  - Checkpoint 2 (core rules): `9add21f`
  - Checkpoint 3 (playbook): `5e0124f`
  - Checkpoint 4 (profile wiring): `ac502da`
  - Checkpoint 5 (quality gates + formatting): `b2b597c`
