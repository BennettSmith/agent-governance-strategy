---
id: "ADR-0001"
title: "Retain feature branch plans in-repo under Docs/Plans"
status: accepted # proposed | accepted | deprecated | superseded
date: 2026-02-05
deciders:
  - "Engineering"
consulted:
  - ""
informed:
  - ""
tags:
  - "process"
  - "documentation"
related:
  - "Docs/Plans/Plan.Template.md"
  - "Constitution.md"
  - "AGENTS.md"
---

## Context

Agents are required to begin work with a written plan and to proceed in small, reviewable steps.
Historically, plans either lived in ephemeral locations (PR descriptions, chat logs) or were deleted
before merge, which reduced long-term traceability and made it harder to understand why changes were
made, what trade-offs were considered, and how work was staged.

We want a lightweight, durable record of intent and checkpoints that:

- is colocated with the codebase
- is easy to find given a feature branch
- does not require a separate system to manage
- supports future audits, onboarding, and follow-up work

Assumptions:

- Not every plan warrants a full decision record; MADRs are reserved for system-shape and long-term constraints.
- Plans should be discoverable by branch name and remain small/operational.

In scope:

- Where plans live, naming, and the minimal workflow

Out of scope:

- CI enforcement details (tooling-specific)
- PR templates or Git hosting configuration

## Decision drivers

- Traceability of intent and checkpoints
- Low friction for authors and reviewers
- Consistent naming and discoverability
- Minimal governance overhead (plans vs ADRs)

## Considered options

1. **Keep plans only in feature branches** — create `Plan.md` in repo root and delete/reset before merge
2. **Keep plans in PR descriptions only** — rely on the hosting platform for plan history
3. **Retain plans in-repo** — store plans under `Docs/Plans/` and merge them to main

## Decision

We will retain feature branch plans in-repo by storing them at `Docs/Plans/<branch-name>.md` and
merging them to `main` as part of the branch.

The plan filename must match the feature branch name. If the branch name contains `/`, it is mirrored
as subfolders under `Docs/Plans/` (e.g. `feat/identity-add-foo` → `Docs/Plans/feat/identity-add-foo.md`).

Branch names follow `<type>/<area>-<short-slug>`, where `<area>` is usually a bounded context slug. For work that is
not related to a bounded context (e.g., composition root / app wiring / routing), use `app` as the `<area>` segment.

## Consequences

### Positive

- Durable record of intent and staging/checkpoints
- Easier review context and future archaeology without scraping PR comments
- Simple lookup: branch name → plan path

### Negative / trade-offs

- More documents in the repo (noise) and occasional stale content
- Requires light discipline to keep plans concise and updated

### Risks and mitigations

- **Risk**: Plans accumulate and become hard to navigate.
  - Mitigation: enforce naming convention and keep plan scope operational; use ADRs for long-term decisions.
- **Risk**: Plans are treated as canonical specs and drift from reality.
  - Mitigation: plans are working documents; final state is reflected in code, tests, and ADRs/use cases.

## Rationale

Retaining plans in `Docs/Plans/` strikes the best balance between traceability and operational
friction. It avoids dependence on a specific Git host, keeps planning artifacts close to the code,
and preserves useful context without requiring every change to produce an ADR.

## Implementation notes

- Add `Docs/Plans/Plan.Template.md`
- Require plans in `Constitution.md` and provide usage guidance in `AGENTS.md`
- (Optional, follow-up) Add CI checks to ensure a plan exists for a branch and matches naming rules

## Alternatives (details)

### Keep plans only in feature branches

- Pros:
  - Keeps `main` cleaner
- Cons:
  - Loses long-term context; reviewers and future maintainers must hunt through PRs or local branches

### Keep plans in PR descriptions only

- Pros:
  - No repo artifacts
- Cons:
  - Platform-dependent, harder to search/version, and not reliably available offline

## Links

- Related Use Case(s): <!-- n/a in this builder repo; templates (if any) are profile-emitted from `Governance/Templates/` -->
- Related Domain doc(s): <!-- n/a in this builder repo; templates (if any) are profile-emitted from `Governance/Templates/` -->
- Related ADR(s): `Docs/Decisions/`
- External references: <!-- -->

## Change log

- 2026-02-05 — Proposed

