# Governance exceptions

This playbook describes how to record and manage **exceptions** to governance rules in a target repository.

Core policy (already decided in this repo): **exceptions live in MADRs** under `Docs/Decisions/`. Do not introduce a separate “exceptions registry”.

Important: the **exception MADR is created in the target repo** (the repo taking on the exception), not in this governance repo. This governance repo defines the policy and provides the playbook/template; each target repo records its own exception decisions locally.

## What counts as an exception

An exception is any intentional deviation from centrally defined governance expectations that would otherwise be enforced or treated as required in the target repo. Examples:

- temporarily allowing a dependency / import that would normally be forbidden by architecture rules
- deferring a required quality gate in CI for a bounded period while remediation work is underway
- adopting a profile but opting out of a specific emitted rule due to platform constraints

If it’s just repo-owned notes, clarifications, or additional local rules, it likely belongs in **Local Addenda (project-owned)** rather than an “exception”.

## When exceptions are permissible

Exceptions should be rare. Prefer to:

- comply with the governance rule as-is
- implement the rule in stages (tighten over time)
- isolate the constraint behind an explicit boundary (so violations don’t spread)

Use an exception when there is a concrete constraint (technical, organizational, regulatory, platform) and a clear, time-bounded plan to return to the baseline.

## Required record: MADR

Create a decision record using the template at `Docs/Decisions/MADR.Template.md`.

In the MADR, include (at minimum) the following information (put it in the most appropriate section(s), e.g. Context / Decision / Risks and mitigations / Implementation notes):

- **Rule being excepted**: link to the relevant governance doc / section (or quote the exact requirement).
- **Scope**: what repo(s), module(s), path(s), or component(s) the exception applies to.
- **Owner**: a person/team accountable for remediation and follow-through.
- **Justification**: the concrete constraint that makes compliance infeasible right now.
- **Mitigations**: what compensating controls reduce risk while the exception exists.
- **Expiry / review date**: a specific date or milestone when the exception must be reviewed, renewed, or removed.
- **Exit criteria**: what “back to compliant” means and how it will be verified.

## How exceptions relate to “Local Addenda (project-owned)”

- **Local Addenda** is for project-owned notes and adaptations that are expected to persist without requiring central governance changes.
- **Exceptions** are explicit, intentional deviations from governance expectations and must be captured in a MADR so they are reviewable, time-bounded, and owned.

If you reference an exception in Local Addenda, treat it as a pointer only (e.g., “See ADR-0123”), not the source of truth.

## Operational checklist (target repo)

1. Create a new MADR under `Docs/Decisions/` describing the exception.
2. Make the minimal code/config change needed to implement the exception.
3. Add any compensating checks (mitigations) to reduce risk.
4. Ensure the repo remains green on its quality gates (typically `make ci`) and that governance verification is still run (even if it allows an exception by design).
5. Track the review/expiry date and remove the exception promptly once exit criteria are met.
