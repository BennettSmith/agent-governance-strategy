## Quality gates (core)

- Tests must be written before production code.
- Code coverage must remain above 85%.
- Generated code (openapi) is not tested. 
- Generated code (openapi) is excluded from code coverage requirement.
- `make ci` must pass with no failures.
- Quality gates must pass for a task to be considered complete.

## Legacy & refactoring (core)

- A change is a **refactor** only if it makes **no externally observable behavior change** at defined boundaries (persistence contracts, API interactions, user-visible UI behavior, and any other profile-defined external boundary).
- Refactors must **not** be mixed with behavior changes (features/fixes) in the same branch or pull request. If both are needed, do the refactor first, then do the behavior change as a separate step.
- Before changing legacy internals, preserve current behavior with a **safety net** (characterization tests, contract tests, or golden-master tests) at the relevant boundary.
- Refactoring steps must be **small and reversible**. Each step must:
  - keep the prior behavior intact
  - pass all quality gates
  - be trivially revertible (small diff, or guarded by a flag, or a parallel-change/expand-contract approach where the old path remains available until cutover)
- New code introduced during refactoring must follow the target architecture; legacy code may remain only when **isolated behind explicit boundaries** (seams/adapters/facades/boundaries) so violations do not spread.
- Follow `Docs/Refactoring/Legacy-Refactoring-Playbook.md` for approved legacy refactoring patterns and checklists.

## Agents (core)

- Agents must follow the architecture for the target profile.
- All changes must be made in feature branches.
- Agents must request human acceptance before checkpoint commits.
- Agents must not mark plan tasks complete without approval.
- After an approved task is complete, agents must update the branch plan to reflect completion (e.g., set `status: completed`, check off completed checkpoints and quality gates, and add PR/commit references) before declaring the work “done” or pushing final updates.
- Agents must not merge branches under any circumstances.
