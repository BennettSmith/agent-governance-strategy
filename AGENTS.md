# AGENTS

This repository is governed by the documents below. Agents must follow them when planning or making changes.

## Read-first (governing docs)

Read governing docs before starting work.

### Scope selection (embedded repos)

- **Default**: Use the repository root key-three.
- **Embedded scope**: If you are working under a subtree that contains a nearer `AGENTS.md` (for example `tools/gov/AGENTS.md`), follow that scope’s instructions, including reading that scope’s key-three first.
- If anything is ambiguous, stop and ask for human direction.

### Root key-three

1. `Non-Negotiables.md`
2. `Architecture.md`
3. `Constitution.md`

### Precedence (conflicts)

- `Non-Negotiables.md` overrides all other documents
- `Architecture.md` overrides `Constitution.md` on matters of system shape
- `Constitution.md` guides behavior when other documents are silent
- If anything is ambiguous, stop and ask for human direction

## Working style (how to proceed)

- Start with a written plan (markdown) before implementation.
- Prefer small, reviewable changes over large refactors.
- If unsure, propose options and ask for a decision; do not make irreversible choices.
- Follow offline-first and explicit-boundary expectations described in the governing docs.

## Plans (required)

- Plans are kept in `Docs/Plans/` and are committed as part of the feature branch history.
- Create a plan at `Docs/Plans/<branch-name>.md` by copying `Docs/Plans/Plan.Template.md`.
- The plan filename must match the feature branch name.
  - If the branch contains `/`, mirror it as subfolders under `Docs/Plans/`.
  - Example: branch `feat/identity-add-session-restore` → plan `Docs/Plans/feat/identity-add-session-restore.md`
- Update the plan as work progresses (especially at checkpoints).

## Branch naming (required)

Feature branches must follow:

- `<type>/<area>-<short-slug>` (lowercase, hyphenated)
- Allowed `<type>` values: `feat`, `fix`, `docs`, `chore`, `refactor`, `spike`
- `<area>` is usually a bounded context slug (e.g., `identity`)
- Use `app` as the `<area>` segment when the change is not related to a bounded context (e.g., composition root / app wiring / routing)
- Other `<area>` values may be used when appropriate, but should be rare and must be documented in the plan

Examples:

- `feat/identity-add-session-restore`
- `fix/sync-handle-409-conflict`
- `docs/identity-update-usecase-catalog`
- `refactor/app-rework-root-router`

## Where things live

- **Architecture overview**: `Architecture.md`
- **Decision records (MADR)**: `Docs/Decisions/` (template: `Docs/Decisions/MADR.Template.md`)
- **Plans**: `Docs/Plans/` (template: `Docs/Plans/Plan.Template.md`)
- **Governance source of truth**: `Governance/` (core fragments, profiles, templates, and playbooks)
- **Archived examples (builder repo only)**: `Docs/Archive/`

## When to create new docs

- If a change affects system shape or long-term constraints, capture a decision record in `Docs/Decisions/`.
- If adding/changing governance behavior in this repo, prefer updating the appropriate `Governance/` source (fragments, profiles, templates, playbooks) and capture decisions in `Docs/Decisions/` when they affect long-term constraints.
