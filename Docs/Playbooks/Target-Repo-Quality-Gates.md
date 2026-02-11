# Target repo quality gates

This playbook describes how teams should implement **repeatable, tool-enforced quality gates** in target repositories, beyond the governance rules themselves.

Governance provides shared expectations (and verification of managed blocks). Target repos still need **repo-owned** enforcement for product code quality, architecture boundaries, and delivery checks.

Related playbooks:

- `Docs/Playbooks/Governance-Upgrades.md`
- `Docs/Playbooks/Governance-Exceptions.md`

## Principles

- **Make it repeatable**: the same command should behave the same on every machine and in CI.
- **Make it fast**: prefer checks that finish quickly and can run on every change.
- **Make it deterministic**: avoid checks that depend on network, time, or nondeterministic ordering.
- **Fail loudly**: errors should be actionable (clear messages, links, and next steps).
- **Prefer static enforcement**: compile-time and static analysis checks are usually more reliable than heuristics.

## Golden commands (repo-owned)

Most teams standardize on:

- `make fmt` — formatting (may write files)
- `make lint` — lint/static analysis (no writes)
- `make test` — tests (no writes)
- `make ci` — the CI gate that runs the full required set (no writes)

Your `ci` target is the contract: if `make ci` is green, the change is acceptable.

## Governance verification fits alongside product checks

Governance verification should be explicit and CI-safe:

- In CI, run **verify only** (no sync / no working tree writes).
- Keep governance checks **additive** to product checks, not a replacement.

Recommended wiring in a target repo:

- `make gov-ci` → runs `agent-gov verify ...` (or `make gov-verify`)
- `make ci` includes `gov-ci` explicitly (e.g., `ci: gov-ci`)

This avoids “attach to ci” magic inside shared includes and keeps product stacks free to define their own `ci` contract.

## What to put in `make ci`

As a starting point, `make ci` often includes:

- `make fmt-check` (or a formatter run in “check” mode)
- `make lint`
- `make test`
- `make gov-ci` (governance verify-only)

Keep `make fmt` available for developers (may write), but keep `make ci` write-free.

## Architecture enforcement

Architecture enforcement should be:

- **explicit** (document what is forbidden and why)
- **mechanical** (enforced by a tool, not by review memory)
- **scoped** (enforce at boundaries you can describe and test)

### Profile example playbooks

Core governance keeps this playbook tool-neutral. Platform- and stack-specific examples are emitted by profiles as additional playbooks.

If your profile emits them, see:

- iOS (SwiftLint examples): `Docs/Playbooks/Quality-Gates-iOS-Architecture-Enforcement.md`
- Go (golangci-lint / depguard examples): `Docs/Playbooks/Quality-Gates-Go-Architecture-Enforcement.md`
- Kotlin/JVM (Detekt / ArchUnit examples): `Docs/Playbooks/Quality-Gates-Kotlin-Architecture-Enforcement.md`

If your stack is not covered by an existing profile, treat that as a signal to add a profile-specific playbook rather than placing stack-specific examples in core.

### Reliability guidance for architecture rules

- Prefer **module boundaries** (targets/packages) for hard isolation; use lint rules as fast feedback.
- Avoid overly clever regexes; choose conventions that are easy to match.
- Minimize false positives; developers will disable noisy rules.
- Provide a clear escape hatch:
  - allow inline disables only with justification
  - record long-lived exceptions via a MADR (see `Docs/Playbooks/Governance-Exceptions.md`)

## Operational guidance

- Start with a small, enforceable set of gates and tighten over time.
- Keep a single source of truth: `make ci` plus a short doc describing what it runs.
- Treat flaky checks as production incidents: fix or remove quickly.
- When introducing a new required gate, provide:
  - a migration plan
  - a clear failure message
  - an owner responsible for keeping it healthy
