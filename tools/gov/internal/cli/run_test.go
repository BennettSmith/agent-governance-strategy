package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_NoArgs_ShowsUsage(t *testing.T) {
	var out, errOut bytes.Buffer
	code := Run([]string{"agent-gov"}, &out, &errOut)
	if code != 2 {
		t.Fatalf("expected 2, got %d", code)
	}
	if !strings.Contains(errOut.String(), "agent-gov <command>") {
		t.Fatalf("expected usage in stderr, got:\n%s", errOut.String())
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	var out, errOut bytes.Buffer
	code := Run([]string{"agent-gov", "nope"}, &out, &errOut)
	if code != 2 {
		t.Fatalf("expected 2, got %d", code)
	}
	if !strings.Contains(errOut.String(), "unknown command") {
		t.Fatalf("expected unknown command error, got:\n%s", errOut.String())
	}
}

func TestRun_Help_PrintsUsage(t *testing.T) {
	var out, errOut bytes.Buffer
	code := Run([]string{"agent-gov", "help"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("expected 0, got %d", code)
	}
	if !strings.Contains(out.String(), "Commands:") {
		t.Fatalf("expected usage in stdout, got:\n%s", out.String())
	}
}

func TestRun_DashH_PrintsUsage(t *testing.T) {
	var out, errOut bytes.Buffer
	code := Run([]string{"agent-gov", "-h"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("expected 0, got %d", code)
	}
	if !strings.Contains(out.String(), "agent-gov <command>") {
		t.Fatalf("expected usage, got:\n%s", out.String())
	}
}

func TestRun_VersionCommand_PrintsVersion(t *testing.T) {
	var out, errOut bytes.Buffer
	code := Run([]string{"agent-gov", "version"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("expected 0, got %d stderr=%s", code, errOut.String())
	}
	if strings.TrimSpace(out.String()) == "" {
		t.Fatalf("expected version output in stdout, got empty")
	}
}

func TestRun_DashDashVersion_PrintsVersion(t *testing.T) {
	var out, errOut bytes.Buffer
	code := Run([]string{"agent-gov", "--version"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("expected 0, got %d stderr=%s", code, errOut.String())
	}
	if strings.TrimSpace(out.String()) == "" {
		t.Fatalf("expected version output in stdout, got empty")
	}
}

func TestRun_Build_EndToEnd_FromLocalTaggedGovernanceRepo(t *testing.T) {
	// This is an end-to-end CLI test that exercises:
	// config parsing, local repo path resolution, source fetch, profile loading, and build output.
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	outDir := filepath.Join(tmp, "out")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "NonNegotiables.Core.md"), "CORE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "Constitution.Core.md"), "CONST\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "NonNegotiables.Profile.md"), "PROFILE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "Constitution.Profile.md"), "C-PROFILE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "Architecture.Profile.md"), "ARCH\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: backend-go-hex
documents:
  - output: Non-Negotiables.md
    fragments:
      - ../../Core/NonNegotiables.Core.md
      - ./NonNegotiables.Profile.md
  - output: Constitution.md
    fragments:
      - ../../Core/Constitution.Core.md
      - ./Constitution.Profile.md
  - output: Architecture.md
    fragments:
      - ./Architecture.Profile.md
`)+"\n")

	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "gov")
	mustRun(t, srcRepo, "git", "tag", "v0.0.1")

	targetRepo := filepath.Join(tmp, "target")
	if err := os.MkdirAll(filepath.Join(targetRepo, ".governance"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cfgPath := filepath.Join(targetRepo, ".governance", "config.yaml")
	// config lives in .governance, so use ".." to refer to the target repo root (unused here)
	// and point source.repo at the governance source repo.
	writeFile(t, cfgPath, strings.TrimSpace(`
schemaVersion: 1
source:
  repo: `+srcRepo+`
  ref: "v0.0.1"
  profile: "backend-go-hex"
paths:
  docsRoot: "."
sync:
  managedBlockPrefix: "GOV"
  localAddendaHeading: "Local Addenda (project-owned)"
`)+"\n")

	// Run build from within the module, pointing at the config, writing to outDir.
	var outBuf, errBuf bytes.Buffer
	code := Run([]string{"agent-gov", "build", "--config", cfgPath, "--out", outDir}, &outBuf, &errBuf)
	if code != 0 {
		t.Fatalf("expected 0, got %d stderr=%s", code, errBuf.String())
	}
	if _, err := os.Stat(filepath.Join(outDir, "Non-Negotiables.md")); err != nil {
		t.Fatalf("expected output doc: %v", err)
	}
}

func TestRun_InitSyncVerify_EndToEnd_FromLocalTaggedGovernanceRepo(t *testing.T) {
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	cache := filepath.Join(tmp, "cache")
	target := filepath.Join(tmp, "target")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	// Minimal profile with 1 doc.
	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "NonNegotiables.Core.md"), "COREv1\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "NonNegotiables.Profile.md"), "PROFILEv1\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: backend-go-hex
documents:
  - output: Non-Negotiables.md
    fragments:
      - ../../Core/NonNegotiables.Core.md
      - ./NonNegotiables.Profile.md
`)+"\n")

	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "v1")
	mustRun(t, srcRepo, "git", "tag", "v0.0.1")

	if err := os.MkdirAll(filepath.Join(target, ".governance"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cfgPath := filepath.Join(target, ".governance", "config.yaml")
	writeFile(t, cfgPath, strings.TrimSpace(`
schemaVersion: 1
source:
  repo: `+srcRepo+`
  ref: "v0.0.1"
  profile: "backend-go-hex"
paths:
  docsRoot: "."
  cacheDir: `+cache+`
sync:
  managedBlockPrefix: "GOV"
  localAddendaHeading: "Local Addenda (project-owned)"
`)+"\n")

	// Run init from target working directory (RepoRoot=".").
	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	if err := os.Chdir(target); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	var outBuf, errBuf bytes.Buffer
	if code := Run([]string{"agent-gov", "init", "--config", cfgPath}, &outBuf, &errBuf); code != 0 {
		t.Fatalf("init code=%d stderr=%s", code, errBuf.String())
	}
	if _, err := os.Stat(filepath.Join(target, "Non-Negotiables.md")); err != nil {
		t.Fatalf("expected initialized doc: %v", err)
	}
	outBuf.Reset()
	errBuf.Reset()
	if code := Run([]string{"agent-gov", "verify", "--config", cfgPath}, &outBuf, &errBuf); code != 0 {
		t.Fatalf("verify code=%d stderr=%s", code, errBuf.String())
	}

	// Running sync without changing the source should still succeed.
	outBuf.Reset()
	errBuf.Reset()
	if code := Run([]string{"agent-gov", "sync", "--config", cfgPath}, &outBuf, &errBuf); code != 0 {
		t.Fatalf("sync code=%d stderr=%s", code, errBuf.String())
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func mustRun(t *testing.T, dir string, exe string, args ...string) {
	t.Helper()
	cmd := execCommandContext(context.Background(), exe, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v failed: %v\n%s", exe, args, err, string(out))
	}
}
