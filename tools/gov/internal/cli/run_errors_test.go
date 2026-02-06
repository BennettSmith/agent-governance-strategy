package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_BuildRequiresOut(t *testing.T) {
	var out, errOut bytes.Buffer
	code := Run([]string{"agent-gov", "build"}, &out, &errOut)
	if code != 2 {
		t.Fatalf("expected 2, got %d", code)
	}
	if !strings.Contains(errOut.String(), "--out is required") {
		t.Fatalf("expected missing out error, got:\n%s", errOut.String())
	}
}

func TestRun_BuildConfigError(t *testing.T) {
	var out, errOut bytes.Buffer
	code := Run([]string{"agent-gov", "build", "--config", "/does/not/exist.yaml", "--out", "/tmp/out"}, &out, &errOut)
	if code != 2 {
		t.Fatalf("expected 2, got %d", code)
	}
	if !strings.Contains(errOut.String(), "config error") {
		t.Fatalf("expected config error, got:\n%s", errOut.String())
	}
}

func TestRun_VerifyReportsIssuesWhenDocsMissing(t *testing.T) {
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	target := filepath.Join(tmp, "target")
	cache := filepath.Join(tmp, "cache")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "NonNegotiables.Core.md"), "CORE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "NonNegotiables.Profile.md"), "PROFILE\n")
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

	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	if err := os.Chdir(target); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	var out, errOut bytes.Buffer
	code := Run([]string{"agent-gov", "verify", "--config", cfgPath}, &out, &errOut)
	if code != 1 {
		t.Fatalf("expected 1, got %d", code)
	}
	if !strings.Contains(errOut.String(), "verification failed") {
		t.Fatalf("expected verification failure output, got:\n%s", errOut.String())
	}
	if !strings.Contains(errOut.String(), "Non-Negotiables.md") {
		t.Fatalf("expected doc listed as issue, got:\n%s", errOut.String())
	}
}

func TestResolveRepoPathIfLocal_ResolvesRelativeToConfigDir(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, ".governance", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// Create a fake repo directory next to config dir.
	repoDir := filepath.Join(tmp, "repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}

	got := resolveRepoPathIfLocal(cfgPath, "../repo")
	if got == "../repo" {
		t.Fatalf("expected resolved path, got %q", got)
	}
	if !strings.Contains(got, "repo") {
		t.Fatalf("expected repo in path, got %q", got)
	}
}

func TestResolveRepoPathIfLocal_ReturnsAbsoluteRepoIfExists(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, ".governance", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	repoDir := filepath.Join(tmp, "repoabs")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}

	got := resolveRepoPathIfLocal(cfgPath, repoDir)
	if got != repoDir {
		t.Fatalf("expected same absolute path, got %q", got)
	}
}

func TestResolveRepoPathIfLocal_FallsBackToCwdRelative(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, ".governance", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cwd := filepath.Join(tmp, "cwd")
	if err := os.MkdirAll(cwd, 0o755); err != nil {
		t.Fatalf("mkdir cwd: %v", err)
	}
	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	if err := os.Chdir(cwd); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	if err := os.MkdirAll("reporel", 0o755); err != nil {
		t.Fatalf("mkdir reporel: %v", err)
	}

	got := resolveRepoPathIfLocal(cfgPath, "reporel")
	if !strings.Contains(got, "reporel") {
		t.Fatalf("expected reporel in path, got %q", got)
	}
}

func TestRun_InvalidFlag_Returns2(t *testing.T) {
	var out, errOut bytes.Buffer
	code := Run([]string{"agent-gov", "build", "--nope"}, &out, &errOut)
	if code != 2 {
		t.Fatalf("expected 2, got %d", code)
	}
}

func TestRun_VerifyFailsWhenSourceRepoInvalid(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "target")
	cache := filepath.Join(tmp, "cache")
	if err := os.MkdirAll(filepath.Join(target, ".governance"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cfgPath := filepath.Join(target, ".governance", "config.yaml")
	writeFile(t, cfgPath, strings.TrimSpace(`
schemaVersion: 1
source:
  repo: /does/not/exist
  ref: "v0.0.1"
  profile: "mobile-clean-ios"
paths:
  docsRoot: "."
  cacheDir: `+cache+`
sync:
  managedBlockPrefix: "GOV"
  localAddendaHeading: "Local Addenda (project-owned)"
`)+"\n")

	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	if err := os.Chdir(target); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	var out, errOut bytes.Buffer
	code := Run([]string{"agent-gov", "verify", "--config", cfgPath}, &out, &errOut)
	if code != 1 {
		t.Fatalf("expected 1, got %d", code)
	}
	if !strings.Contains(errOut.String(), "verify failed") {
		t.Fatalf("expected verify failed message, got:\n%s", errOut.String())
	}
}

func TestRun_SubcommandInternalErrorPath(t *testing.T) {
	// Directly exercise the default branch in runSubcommand.
	var out, errOut bytes.Buffer
	code := runSubcommand("wat", []string{}, &out, &errOut)
	if code != 1 {
		t.Fatalf("expected 1, got %d", code)
	}
	if !strings.Contains(errOut.String(), "internal error") {
		t.Fatalf("expected internal error, got:\n%s", errOut.String())
	}
}

func TestRun_BuildReturns1OnBuildFailure(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "target")
	cache := filepath.Join(tmp, "cache")
	if err := os.MkdirAll(filepath.Join(target, ".governance"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cfgPath := filepath.Join(target, ".governance", "config.yaml")
	// Point at an invalid governance repo to force builder.Build failure.
	writeFile(t, cfgPath, strings.TrimSpace(`
schemaVersion: 1
source:
  repo: /does/not/exist
  ref: "v0.0.1"
  profile: "backend-go-hex"
paths:
  docsRoot: "."
  cacheDir: `+cache+`
sync:
  managedBlockPrefix: "GOV"
  localAddendaHeading: "Local Addenda (project-owned)"
`)+"\n")

	var out, errOut bytes.Buffer
	code := Run([]string{"agent-gov", "build", "--config", cfgPath, "--out", filepath.Join(tmp, "out")}, &out, &errOut)
	if code != 1 {
		t.Fatalf("expected 1, got %d", code)
	}
	if !strings.Contains(errOut.String(), "build failed") {
		t.Fatalf("expected build failed output, got:\n%s", errOut.String())
	}
}

func TestRun_SyncReturns1WhenDocsMissing(t *testing.T) {
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	target := filepath.Join(tmp, "target")
	cache := filepath.Join(tmp, "cache")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "NonNegotiables.Core.md"), "CORE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "NonNegotiables.Profile.md"), "PROFILE\n")
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

	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	if err := os.Chdir(target); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	var out, errOut bytes.Buffer
	code := Run([]string{"agent-gov", "sync", "--config", cfgPath}, &out, &errOut)
	if code != 1 {
		t.Fatalf("expected 1, got %d", code)
	}
	if !strings.Contains(errOut.String(), "sync failed") {
		t.Fatalf("expected sync failed, got:\n%s", errOut.String())
	}
}
