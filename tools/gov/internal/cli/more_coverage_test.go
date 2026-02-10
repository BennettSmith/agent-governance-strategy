package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_Build_RequiresOutFlag(t *testing.T) {
	var outBuf, errBuf bytes.Buffer
	code := Run([]string{"agent-gov", "build", "--config", "does-not-matter.yaml"}, &outBuf, &errBuf)
	if code != 2 {
		t.Fatalf("expected 2, got %d stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(errBuf.String(), "--out is required") {
		t.Fatalf("expected --out error, got:\n%s", errBuf.String())
	}
}

func TestRepoRootForConfig_WhenNotInGovernanceDir_UsesConfigDir(t *testing.T) {
	tmp := t.TempDir()
	cfg := filepath.Join(tmp, "somewhere", "config.yaml")
	got := repoRootForConfig(cfg)
	wantDir := filepath.Dir(cfg)
	wantAbs, _ := filepath.Abs(wantDir)
	if filepath.Clean(got) != filepath.Clean(wantAbs) {
		t.Fatalf("expected %q, got %q", wantAbs, got)
	}
}

func TestRepoRootForConfig_WhenInGovernanceDir_UsesParentDir(t *testing.T) {
	tmp := t.TempDir()
	cfg := filepath.Join(tmp, ".governance", "config.yaml")
	got := repoRootForConfig(cfg)
	wantAbs, _ := filepath.Abs(tmp)
	if filepath.Clean(got) != filepath.Clean(wantAbs) {
		t.Fatalf("expected %q, got %q", wantAbs, got)
	}
}

func TestRepoRootForConfig_WhenAbsFails_FallsBackToRelative(t *testing.T) {
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	wd := t.TempDir()
	if err := os.Chdir(wd); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	// Remove the working directory so filepath.Abs (which calls Getwd for relative paths)
	// is likely to fail.
	_ = os.RemoveAll(wd)
	defer func() { _ = os.Chdir(oldCwd) }()

	got := repoRootForConfig(".governance/config.yaml")
	if got == "" {
		t.Fatalf("expected non-empty fallback")
	}
}

func TestResolveRepoPathIfLocal_ResolvesRelativeToConfigDir_ForCoverage(t *testing.T) {
	tmp := t.TempDir()
	cfg := filepath.Join(tmp, ".governance", "config.yaml")
	repoDir := filepath.Join(tmp, "govsrc")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	got := resolveRepoPathIfLocal(cfg, "../govsrc")
	if !filepath.IsAbs(got) {
		t.Fatalf("expected abs path, got %q", got)
	}
	if filepath.Clean(got) != filepath.Clean(repoDir) {
		t.Fatalf("expected %q, got %q", repoDir, got)
	}
}

func TestIsTTY_FalseForTempFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "notty")
	if err != nil {
		t.Fatalf("temp: %v", err)
	}
	defer func() { _ = f.Close() }()
	if isTTY(f) {
		t.Fatalf("expected false")
	}
}

func TestListProfilesFromCheckout_SkipsNonDirs(t *testing.T) {
	tmp := t.TempDir()
	checkout := filepath.Join(tmp, "co")
	if err := os.MkdirAll(filepath.Join(checkout, "Governance", "Profiles"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// Non-dir entry.
	if err := os.WriteFile(filepath.Join(checkout, "Governance", "Profiles", "README.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	// Valid profile dir.
	if err := os.MkdirAll(filepath.Join(checkout, "Governance", "Profiles", "p"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(checkout, "Governance", "Profiles", "p", "profile.yaml"), []byte("schemaVersion: 1\nid: p\ndocuments: []\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	infos, err := listProfilesFromCheckout(checkout)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(infos) != 1 || infos[0].ID != "p" {
		t.Fatalf("unexpected infos: %+v", infos)
	}
}
