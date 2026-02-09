# agent-governance-strategy

This repository is a **governance source + toolchain**. It is not a product/runtime system.

- **Governance sources** live in `Governance/` (core fragments, profiles, templates, playbooks).
- The CLI tool **`agent-gov`** lives in `tools/gov/` and is used to **init / sync / verify** governance docs in a target repo using **managed blocks**.

## What `agent-gov` does

`agent-gov` generates and maintains governance documents in a target repo using:

- **Managed blocks**: deterministic, tool-owned sections that can be updated in-place.
- **Local addenda**: a project-owned section (by default headed **“Local Addenda (project-owned)”**) that is preserved and never overwritten.

Commands:

- `init`: create governance docs with managed blocks + local addenda
- `sync`: update managed blocks in-place
- `verify`: check that managed blocks match expected content (CI-friendly)
- `build`: assemble a governance bundle into an output folder (for inspection/artifacts)

## Recommended usage (apply governance to another repo)

### 1) Use a pinned `agent-gov` binary (recommended for teams)

To avoid “every repo has a different tool snapshot”, we recommend shipping `agent-gov` as a **versioned binary** and having each target repo **pin** the version it uses.

Recommended tag format for the tool:

- `agent-gov/v0.4.0` (SemVer)

In a target repo, store the downloaded binary at a repo-local path (and gitignore it), for example:

- `tools/bin/agent-gov`

The target repo then provides a small wrapper (Makefile/script) that:

- downloads the pinned version for the current OS/arch if missing
- runs `tools/bin/agent-gov ...`

### 2) Add `.governance/config.yaml` to the target repo (team-safe default)

In the target repo, create `.governance/config.yaml`:

```yaml
schemaVersion: 1

source:
  # Team-safe: use a shared remote URL so everyone can run the same config.
  repo: "git@github.com:<org>/agent-governance-strategy.git" # or https://...

  # Strongly recommended: pin to an immutable tag or commit SHA for repeatability.
  # Avoid moving refs like HEAD unless you explicitly want “latest on every run”.
  ref: "gov/v2026.02.09"

  # Choose a profile ID from this repo under `Governance/Profiles/`.
  profile: "docs-only"

paths:
  # Where to emit docs inside the target repo. Defaults to "." when omitted.
  docsRoot: "."

sync:
  managedBlockPrefix: "GOV"
  localAddendaHeading: "Local Addenda (project-owned)"
```

Why remote URLs matter for teams:

- `source.repo` is used as a **git clone source**. If you commit a machine-local path (like `/Users/alice/...`), it will not work on other machines.

### 3) Run `init`, then `verify`

From anywhere inside the target repo:

- Initialize governance docs:

```bash
tools/bin/agent-gov init --config .governance/config.yaml
```

- Verify (good for CI):

```bash
tools/bin/agent-gov verify --config .governance/config.yaml
```

- Sync later (after updating `source.ref`):

```bash
tools/bin/agent-gov sync --config .governance/config.yaml
```

Notes:

- If you omit `--config`, `agent-gov` **auto-discovers** the nearest `.governance/config.yaml` by walking upward from the current working directory.
- You can always be explicit with `--config .governance/config.yaml`.

### Example Makefile snippet for target repos (pinned binary)

Below is a minimal pattern target repos can adopt. It downloads a pinned `agent-gov` binary into `tools/bin/agent-gov` and then uses it.

You will need to replace `<org>` with your GitHub org/user (and ensure releases publish assets with the expected names).

```make
AGENT_GOV_VERSION ?= agent-gov/v0.4.0
AGENT_GOV_BIN ?= tools/bin/agent-gov

.PHONY: agent-gov gov-init gov-sync gov-verify

agent-gov:
	@mkdir -p $$(dirname "$(AGENT_GOV_BIN)")
	@if [ ! -x "$(AGENT_GOV_BIN)" ]; then \
	  os="$$(uname -s | tr '[:upper:]' '[:lower:]')"; \
	  arch="$$(uname -m)"; \
	  if [ "$$arch" = "x86_64" ]; then arch="amd64"; fi; \
	  if [ "$$arch" = "aarch64" ]; then arch="arm64"; fi; \
	  asset="agent-gov_$${os}_$${arch}"; \
	  url="https://github.com/<org>/agent-governance-strategy/releases/download/$(AGENT_GOV_VERSION)/$${asset}"; \
	  echo "downloading $${url}"; \
	  curl -fsSL "$${url}" -o "$(AGENT_GOV_BIN)"; \
	  chmod +x "$(AGENT_GOV_BIN)"; \
	fi

gov-init: agent-gov
	@$(AGENT_GOV_BIN) init --config .governance/config.yaml

gov-sync: agent-gov
	@$(AGENT_GOV_BIN) sync --config .governance/config.yaml

gov-verify: agent-gov
	@$(AGENT_GOV_BIN) verify --config .governance/config.yaml
```

## Releasing `agent-gov` (maintainers)

Pushing a tag matching `agent-gov/v*` triggers CI to build and publish release assets for a small OS/arch set.

### 1) Cut a tag

From a clean commit on `main`:

```bash
git tag agent-gov/vX.Y.Z
git push origin agent-gov/vX.Y.Z
```

### 2) Confirm assets

The GitHub Release for the tag should contain assets named:

- `agent-gov_darwin_amd64`
- `agent-gov_darwin_arm64`
- `agent-gov_linux_amd64`

## Local development workflow (optional override)

When authoring governance changes locally, you may want to run a target repo against a local checkout of this governance source.

Recommended pattern:

- Keep the committed `.governance/config.yaml` pointing at the **remote URL** (team-safe).
- Add a **gitignored** dev config, e.g. `.governance/config.dev.yaml`, that points to your local checkout:

```yaml
schemaVersion: 1
source:
  repo: "/path/to/your/local/agent-governance-strategy"
  ref: "HEAD"
  profile: "docs-only"
paths:
  docsRoot: "."
sync:
  managedBlockPrefix: "GOV"
  localAddendaHeading: "Local Addenda (project-owned)"
```

Then run:

```bash
tools/bin/agent-gov sync --config .governance/config.dev.yaml
```

## Profiles

Profiles are defined under `Governance/Profiles/<id>/profile.yaml`. v1 includes:

- `docs-only`
- `backend-go-hex`
- `mobile-clean-ios`

## Contributing to this repo

This repo uses `make` targets to run checks:

- `make preflight`
- `make ci` (format, tests, coverage gate, and a CLI smoke test)

For repo working agreements and quality gates, see:

- `Non-Negotiables.md`
- `Architecture.md`
- `Constitution.md`

