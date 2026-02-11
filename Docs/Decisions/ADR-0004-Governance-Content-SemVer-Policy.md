---
id: "ADR-0004"
title: "Version governance content with SemVer tags (gov/vX.Y.Z)"
status: accepted # proposed | accepted | deprecated | superseded
date: 2026-02-11
deciders:
  - "Engineering"
consulted:
  - ""
informed:
  - ""
tags:
  - "governance"
  - "versioning"
related:
  - "README.md"
  - "Non-Negotiables.md"
  - "Architecture.md"
  - "Constitution.md"
  - "Docs/Plans/docs/app-governance-approach-docs.md"
---

## Context

Target repositories consume governance content from this repo via `.governance/config.yaml` (`source.repo`, `source.ref`, `source.profile`). For operational reliability and auditability, target repos should be able to pin governance content to an immutable, reviewable identifier.

This repo also ships a separate, vendorable CLI tool (`agent-gov`) that evolves on its own cadence. Teams need a clear model for:

- what gets versioned together vs independently
- how to communicate compatibility and expected remediation
- what constitutes a breaking change for governance content

## Decision drivers

- Repeatability: target repos must be able to re-sync/verify deterministically.
- Clarity: distinguish between “tool changes” and “governance content changes”.
- Auditability: be able to answer “what governance produced these managed blocks?”
- Operability: upgrades should be predictable and support safe rollback.

## Considered options

1. **Date-based tags** for governance content (e.g., `gov/vYYYY.MM.DD`)
2. **Commit SHA only** pins (no tags)
3. **SemVer tags** for governance content (e.g., `gov/vX.Y.Z`) and separate SemVer tags for the tool (e.g., `agent-gov/vX.Y.Z`)

## Decision

We will version **governance content** with SemVer tags `gov/vX.Y.Z`.

- The `gov/vX.Y.Z` tag represents **a single governance release** covering **core + all profiles** at that point in time.
- The `agent-gov` tool continues to use separate SemVer tags `agent-gov/vX.Y.Z`.
- Target repos should pin:
  - tool version via their tool pin (e.g., `AGENT_GOV_TAG=agent-gov/vX.Y.Z`)
  - governance content via `.governance/config.yaml` `source.ref: gov/vX.Y.Z` (or an immutable commit SHA)

## Consequences

### Positive

- Clear mental model: tool and content are distinct, each pinned independently.
- Governance upgrades become communicable (“upgrade from `gov/v1.2.0` to `gov/v1.3.0`”).
- SemVer provides a shared language for breaking vs non-breaking changes.

### Negative / trade-offs

- A single governance content version spans all profiles, so a breaking change in any supported profile requires a major bump for the whole governance release line.
- Maintainers must be disciplined about applying breaking-change criteria consistently.

### Risks and mitigations

- **Risk**: Consumers misinterpret MINOR/PATCH bumps as “no action required”.
  - Mitigation: define explicit breaking-change criteria and keep upgrade playbooks short and procedural.
- **Risk**: Profile-specific changes unintentionally break another profile’s consumers.
  - Mitigation: treat “breaking for any supported profile” as breaking for the governance release; keep changes small and verified.

## Rationale

SemVer tags strike the best balance between stability and operability. They are more informative than SHAs alone and clearer about compatibility expectations than date stamps, while still allowing target repos to pin immutably.

## Implementation notes

- Document the tool/content pin split in `README.md`.
- Prefer immutable tags for `source.ref` in `.governance/config.yaml`; allow commit SHAs as an escape hatch.
- In CI for target repos, prefer verify-only wiring (`agent-gov verify` / `make gov-ci`) and avoid sync-on-CI behavior.

## Breaking change criteria (governance content)

A governance content change requires a **MAJOR** bump if a repo pinned to the prior version is likely to require **human remediation** beyond:

- bumping `source.ref`, and
- re-running `agent-gov sync` to update managed blocks.

Breaking changes include (non-exhaustive):

- New **required repo-side changes** beyond managed-block updates (new required files/paths/layout).
- New or stricter **required quality gates** likely to fail existing repos by default.
- Managed-block contract changes (markers/prefix semantics) that invalidate existing managed blocks or local addenda preservation expectations.
- Template contract changes that force migration of existing docs to remain compliant.
- Profile meaning changes where an existing conforming repo becomes non-conforming without repo-side changes.
- Because governance content is versioned as one release covering core + all profiles: breaking for **any** supported profile is breaking for the governance release.

## Change log

- 2026-02-11 — Accepted

