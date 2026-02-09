.DEFAULT_GOAL := help

.PHONY: help ci fmt test coverage gov-smoke preflight gov-preflight gov-preflight-gocli

GOV_MIN_COVERAGE ?= 85
# Many Go environments set `GOFLAGS=-mod=vendor` globally to enforce vendoring.
# This repo's `tools/gov` module does not commit a vendor tree, so we force module
# mode for Go commands invoked by this Makefile to keep `make ci` reliable.
TOOLS_GOV_GOFLAGS ?= -mod=mod

help: ## Show available make targets
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_.-]+:.*##/ {printf "  %-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

preflight: gov-preflight ## Sanity check branch + baseline

gov-preflight: ## Run agent-gov preflight (root scope)
	@cd tools/gov && GOFLAGS="$(TOOLS_GOV_GOFLAGS)" go run ./cmd/agent-gov preflight \
	  --require "Makefile" \
	  --require "Governance/Profiles/backend-go-hex/profile.yaml" \
	  --require "Governance/Profiles/mobile-clean/profile.yaml" \
	  --require "tools/gov/go.mod"

gov-preflight-gocli: ## Run agent-gov preflight (embedded tools/gov scope)
	@cd tools/gov && GOFLAGS="$(TOOLS_GOV_GOFLAGS)" go run ./cmd/agent-gov preflight \
	  --require "go.mod" \
	  --require "cmd/agent-gov/main.go"

ci: fmt test coverage gov-smoke ## Run all CI checks

fmt: ## Format Go sources
	@echo "Formatting Go sources"
	@cd tools/gov && gofmt -w $$(git ls-files '*.go')

test: ## Run Go tests
	@echo "Running Go tests"
	@cd tools/gov && GOFLAGS="$(TOOLS_GOV_GOFLAGS)" go test ./...

coverage: ## Enforce minimum CLI test coverage
	@echo "Checking tools/gov coverage >= $(GOV_MIN_COVERAGE)%"
	@cd tools/gov && rm -f coverage.out coverage.txt
	@cd tools/gov && GOFLAGS="$(TOOLS_GOV_GOFLAGS)" go test ./... -coverprofile=coverage.out >/dev/null
	@cd tools/gov && GOFLAGS="$(TOOLS_GOV_GOFLAGS)" go tool cover -func=coverage.out | tee coverage.txt >/dev/null
	@cd tools/gov && awk '/^total:/ {gsub(/%/,"",$$3); pct=$$3} END { if (pct+0 < $(GOV_MIN_COVERAGE)) { printf "coverage %.1f%% is below %d%%\\n", pct, $(GOV_MIN_COVERAGE); exit 1 } else { printf "coverage %.1f%%\\n", pct; } }' coverage.txt

gov-smoke: ## Smoke test agent-gov init/verify
	@echo "Smoke test agent-gov init/verify"
	@set -eu; \
	tmp="$$(mktemp -d)"; \
	repo_root="$$(pwd)"; \
	( cd tools/gov && GOFLAGS="$(TOOLS_GOV_GOFLAGS)" go build -o "$$tmp/agent-gov" ./cmd/agent-gov ); \
	mkdir -p "$$tmp/target/.governance"; \
	printf '%s\n' \
	  'schemaVersion: 1' \
	  'source:' \
	  "  repo: \"$$repo_root\"" \
	  '  ref: "HEAD"' \
	  '  profile: "mobile-clean-ios"' \
	  'paths:' \
	  '  docsRoot: "."' \
	  'sync:' \
	  '  managedBlockPrefix: "GOV"' \
	  '  localAddendaHeading: "Local Addenda (project-owned)"' \
	  > "$$tmp/target/.governance/config.yaml"; \
	cd "$$tmp/target"; \
	"$$tmp/agent-gov" init --config .governance/config.yaml >/dev/null; \
	"$$tmp/agent-gov" verify --config .governance/config.yaml >/dev/null; \
	echo "ok"

