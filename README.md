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

- `agent-gov/v0.1.0` (SemVer)

In a target repo, store the downloaded binary at a repo-local path (and gitignore it), for example:

- `tools/bin/agent-gov`

The target repo then provides a small wrapper (Makefile/script) that:

- downloads the pinned version for the current OS/arch if missing
- runs `tools/bin/agent-gov ...`

#### Copy/paste bootstrap (GitLab Releases)

This one-liner downloads the pinned tool into `tools/bin/agent-gov`:

```bash
AGENT_GOV_TAG="agent-gov/v1.1.0" AGENT_GOV_GITLAB_REPO="bsmith.quanata/agent-governance-strategy" bash -c 'set -euo pipefail; bin="tools/bin/agent-gov"; dir="$(dirname "$bin")"; mkdir -p "${dir}"; os="$(uname -s | tr "[:upper:]" "[:lower:]")"; arch="$(uname -m)"; [ "$arch" = "x86_64" ] && arch="amd64"; [ "$arch" = "aarch64" ] && arch="arm64"; asset="agent-gov_${os}_${arch}"; echo "downloading ${asset} from ${AGENT_GOV_GITLAB_REPO}@${AGENT_GOV_TAG}"; glab release download "${AGENT_GOV_TAG}" -R "${AGENT_GOV_GITLAB_REPO}" --asset-name "${asset}" -D "${dir}"; mv -f "${dir}/${asset}" "${bin}"; chmod +x "${bin}"; "${bin}" --version'
```

Notes:

- Requires GitLab CLI (`glab`) and authentication (recommended: `glab auth login`). This works cleanly for private repos.
After downloading:

- Add `tools/bin/agent-gov` to `.gitignore`
- Create `.governance/config.yaml` (see below), then run `tools/bin/agent-gov init --config .governance/config.yaml`

### 2) Add `.governance/config.yaml` to the target repo (team-safe default)

In the target repo, you can generate `.governance/config.yaml` using the CLI:

```bash
tools/bin/agent-gov bootstrap \
  --config .governance/config.yaml \
  --source-repo "git@github.com:<org>/agent-governance-strategy.git" \
  --source-ref "gov/v2026.02.09" \
  --profile "docs-only" \
  --non-interactive
```

Helpful discovery commands:

- List profiles at a given repo/ref:

```bash
tools/bin/agent-gov bootstrap \
  --source-repo "git@github.com:<org>/agent-governance-strategy.git" \
  --source-ref "gov/v2026.02.09" \
  --profile "docs-only" \
  --list-profiles
```

- One-shot: write config and immediately run `init`:

```bash
tools/bin/agent-gov bootstrap \
  --config .governance/config.yaml \
  --source-repo "git@github.com:<org>/agent-governance-strategy.git" \
  --source-ref "gov/v2026.02.09" \
  --profile "docs-only" \
  --non-interactive \
  --run-init
```

Or create `.governance/config.yaml` by hand:

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

### Makefile integration for target repos (shared include)

Instead of copy/pasting Make targets into every repo, we provide a shared Makefile include:

- Governance source: `Governance/Templates/Make/agent-gov.mk`
- Emitted into target repos by profiles as: `tools/make/agent-gov.mk`

In a target repo, your top-level `Makefile` can be as small as:

```make
# Optional include (present after you run `agent-gov init/sync`, or if you vendor it yourself).
-include tools/make/agent-gov.mk

# Pin the tool tag (SemVer tag in this repo, e.g. agent-gov/v0.1.0)
AGENT_GOV_TAG ?= agent-gov/v0.1.0

# Where to place the downloaded binary
AGENT_GOV_BIN ?= tools/bin/agent-gov

# Optional (recommended) explicit config path; omit to rely on auto-discovery
GOV_CONFIG ?= .governance/config.yaml

# Choose download source: github (default) or gitlab
AGENT_GOV_SOURCE ?= github

# GitHub Releases settings (required when AGENT_GOV_SOURCE=github)
AGENT_GOV_GITHUB_ORG ?= <org>
# AGENT_GOV_GITHUB_REPO ?= agent-governance-strategy

# For GitLab Generic Package Registry (when AGENT_GOV_SOURCE=gitlab), set:
# GITLAB_HOST ?= gitlab.com
# GITLAB_PROJECT_ID ?= 12345678
# AGENT_GOV_PKG ?= agent-gov
# GITLAB_PKG_USERNAME ?= token
# GITLAB_PKG_TOKEN ?= <deploy-token>
```

Once included, you get standard targets:

- `make gov-init`
- `make gov-sync`
- `make gov-verify`
- `make gov-build`

GitLab note (private projects): downloads must be authenticated. A team-safe approach is a **deploy token** with `read_package_registry` stored as masked CI/CD variables in the consuming repo. This repo’s GitLab release flow publishes binaries to the **Generic Package Registry** under package name:

- `agent-gov` for tags `agent-gov/v*`
- `agent-gov-test` for tags `agent-gov/test/v*`

If you need to bootstrap the include file initially, you can vendor/copy it from this repo’s `Governance/Templates/Make/agent-gov.mk` (pinned to the same ref as your `.governance/config.yaml`).

## Releasing `agent-gov` (maintainers)

Pushing a tag matching `agent-gov/v*` triggers CI to build and publish release assets for a small OS/arch set.

For **safe in-repo testing**, you can push a tag matching `agent-gov/test/v*`. Test-tag releases are created as **draft prereleases**.

### 1) Cut a tag

From a clean commit on `main`:

```bash
git tag agent-gov/vX.Y.Z
git push origin agent-gov/vX.Y.Z
```

Safe test example:

```bash
git tag agent-gov/test/v0.0.0-test1
git push origin agent-gov/test/v0.0.0-test1
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

