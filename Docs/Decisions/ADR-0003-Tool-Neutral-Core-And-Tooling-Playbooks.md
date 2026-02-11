---
id: "ADR-0003"
title: "Keep core governance tool-neutral; put tool-specific steps in playbooks"
status: proposed # proposed | accepted | deprecated | superseded
date: 2026-02-10
deciders:
  - "Engineering"
consulted:
  - ""
informed:
  - ""
tags:
  - "governance"
  - "tooling"
related:
  - "Non-Negotiables.md"
  - "Architecture.md"
  - "Constitution.md"
  - "Governance/README.md"
---

## Context

Governance documents in target repositories are generated from this repo’s sources (`Governance/`) using managed blocks and deterministic sync. Projects may add local addenda, but when those addenda represent broadly applicable governance rules, keeping them “local” creates drift and inconsistent expectations across repositories.

Separately, some governance rules require operational actions (for example, setting/updating a PR/MR description). If core governance embeds specific tools (e.g., GitLab’s `glab`), the core becomes less portable across hosting providers and teams that use different toolchains.

We need a consistent approach that:

- keeps core governance portable across Git hosts
- still provides precise operational guidance where needed
- avoids re-introducing tool names into core docs over time

## Decision drivers

- Portability across Git hosts (GitHub, GitLab) and tooling preferences
- Clear separation of “what must be true” vs “how to do it”
- Maintainability: changes to tool commands should not require redefining core governance
- Deterministic generation/sync of managed blocks

## Considered options

1. **Embed tool-specific commands in core governance** — e.g., require `glab` in `Non-Negotiables.md`
2. **Keep core tool-neutral; put tool-specific steps in playbooks** — core states required outcomes, playbooks show how
3. **Keep all of this as local addenda** — allow each repo to decide independently

## Decision

We will keep **core governance tool-neutral** and place **tool-specific operational guidance** in playbooks.

- Core fragments in `Governance/Core/` may state required outcomes (e.g., “PR/MR description must be meaningful and verified”) but must not mandate or name a specific CLI/tool.
- Playbooks may name tools and include concrete commands (e.g., `glab`, `gh`) and will be emitted by profiles into target repositories under `Docs/Playbooks/`.
- When a repo type needs platform-specific operational guidance, the profile (or a profile overlay) should emit the relevant playbooks.

## Consequences

### Positive

- Core governance remains portable and consistent across Git hosts
- Tooling changes are localized to playbooks (smaller diffs, clearer intent)
- Projects can adopt different toolchains without weakening core expectations

### Negative / trade-offs

- Requires profiles to emit the relevant playbooks (slightly more wiring)
- Readers may need to consult both core docs and playbooks for end-to-end guidance

### Risks and mitigations

- **Risk**: Tool-neutral core becomes vague and less enforceable.
  - Mitigation: core rules must specify concrete outcomes and verification requirements; playbooks supply exact commands.
- **Risk**: Tool guidance drifts between GitHub/GitLab playbooks.
  - Mitigation: keep playbooks minimal and focused on the same outcome; prefer shared structure and review together.

## Rationale

Option 2 best matches the builder architecture: core fragments represent universal governance, while playbooks are explicitly intended for operational, platform/tooling guidance. This keeps the governing “contract” stable while allowing tooling to vary.

## Implementation notes

- Promote widely applicable local addenda into `Governance/Core/NonNegotiables.Core.md`.
- Replace tool-specific mentions in core with tool-neutral outcome requirements.
- Add GitHub/GitLab CLI playbooks under `Governance/Core/Playbooks/` and emit them via profiles to `Docs/Playbooks/`.
- Remove redundant local addenda from generated `Non-Negotiables.md` in this repo after syncing.

## Change log

- 2026-02-10 — Proposed
