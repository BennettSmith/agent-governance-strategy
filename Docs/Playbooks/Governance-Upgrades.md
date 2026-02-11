# Governance upgrades

This playbook describes how to safely upgrade a target repository’s governance to a newer version of this repo’s published governance content and/or a newer `agent-gov` tool version.

## Preconditions

- Your target repo has governance configured via `.governance/config.yaml`.
- You can run `agent-gov` in the target repo (directly, or via your Makefile wrapper).
- Your normal repo quality gates exist and are runnable (typically `make ci`).

## Two independent upgrades

There are two things you can upgrade independently:

- **Governance content**: the `source.ref` in `.governance/config.yaml` (tagged as `gov/vX.Y.Z`).
- **Tool**: the `agent-gov` binary version you download/use (tagged as `agent-gov/vX.Y.Z`).

In general: upgrade governance content first (to keep policy/docs aligned), then upgrade the tool if needed.

## Upgrade governance content (recommended starting point)

1. Choose the new governance content tag:
   - Example: `gov/v1.2.3`
2. Update `.governance/config.yaml`:
   - Set `source.ref` to the new tag (or an immutable commit SHA).
3. Sync managed blocks:
   - `agent-gov sync --config .governance/config.yaml`
4. Verify without writing:
   - `agent-gov verify --config .governance/config.yaml`
5. Run your repo quality gates:
   - `make ci`

Notes:

- `sync` updates **managed blocks** only. It should not overwrite your **Local Addenda (project-owned)** sections.
- `verify` is CI-friendly: it should not modify the working tree.

## Upgrade the tool (`agent-gov`)

1. Choose the new tool tag:
   - Example: `agent-gov/v1.2.3`
2. Update your repo’s tool pin (Makefile/script/env var) to that tag.
3. Download the new binary (or let your wrapper do it).
4. Sanity check:
   - `agent-gov --version`
5. Re-run:
   - `agent-gov verify --config .governance/config.yaml`
   - `make ci`

## CI wiring recommendation

In target repos, keep governance verification explicit and CI-safe:

- Prefer a CI target that runs **verify only** (no sync), e.g.:
  - `make gov-ci` → `agent-gov verify ...`
- Then wire into your repo’s `ci` target explicitly:
  - `ci: gov-ci`

## Rollback strategy

If an upgrade makes the repo fail verification or quality gates:

1. Revert the pin(s):
   - `source.ref` back to the previous `gov/vX.Y.Z` (and/or revert the `agent-gov/vX.Y.Z` tool tag).
2. Re-run sync:
   - `agent-gov sync --config .governance/config.yaml`
3. Re-verify and re-run quality gates:
   - `agent-gov verify --config .governance/config.yaml`
   - `make ci`

## Breaking changes (governance content)

Governance content uses SemVer tags (`gov/vX.Y.Z`). Treat a governance content change as **breaking** (MAJOR bump) if a repo pinned to the previous version is likely to require **human remediation** beyond:

- bumping `source.ref`, and
- re-running `agent-gov sync`
