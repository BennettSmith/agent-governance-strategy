---
branch: "docs/app-governance-approach-docs"
status: active
---

## Summary

Improve and clarify the documentation for the multi-repo governance approach in this repository so new adopters can quickly understand **what this is**, **why it exists**, **how to apply it**, and **how to evolve it safely** (pinning, upgrades, exceptions, and enforcement).

## Constraints

- Follow repository governing docs: `Non-Negotiables.md`, `Architecture.md`, `Constitution.md`.
- Keep changes **small and reviewable**; focus on docs-first improvements and eliminate ambiguity.
- Preserve the builder architecture framing: governance source-of-truth in `Governance/`, applied via `agent-gov` using managed blocks + local addenda.
- Keep guidance **tool-neutral** where feasible (GitHub/GitLab), but include concrete examples and “golden command” expectations.
- `make ci` must pass before completion (even though this is docs-focused).

## Scope

### In scope

- **Planning source material (temporary, untracked)**
  - Use the local, untracked scratch document `multi-repo-governance-strategy.md` as source material to identify gaps and extract the best ideas.
  - Incorporate adopted ideas into tracked docs (README + playbooks) and/or this plan’s “gaps/open questions”.
  - Delete `multi-repo-governance-strategy.md` once planning/extraction is complete; it must not be committed.

- **Golden commands (execution step; not during planning)**
  - Add canonical targets to the root `Makefile`: `fmt`, `lint`, and `ci` (and any supporting `*-check` targets as needed).
  - Ensure markdown in this repo is consistently **formatted** and **linted** via those golden commands.
  - Decide on a repo-appropriate markdown toolchain and installation strategy (prefer deterministic, low-friction, CI-friendly).
  - Clarify and pin Node runtime expectations for markdown tooling by pinning the Node **major** line (Option B), and documenting it for developers and CI.
  - Document the expected behavior in `README.md` (what each command does) once implemented.
  - Improve target-repo developer UX for governance verification by adding a `gov-ci` target to `Governance/Templates/Make/agent-gov.mk` and documenting the recommended integration:
    - in consuming repos: `ci: gov-ci` (explicit wiring; avoids the include defining stack-specific `ci/fmt/lint/test`)
    - `gov-ci` should remain governance-scoped (e.g., run `gov-verify`) and avoid surprising behavior changes.

- **README improvements**
  - Add an opening section: what this repo is, why we built it, and goals/non-goals.
  - Add/clarify a “mental model” section (builder repo → profiles/fragments → `agent-gov` → managed blocks/local addenda).
  - Clarify **pinning**:
    - tool pin (`agent-gov/vX.Y.Z`)
    - governance content pin (`source.ref` like `gov/vYYYY.MM.DD` or commit SHA)
  - Add a clear recommendation for “golden commands” in target repos (esp. `make ci`) and how governance verification fits.

- **Playbooks**
  - Add `Docs/Playbooks/Governance-Upgrades.md` describing:
    - how to bump pins safely
    - expected verification steps (`verify`, `make ci`)
    - rollback strategy
    - what constitutes a breaking change for governance/tooling
  - Add `Docs/Playbooks/Governance-Exceptions.md` describing:
    - when exceptions are permissible
    - required decision record (MADR) + fields (owner, expiry, mitigations)
    - how exceptions relate to “Local Addenda (project-owned)”
  - Add `Docs/Playbooks/Target-Repo-Quality-Gates.md` describing how teams should implement **repeatable, tool-enforced** quality gates in target repos (beyond governance rules and agent behavior), including:
    - wiring repo-owned enforcement into `make lint` / `make ci`
    - examples of architecture enforcement mechanisms (e.g., SwiftLint rules for forbidden imports/layering, module boundary checks, dependency rules)
    - guidance on reliability: fast, deterministic, CI-friendly checks; avoid flaky heuristics; prefer compile-time/static analysis where possible
    - how to pair governance verification (`gov-ci` / `gov-verify`) with product checks without coupling the governance include to stack-specific tools

- **Decision record**
  - Add a MADR under `Docs/Decisions/` capturing the governance content tagging/versioning policy:
    - governance content tags: `gov/vX.Y.Z` (SemVer, covers core + all profiles)
    - tool tags remain: `agent-gov/vX.Y.Z`
    - high-level “breaking change” criteria for governance content

- **Alternative approach doc integration**
  - Update `multi-repo-governance-strategy.md` with a short “mapping to this repo’s approach” section:
    - what we already do (managed blocks + `verify`)
    - what we adopt (SemVer/pinning framing, golden commands, exception discipline)
  - Add/expand a “gaps/open questions” section that becomes a living backlog for the governance strategy.

- **Cross-linking + consistency pass**
  - Ensure terminology is consistent across README + playbooks + the alternative doc (managed blocks, local addenda, pinning, enforcement).
  - Ensure docs point readers to the binding governing docs (the big three) and the relevant playbooks.

### Out of scope

- Changing `.governance/config.yaml` schema or `agent-gov` CLI behavior.
- Implementing new enforcement in the `agent-gov` CLI/toolchain beyond the planned repo-level golden commands (no new/changed CLI behavior).
- Adding new governance profiles (unless required to clarify docs with a minimal example).

## Approach

- Lead with an operator-friendly narrative: “how to adopt” + “how to keep governance consistent over time”.
- Make implicit assumptions explicit:
  - what is enforced vs guidance
  - what is pinned and why
  - what is allowed to vary per repo (local addenda, config knobs) vs what is controlled centrally
- Keep playbooks short, procedural, and copy/paste friendly.

## Checkpoints

- [x] Checkpoint 1 — Add/shape the new README opening section (what/why/goals/non-goals) and a short mental model
- [x] Checkpoint 2 — Clarify README pinning/versioning story and introduce “golden commands” guidance for target repos
- [x] Checkpoint 3 — Add playbook: `Docs/Playbooks/Governance-Upgrades.md`
- [x] Checkpoint 4 — Add playbook: `Docs/Playbooks/Governance-Exceptions.md`
- [ ] Checkpoint 5 — Add playbook: `Docs/Playbooks/Target-Repo-Quality-Gates.md` (tool-enforced architecture/quality checks; SwiftLint example)
- [ ] Checkpoint 6 — Add decision record (MADR) for governance content SemVer tagging policy (`gov/vX.Y.Z`)
- [ ] Checkpoint 7 — Extract/adopt best ideas from scratch notes (`multi-repo-governance-strategy.md`) into tracked docs; then delete the scratch file (must remain untracked)
- [ ] Checkpoint 8 — Implement golden commands for this repo (add markdown fmt/lint; wire into `make ci`)
  - Add `gov-ci` to `Governance/Templates/Make/agent-gov.mk` and update README to recommend `ci: gov-ci` in consuming repos
- [ ] Checkpoint 9 — Consistency + cross-links pass (README ↔ playbooks ↔ governing docs)
- [ ] Checkpoint 10 — Run quality gates and ensure docs render cleanly (links, headings, no contradictory guidance)
- [ ] Final checkpoint — PR wrap-up (final approval gate)

## Completion checklist

- [ ] Set frontmatter `status: completed`
- [ ] Check off completed checkpoint(s) above and add PR/commit references as needed
- [ ] Check off all quality gates below (or document any exceptions)

## Quality gates / test plan

- [ ] `make ci`
- [ ] `make fmt` updates Go + markdown formatting as expected (no unexpected file churn)
- [ ] `make lint` lints markdown (and Go, if applicable) with clear failure output
- [ ] Manual read-through: README “happy path” is clear and complete for a new reader
- [ ] Links check (spot-check all new internal links and key external references)

## Notes / open questions

- Decision: use **SemVer tags for governance content** (core + all profiles) as `gov/vX.Y.Z`, separate from tool tags `agent-gov/vX.Y.Z`.
- Decision: “breaking change” criteria for governance content (requires **MAJOR** bump) — treat as breaking if a repo pinned to the prior version likely needs **human remediation** beyond bumping `source.ref` and re-running `gov-sync` to return to green/compliant. Breaking includes:
  - New **required repo changes** (new required files/paths/layout beyond managed-block updates)
  - Stricter **required quality gates** likely to fail existing repos by default (new mandatory checks/threshold increases)
  - Managed-block contract changes (prefix/heading/verify semantics) that invalidate existing managed blocks or local addenda preservation expectations
  - Template contract changes that force migration of existing docs to remain compliant
  - Profile meaning changes where an existing conforming repo becomes non-conforming without repo-side changes
  - Policy: since governance content uses one version for core + all profiles, a change that is breaking for **any** supported profile is a breaking governance release
- Decision: exceptions live in the **decision record (MADR)**. Do not introduce a separate exceptions registry.
- Decision: standardize markdown golden commands on **Prettier** (format) + **markdownlint-cli2** (lint), aligned with `.markdownlint.jsonc` and common Cursor/editor setups.
- Decision: pin/install deterministically via **npm + committed `package-lock.json`**:
  - add a minimal `package.json` with devDependencies for `prettier` and `markdownlint-cli2`
  - commit `package-lock.json`
  - use `npm ci` in CI (and run via npm scripts and/or `npx`)
- Decision: pin Node runtime by **major** line (Option B). Add `.nvmrc` (e.g. `22`) and document “Node 22.x + npm required” for the markdown tooling path.
- Decision: `gov-ci` in `Governance/Templates/Make/agent-gov.mk` will be **CI-safe** and run **`gov-verify` only** (no sync / no working tree changes). Optionally print “next steps” guidance on failure (e.g., run `make gov-sync`, then re-run `make gov-ci`).
- Decision: keep integration **purely explicit** in consuming repos (recommended pattern: `ci: gov-ci`). Do not add any “attach to ci” behavior in the shared include.
- Decision: capture the governance content SemVer policy as a durable decision record in `Docs/Decisions/` (MADR).
