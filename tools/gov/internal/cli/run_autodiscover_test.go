package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_Init_AutoDiscoversConfig_LogsSelectedPath(t *testing.T) {
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	target := filepath.Join(tmp, "target")
	sub := filepath.Join(target, "subdir")
	cache := filepath.Join(tmp, "cache")

	// Governance source repo.
	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "NonNegotiables.Core.md"), "COREv1\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "p", "NonNegotiables.Profile.md"), "PROFILEv1\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "p", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: p
documents:
  - output: Non-Negotiables.md
    fragments:
      - ../../Core/NonNegotiables.Core.md
      - ./NonNegotiables.Profile.md
`)+"\n")
	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "v1")
	mustRun(t, srcRepo, "git", "tag", "v0.0.1")

	// Target repo with config at root, but run from a subdir so auto-discovery logs.
	if err := os.MkdirAll(filepath.Join(target, ".governance"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cfgPath := filepath.Join(target, ".governance", "config.yaml")
	writeFile(t, cfgPath, strings.TrimSpace(`
schemaVersion: 1
source:
  repo: `+srcRepo+`
  ref: "v0.0.1"
  profile: "p"
paths:
  docsRoot: "."
  cacheDir: `+cache+`
sync:
  managedBlockPrefix: "GOV"
  localAddendaHeading: "Local Addenda (project-owned)"
`)+"\n")

	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	if err := os.Chdir(sub); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	var outBuf, errBuf bytes.Buffer
	code := Run([]string{"agent-gov", "init"}, &outBuf, &errBuf)
	if code != 0 {
		t.Fatalf("init code=%d stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(errBuf.String(), "using config:") {
		t.Fatalf("expected config discovery log, got:\n%s", errBuf.String())
	}
	if _, err := os.Stat(filepath.Join(target, "Non-Negotiables.md")); err != nil {
		t.Fatalf("expected initialized doc: %v", err)
	}
}

// Anchor unused import in older Go toolchains (keeps file consistent with other tests).
var _ = context.Background
