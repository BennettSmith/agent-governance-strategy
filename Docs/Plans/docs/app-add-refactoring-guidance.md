# Summary
This change adds explicit guidance for refactoring legacy code in an orderly, reversible, small-step way, without weakening the target architecture constraints.

It introduces a refactoring playbook and tightens the governing documents to clearly separate refactors from behavior changes.

## Constraints

- `Non‑Negotiables.md` overrides all other documents.
- Guidance must preserve Clean Architecture boundaries, use-case primacy, and offline-first expectations.
- Refactoring must be treated as a distinct class of work, with clear reversibility requirements and safety nets.

## Scope

### In scope

- Add explicit “refactor vs behavior change” definitions and non-mixing rule to governing docs.
- Add a refactoring protocol focused on small, reversible steps (characterization tests, seams, parallel change).
- Introduce `Docs/Refactoring/Legacy-Refactoring-Playbook.md` referenced by the governing docs.

### Out of scope

- Adding or modifying application use cases, domain documents, or decisions outside refactoring guidance.
- Tooling/CI implementation (e.g., automated enforcement) beyond documentation.

## Approach

- Add a new `Legacy & refactoring` section to `Non‑Negotiables.md` with enforceable constraints.
- Add a `Legacy refactoring protocol` section to `Constitution.md` describing how refactor work is planned and executed.
- Add a small cross-reference in `Architecture.md` to direct legacy modernization work to the playbook.
- Create a playbook doc with approved patterns, checklists, and examples.

## Checkpoints

- [ ] Checkpoint 1 — Add plan + playbook scaffold and cross-links
- [ ] Checkpoint 2 — Update `Non‑Negotiables.md` with enforceable refactor rules
- [ ] Checkpoint 3 — Update `Constitution.md` and `Architecture.md` with protocol + references

## Quality gates / test plan

- [ ] All markdown links resolve and referenced files exist
- [ ] Wording is consistent with precedence rules and does not conflict with existing constraints

## Notes / open questions

- This repository is not currently a git repo; in a normal workflow, this plan would be committed on branch `docs/app-add-refactoring-guidance`.

