package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBootstrap_EnvDefaults_EndToEnd_ThenInitAndVerify(t *testing.T) {
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	target := filepath.Join(tmp, "target")
	cache := filepath.Join(tmp, "cache")

	// Governance source repo with a minimal profile.
	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "NonNegotiables.Core.md"), "COREv1\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "docs-only", "Architecture.Profile.md"), "ARCHv1\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "docs-only", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: docs-only
description: Docs-only profile
documents:
  - output: Architecture.md
    fragments:
      - ./Architecture.Profile.md
`)+"\n")
	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "gov")
	mustRun(t, srcRepo, "git", "tag", "gov/v0.0.1")

	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Run bootstrap from target dir using env defaults.
	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	if err := os.Chdir(target); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	t.Setenv("AGENT_GOV_SOURCE_REPO", srcRepo)
	t.Setenv("AGENT_GOV_SOURCE_REF", "gov/v0.0.1")
	t.Setenv("AGENT_GOV_PROFILE", "docs-only")
	t.Setenv("AGENT_GOV_DOCS_ROOT", ".")
	// Cache dir can be supplied via env for fetch behavior, but should not be written
	// into the generated config unless explicitly requested via --cache-dir.
	t.Setenv("AGENT_GOV_CACHE_DIR", cache)

	var outBuf, errBuf bytes.Buffer
	if code := Run([]string{"agent-gov", "bootstrap", "--config", ".governance/config.yaml", "--non-interactive"}, &outBuf, &errBuf); code != 0 {
		t.Fatalf("bootstrap code=%d stderr=%s", code, errBuf.String())
	}

	cfgPath := filepath.Join(target, ".governance", "config.yaml")
	cfgBytes, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if strings.Contains(string(cfgBytes), "cacheDir:") {
		t.Fatalf("did not expect cacheDir written by env default, got:\n%s", string(cfgBytes))
	}

	outBuf.Reset()
	errBuf.Reset()
	if code := Run([]string{"agent-gov", "init", "--config", cfgPath}, &outBuf, &errBuf); code != 0 {
		t.Fatalf("init code=%d stderr=%s", code, errBuf.String())
	}
	if _, err := os.Stat(filepath.Join(target, "Architecture.md")); err != nil {
		t.Fatalf("expected initialized doc: %v", err)
	}

	outBuf.Reset()
	errBuf.Reset()
	if code := Run([]string{"agent-gov", "verify", "--config", cfgPath}, &outBuf, &errBuf); code != 0 {
		t.Fatalf("verify code=%d stderr=%s", code, errBuf.String())
	}
	if strings.TrimSpace(outBuf.String()) != "ok" {
		t.Fatalf("expected ok, got %q", strings.TrimSpace(outBuf.String()))
	}
}

var _ = context.Background

