---
id: "ADR-0002"
title: "Adopt governance-by-profiles builder with managed-block sync"
status: proposed # proposed | accepted | deprecated | superseded
date: 2026-02-06
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
  - "Docs/Plans/feat/app-governance-profiles-builder.md"
---

## Context

This repository currently encodes governance in a single, iOS-leaning set of documents (e.g., references to SwiftUI and Swift Packages). We want a scalable approach that can:

- Support multiple repo types (mobile apps, backend services, web apps) with different architectural profiles.
- Keep core working agreements and principles consistent across repos.
- Allow per-repo/platform drift in implementation details without losing alignment on core governance.
- Provide deterministic, non-destructive re-sync of the “core-managed” parts of governance docs in target repos.

## Decision drivers

- Deterministic updates that do not clobber local addenda
- Ability to support multiple repo types/profiles over time
- Low-friction adoption (works inside target repos, minimal dependencies)
- Auditability (which governance version produced a given doc)
- Offline-friendly in practice (cacheable), but not strictly offline-only

## Considered options

1. **Single platform-neutral key-three + playbooks** in one repo
2. **Per-platform/per-repo copies** of governance docs (manual duplication)
3. **Builder repo with profiles** that generates governance bundles into target repos, with managed-block re-sync

## Decision

We will adopt **option 3**: a **governance-by-profiles builder** strategy.

- This repo will define governance as **Core + Profile + Platform overlay/playbooks + Templates**.
- A Go CLI tool, **`agent-gov`**, will be vendored into target repos and will support:
  - `init`, `sync`, `verify`, `build`
  - reading profile content from a pinned governance repo tag/release
  - caching fetched source using `os.UserCacheDir()` + `/govbuilder`
- Generated docs in target repos will use **managed blocks** with stable markers:
  - `<!-- GOV:BEGIN ... -->` / `<!-- GOV:END ... -->`
  - markers include `sourceRepo`, `sourceRef`, `sourceCommit` for auditability
  - a required “Local Addenda (project-owned)” section remains editable and is never overwritten

v1 profiles:

- `mobile-clean-ios`
- `backend-go-hex`

## Consequences

### Positive

- One upstream source of truth for core governance with deterministic re-sync.
- Profile-specific governance can evolve independently while sharing core.
- Target repos can run governance operations locally via the vendored CLI.

### Negative / trade-offs

- Additional complexity: fragments, profiles, and a CLI tool to maintain.
- Remote fetching requires git access/credentials and may fail without network (unless cached).

### Risks and mitigations

- **Risk**: Profiles diverge unintentionally in core principles.
  - Mitigation: keep core fragments small and shared; CI verifies managed blocks via `agent-gov verify`.
- **Risk**: Tag/release ref is moved.
  - Mitigation: record resolved commit SHA and include it in markers/metadata; prefer immutable release tags.

## Rationale

The builder approach preserves developer ergonomics by generating platform/profile-specific outputs, while still maintaining alignment on core governance principles. Managed-block sync provides a safe mechanism to update core content without overwriting local adaptations.

## Implementation notes

- Add `Governance/` fragment and profile layout.
- Build `tools/gov/` (`agent-gov`) with explicit test coverage goals (>= 85% for CLI package surface).
- Add a root `Makefile` with a canonical `make ci` target for CI status checks.

## Change log

- 2026-02-06 — Proposed
