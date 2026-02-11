---
branch: "docs/app-docs-folder-cleanup"
status: completed
---

# Summary

Clean up stale documentation in this governance builder repository’s `Docs/` folder, while preserving required governance artifacts and aligning refactoring playbooks with the `Governance/` source-of-truth model.

## Constraints

- Do not remove feature-branch plans under `Docs/Plans/`.
- Keep `Docs/Refactoring/Legacy-Refactoring-Playbook.md` at a stable path across profiles (Pattern 1).
- Emit refactoring playbooks as managed documents (`documents:`) so `agent-gov sync` and `agent-gov verify` keep them current (Option A).

## Scope

### In scope

- Archive stale builder-repo examples in `Docs/Domains/` and `Docs/UseCases/`.
- Update references (e.g. `AGENTS.md`, `Docs/Decisions/ADR-0001-*`) that point at removed/archived paths.
- Introduce per-profile refactoring playbook sources under `Governance/Profiles/*/Playbooks/` and emit them as managed documents to `Docs/Refactoring/Legacy-Refactoring-Playbook.md`:
  - `docs-only`: docs/tooling refactoring guidance for this builder repo
  - `mobile-clean`: clean-architecture + use-case-driven guidance for mobile targets
  - `backend-go-hex`: Go hexagonal refactoring guidance for service/CLI targets
- Regenerate docs using `agent-gov init` and ensure `make ci` passes.

### Out of scope

- Changing the underlying decision that feature-branch plans are retained in-repo.
- Removing the `Docs/Refactoring/Legacy-Refactoring-Playbook.md` reference from core governance fragments (we’re keeping the stable path).

## Approach

- Make doc-structure changes small and reversible (move content into `Docs/Archive/` rather than deleting).
- Treat refactoring playbooks as profile-specific managed docs so drift is prevented by normal sync/verify operations.

## Checkpoints

- [x] Checkpoint 1 — Create branch + plan; archive stale docs and update broken links.
- [x] Checkpoint 2 — Add profile-specific managed refactoring playbooks and update manifests.
- [x] Checkpoint 3 — Regenerate docs (`agent-gov init`) and run `make ci`.

## Quality gates / test plan

- [ ] `make ci`

## Notes / open questions

- None

## Progress log

- Stale example docs were moved under `Docs/Archive/` (keeping them for reference while removing them from the “live” docs surface).
- Profile playbooks exist under `Governance/Profiles/*/Playbooks/`, including per-profile `Legacy-Refactoring-Playbook.md`.
- The stable path `Docs/Refactoring/Legacy-Refactoring-Playbook.md` exists in the repo.
