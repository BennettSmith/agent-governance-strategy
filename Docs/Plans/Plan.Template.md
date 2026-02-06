# Plan Template

Use this template to create a plan for a feature branch.

- Copy this file to `Docs/Plans/<branch-name>.md`
- Keep the plan updated as you execute work and create checkpoint commits

> Naming: the plan filename must match the feature branch name. If the branch contains `/`,
> mirror it as subfolders under `Docs/Plans/` (e.g. `feat/identity-add-foo` → `Docs/Plans/feat/identity-add-foo.md`).
>
> Branch names follow `<type>/<area>-<short-slug>`.
> `<area>` is usually a bounded context slug (e.g., `identity`). If the change is not related to a bounded context
> (e.g., composition root / app wiring / routing), use `app` as the `<area>` segment (e.g. `refactor/app-rework-root-router`).

## Summary

<!-- 1-3 sentences describing the user-visible outcome and why we're doing it. -->

## Constraints

<!-- List binding constraints (architecture, offline-first, privacy, etc.). -->

## Scope

### In scope

- <!-- -->

### Out of scope

- <!-- -->

## Approach

<!-- Brief outline of the intended implementation approach. -->

## Checkpoints

<!--
Checkpoints are small, reviewable steps. Update this list as work progresses.
Include links, commit SHAs, or references as useful.
-->

- [ ] Checkpoint 1 — <!-- description -->
- [ ] Checkpoint 2 — <!-- description -->
- [ ] Checkpoint 3 — <!-- description -->

## Quality gates / test plan

<!--
List the checks that must pass for this work to be considered done.
Example: `make ci`, unit tests, documentation updates, etc.
-->

- [ ] <!-- -->

## Notes / open questions

- <!-- -->
