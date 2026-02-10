# Shared Makefile include for consuming `agent-gov` in target repos.
#
# Recommended usage (in your repo's top-level Makefile):
#
#   # Prefer a single source of truth for governance operations:
#   -include tools/make/agent-gov.mk
#
#   # Pin the tool tag (SemVer tag in this repo, e.g. agent-gov/v0.4.0)
#   AGENT_GOV_TAG ?= agent-gov/v0.4.0
#
#   # Where to place the downloaded binary
#   AGENT_GOV_BIN ?= tools/bin/agent-gov
#
#   # Optional (recommended) explicit config path; omit to rely on auto-discovery
#   GOV_CONFIG ?= .governance/config.yaml
#
#   # Choose download source: github (default) or gitlab
#   AGENT_GOV_SOURCE ?= github
#
#   # GitHub Releases settings (required when AGENT_GOV_SOURCE=github)
#   AGENT_GOV_GITHUB_ORG ?= <org>
#   AGENT_GOV_GITHUB_REPO ?= agent-governance-strategy
#
#   # GitLab Generic Package Registry settings (required when AGENT_GOV_SOURCE=gitlab)
#   GITLAB_HOST ?= gitlab.com
#   GITLAB_PROJECT_ID ?= 12345678
#   AGENT_GOV_PKG ?= agent-gov
#   GITLAB_PKG_USERNAME ?= token
#   GITLAB_PKG_TOKEN ?=
#
# Notes:
# - Asset naming must match the release pipeline: agent-gov_<os>_<arch>
# - `AGENT_GOV_TAG` is used as-is for GitHub, and reduced to `vX.Y.Z` for GitLab packages.

ifndef AGENT_GOV_MK_INCLUDED
AGENT_GOV_MK_INCLUDED := 1

AGENT_GOV_TAG ?= agent-gov/v0.4.0
AGENT_GOV_BIN ?= tools/bin/agent-gov
GOV_CONFIG ?= .governance/config.yaml

AGENT_GOV_SOURCE ?= github

AGENT_GOV_GITHUB_ORG ?=
AGENT_GOV_GITHUB_REPO ?= agent-governance-strategy

GITLAB_HOST ?= gitlab.com
GITLAB_PROJECT_ID ?=
AGENT_GOV_PKG ?= agent-gov
GITLAB_PKG_USERNAME ?= token
GITLAB_PKG_TOKEN ?=

# Generic Package Registry versions are a single path segment, so we use just `vX.Y.Z`.
AGENT_GOV_VERSION ?= $(notdir $(AGENT_GOV_TAG))

.PHONY: agent-gov gov-init gov-sync gov-verify gov-build

agent-gov:
	@mkdir -p $$(dirname "$(AGENT_GOV_BIN)")
	@if [ -x "$(AGENT_GOV_BIN)" ]; then exit 0; fi; \
	  os="$$(uname -s | tr '[:upper:]' '[:lower:]')"; \
	  arch="$$(uname -m)"; \
	  if [ "$$arch" = "x86_64" ]; then arch="amd64"; fi; \
	  if [ "$$arch" = "aarch64" ]; then arch="arm64"; fi; \
	  asset="agent-gov_$${os}_$${arch}"; \
	  case "$(AGENT_GOV_SOURCE)" in \
	    github) \
	      if [ -z "$(AGENT_GOV_GITHUB_ORG)" ]; then echo "AGENT_GOV_GITHUB_ORG is required for AGENT_GOV_SOURCE=github"; exit 1; fi; \
	      url="https://github.com/$(AGENT_GOV_GITHUB_ORG)/$(AGENT_GOV_GITHUB_REPO)/releases/download/$(AGENT_GOV_TAG)/$${asset}"; \
	      curl -fsSL "$${url}" -o "$(AGENT_GOV_BIN)"; \
	      ;; \
	    gitlab) \
	      if [ -z "$(GITLAB_PROJECT_ID)" ]; then echo "GITLAB_PROJECT_ID is required for AGENT_GOV_SOURCE=gitlab"; exit 1; fi; \
	      if [ -z "$(GITLAB_PKG_TOKEN)" ]; then echo "GITLAB_PKG_TOKEN is required for AGENT_GOV_SOURCE=gitlab"; exit 1; fi; \
	      url="https://$(GITLAB_HOST)/api/v4/projects/$(GITLAB_PROJECT_ID)/packages/generic/$(AGENT_GOV_PKG)/$(AGENT_GOV_VERSION)/$${asset}"; \
	      curl -fsSL --user "$(GITLAB_PKG_USERNAME):$(GITLAB_PKG_TOKEN)" "$${url}" -o "$(AGENT_GOV_BIN)"; \
	      ;; \
	    *) \
	      echo "AGENT_GOV_SOURCE must be 'github' or 'gitlab' (got: $(AGENT_GOV_SOURCE))"; \
	      exit 1; \
	      ;; \
	  esac; \
	  chmod +x "$(AGENT_GOV_BIN)"

gov-init: agent-gov
	@if [ -n "$(GOV_CONFIG)" ]; then \
	  $(AGENT_GOV_BIN) init --config "$(GOV_CONFIG)"; \
	else \
	  $(AGENT_GOV_BIN) init; \
	fi

gov-sync: agent-gov
	@if [ -n "$(GOV_CONFIG)" ]; then \
	  $(AGENT_GOV_BIN) sync --config "$(GOV_CONFIG)"; \
	else \
	  $(AGENT_GOV_BIN) sync; \
	fi

gov-verify: agent-gov
	@if [ -n "$(GOV_CONFIG)" ]; then \
	  $(AGENT_GOV_BIN) verify --config "$(GOV_CONFIG)"; \
	else \
	  $(AGENT_GOV_BIN) verify; \
	fi

gov-build: agent-gov
	@if [ -n "$(GOV_CONFIG)" ]; then \
	  $(AGENT_GOV_BIN) build --config "$(GOV_CONFIG)"; \
	else \
	  $(AGENT_GOV_BIN) build; \
	fi

endif
