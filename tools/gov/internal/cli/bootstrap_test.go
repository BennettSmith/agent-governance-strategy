package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_Bootstrap_NonInteractive_WritesConfig(t *testing.T) {
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	target := filepath.Join(tmp, "target")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	// Minimal profile with a description.
	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "NonNegotiables.Core.md"), "COREv1\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "docs-only", "Architecture.Profile.md"), "ARCH\n")
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
	cfgPath := filepath.Join(target, ".governance", "config.yaml")

	var outBuf, errBuf bytes.Buffer
	code := Run([]string{
		"agent-gov", "bootstrap",
		"--config", cfgPath,
		"--source-repo", srcRepo,
		"--source-ref", "gov/v0.0.1",
		"--profile", "docs-only",
		"--non-interactive",
	}, &outBuf, &errBuf)
	if code != 0 {
		t.Fatalf("bootstrap code=%d stderr=%s", code, errBuf.String())
	}
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("expected config written: %v", err)
	}
}

func TestRun_Bootstrap_ListProfiles_PrintsIDs(t *testing.T) {
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	cache := filepath.Join(tmp, "cache")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "a", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: a
description: A profile
documents: []
`)+"\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "b", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: b
description: B profile
documents: []
`)+"\n")

	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "gov")
	mustRun(t, srcRepo, "git", "tag", "gov/v0.0.1")

	var outBuf, errBuf bytes.Buffer
	code := Run([]string{
		"agent-gov", "bootstrap",
		"--source-repo", srcRepo,
		"--source-ref", "gov/v0.0.1",
		"--profile", "a",
		"--cache-dir", cache,
		"--non-interactive",
		"--list-profiles",
	}, &outBuf, &errBuf)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, errBuf.String())
	}
	got := outBuf.String()
	if !strings.Contains(got, "a") || !strings.Contains(got, "b") {
		t.Fatalf("expected profile ids in output, got:\n%s", got)
	}
}

func TestRun_Bootstrap_PrintOnly_DoesNotWriteFile(t *testing.T) {
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	cache := filepath.Join(tmp, "cache")
	cfgPath := filepath.Join(tmp, ".governance", "config.yaml")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "docs-only", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: docs-only
documents: []
`)+"\n")
	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "gov")
	mustRun(t, srcRepo, "git", "tag", "gov/v0.0.1")

	var outBuf, errBuf bytes.Buffer
	code := Run([]string{
		"agent-gov", "bootstrap",
		"--config", cfgPath,
		"--source-repo", srcRepo,
		"--source-ref", "gov/v0.0.1",
		"--profile", "docs-only",
		"--cache-dir", cache,
		"--non-interactive",
		"--print",
	}, &outBuf, &errBuf)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, errBuf.String())
	}
	if _, err := os.Stat(cfgPath); err == nil {
		t.Fatalf("did not expect config file to be written")
	}
	if !strings.Contains(outBuf.String(), "schemaVersion: 1") {
		t.Fatalf("expected yaml printed to stdout, got:\n%s", outBuf.String())
	}
}

func TestRun_Bootstrap_ForceOverwrite(t *testing.T) {
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	cache := filepath.Join(tmp, "cache")
	target := filepath.Join(tmp, "target")
	cfgPath := filepath.Join(target, ".governance", "config.yaml")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "docs-only", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: docs-only
documents: []
`)+"\n")
	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "gov")
	mustRun(t, srcRepo, "git", "tag", "gov/v0.0.1")

	writeFile(t, cfgPath, "old: true\n")

	// Without --force should fail.
	{
		var outBuf, errBuf bytes.Buffer
		code := Run([]string{
			"agent-gov", "bootstrap",
			"--config", cfgPath,
			"--source-repo", srcRepo,
			"--source-ref", "gov/v0.0.1",
			"--profile", "docs-only",
			"--cache-dir", cache,
			"--non-interactive",
		}, &outBuf, &errBuf)
		if code == 0 {
			t.Fatalf("expected failure without --force")
		}
	}

	// With --force should overwrite.
	{
		var outBuf, errBuf bytes.Buffer
		code := Run([]string{
			"agent-gov", "bootstrap",
			"--config", cfgPath,
			"--source-repo", srcRepo,
			"--source-ref", "gov/v0.0.1",
			"--profile", "docs-only",
			"--cache-dir", cache,
			"--non-interactive",
			"--force",
		}, &outBuf, &errBuf)
		if code != 0 {
			t.Fatalf("code=%d stderr=%s", code, errBuf.String())
		}
		b, err := os.ReadFile(cfgPath)
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		if !strings.Contains(string(b), "schemaVersion: 1") {
			t.Fatalf("expected new config content, got:\n%s", string(b))
		}
	}
}

func TestRun_Bootstrap_ProfileNotFound(t *testing.T) {
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	cache := filepath.Join(tmp, "cache")
	target := filepath.Join(tmp, "target")
	cfgPath := filepath.Join(target, ".governance", "config.yaml")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "docs-only", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: docs-only
documents: []
`)+"\n")
	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "gov")
	mustRun(t, srcRepo, "git", "tag", "gov/v0.0.1")

	var outBuf, errBuf bytes.Buffer
	code := Run([]string{
		"agent-gov", "bootstrap",
		"--config", cfgPath,
		"--source-repo", srcRepo,
		"--source-ref", "gov/v0.0.1",
		"--profile", "nope",
		"--cache-dir", cache,
		"--non-interactive",
	}, &outBuf, &errBuf)
	if code == 0 {
		t.Fatalf("expected failure for missing profile")
	}
}

func TestRun_Bootstrap_RelativeConfigPath_ResolvesAndWrites(t *testing.T) {
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "docs-only", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: docs-only
documents: []
`)+"\n")
	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "gov")
	mustRun(t, srcRepo, "git", "tag", "gov/v0.0.1")

	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	var outBuf, errBuf bytes.Buffer
	code := Run([]string{
		"agent-gov", "bootstrap",
		"--config", ".governance/config.yaml",
		"--source-repo", srcRepo,
		"--source-ref", "gov/v0.0.1",
		"--profile", "docs-only",
		"--non-interactive",
	}, &outBuf, &errBuf)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, errBuf.String())
	}
	if _, err := os.Stat(filepath.Join(tmp, ".governance", "config.yaml")); err != nil {
		t.Fatalf("expected config written: %v", err)
	}
}

func TestRun_Bootstrap_ListProfiles_PrintsIDsWithoutDescription(t *testing.T) {
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	cache := filepath.Join(tmp, "cache")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "p", "profile.yaml"), "schemaVersion: 1\nid: p\ndocuments: []\n")
	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "gov")
	mustRun(t, srcRepo, "git", "tag", "gov/v0.0.1")

	var outBuf, errBuf bytes.Buffer
	code := Run([]string{
		"agent-gov", "bootstrap",
		"--source-repo", srcRepo,
		"--source-ref", "gov/v0.0.1",
		"--profile", "p",
		"--cache-dir", cache,
		"--non-interactive",
		"--list-profiles",
	}, &outBuf, &errBuf)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, errBuf.String())
	}
	if strings.Contains(outBuf.String(), "\t") {
		t.Fatalf("expected no tab when description empty, got:\n%s", outBuf.String())
	}
}

func TestRun_Bootstrap_NonInteractive_MissingRequiredFlags(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, ".governance", "config.yaml")

	var outBuf, errBuf bytes.Buffer
	code := Run([]string{
		"agent-gov", "bootstrap",
		"--config", cfgPath,
		"--non-interactive",
	}, &outBuf, &errBuf)
	if code != 2 {
		t.Fatalf("expected 2, got %d stderr=%s", code, errBuf.String())
	}
}

func TestPromptBootstrapValues_ReadsMissingFields(t *testing.T) {
	in := bytes.NewBufferString("repo-path\nref-tag\n")
	var outBuf, errBuf bytes.Buffer

	repo := ""
	ref := ""
	profile := "docs-only"
	dRoot := ""

	if err := promptBootstrapValues(in, &outBuf, &errBuf, t.TempDir(), &repo, &ref, &profile, &dRoot); err != nil {
		t.Fatalf("prompt err: %v", err)
	}
	if repo != "repo-path" {
		t.Fatalf("expected repo set, got %q", repo)
	}
	if ref != "ref-tag" {
		t.Fatalf("expected ref set, got %q", ref)
	}
	if dRoot != "." {
		t.Fatalf("expected docsRoot default '.', got %q", dRoot)
	}
}

func TestPromptBootstrapValues_SelectsProfileFromList(t *testing.T) {
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	cache := filepath.Join(tmp, "cache")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "b", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: b
description: B profile
documents: []
`)+"\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "a", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: a
description: A profile
documents: []
`)+"\n")

	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "gov")
	mustRun(t, srcRepo, "git", "tag", "gov/v0.0.1")

	// Provide: sourceRepo, sourceRef, then choose profile "1" from the sorted list (a, b).
	in := bytes.NewBufferString(srcRepo + "\n" + "gov/v0.0.1\n" + "1\n")
	var outBuf, errBuf bytes.Buffer

	repo := ""
	ref := ""
	profile := ""
	dRoot := "."
	if err := promptBootstrapValues(in, &outBuf, &errBuf, cache, &repo, &ref, &profile, &dRoot); err != nil {
		t.Fatalf("prompt err: %v", err)
	}
	if profile != "a" {
		t.Fatalf("expected profile 'a' selected, got %q", profile)
	}
}

func TestDefaultCacheDir_EndsWithGovbuilder(t *testing.T) {
	cd, err := defaultCacheDir()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !strings.HasSuffix(filepath.ToSlash(cd), "/govbuilder") {
		t.Fatalf("expected suffix /govbuilder, got %q", cd)
	}
}

func TestWriteConfigFile_RespectsForce(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, ".governance", "config.yaml")
	if err := writeConfigFile(p, []byte("a: 1\n"), false); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := writeConfigFile(p, []byte("b: 2\n"), false); err == nil {
		t.Fatalf("expected error when file exists without force")
	}
	if err := writeConfigFile(p, []byte("b: 2\n"), true); err != nil {
		t.Fatalf("force write: %v", err)
	}
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(b) != "b: 2\n" {
		t.Fatalf("expected overwritten content, got %q", string(b))
	}
}

func TestRun_Bootstrap_RunInit_WritesDocs(t *testing.T) {
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	target := filepath.Join(tmp, "target")
	cache := filepath.Join(tmp, "cache")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "NonNegotiables.Profile.md"), "PROFILEv1\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "NonNegotiables.Core.md"), "COREv1\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: backend-go-hex
description: Go hex profile
documents:
  - output: Non-Negotiables.md
    fragments:
      - ../../Core/NonNegotiables.Core.md
      - ./NonNegotiables.Profile.md
`)+"\n")

	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "v1")
	mustRun(t, srcRepo, "git", "tag", "gov/v0.0.1")

	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cfgPath := filepath.Join(target, ".governance", "config.yaml")

	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	if err := os.Chdir(target); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	var outBuf, errBuf bytes.Buffer
	code := Run([]string{
		"agent-gov", "bootstrap",
		"--config", cfgPath,
		"--source-repo", srcRepo,
		"--source-ref", "gov/v0.0.1",
		"--profile", "backend-go-hex",
		"--cache-dir", cache,
		"--non-interactive",
		"--run-init",
	}, &outBuf, &errBuf)
	if code != 0 {
		t.Fatalf("bootstrap code=%d stderr=%s", code, errBuf.String())
	}
	if _, err := os.Stat(filepath.Join(target, "Non-Negotiables.md")); err != nil {
		t.Fatalf("expected initialized doc: %v", err)
	}
}

// Reuse mustRun/writeFile helpers from run_test.go (same package).
var _ = context.Background
