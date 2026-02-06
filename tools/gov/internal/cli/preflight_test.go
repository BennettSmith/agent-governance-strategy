package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPreflight_FailsOnMain(t *testing.T) {
	tmp := t.TempDir()
	repo := filepath.Join(tmp, "repo")
	mustRun(t, tmp, "git", "init", repo)
	mustRun(t, repo, "git", "config", "user.email", "test@example.com")
	mustRun(t, repo, "git", "config", "user.name", "Test")
	writeFile(t, filepath.Join(repo, ".governance", "config.yaml"), "schemaVersion: 1\nsource:\n  repo: .\n  ref: \"HEAD\"\n  profile: \"backend-go-hex\"\npaths:\n  docsRoot: \".\"\n")
	writeFile(t, filepath.Join(repo, "README.md"), "x\n")
	mustRun(t, repo, "git", "add", ".")
	mustRun(t, repo, "git", "commit", "-m", "init")

	// Ensure we are on main (default for modern git).
	var out, errOut bytes.Buffer
	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	_ = os.Chdir(repo)
	code := Run([]string{"agent-gov", "preflight"}, &out, &errOut)
	if code == 0 {
		t.Fatalf("expected failure on main")
	}
	if !strings.Contains(errOut.String(), "on main") {
		t.Fatalf("expected main error, got:\n%s", errOut.String())
	}
}

func TestPreflight_FailsWhenOnBranchFromCompletedPlan(t *testing.T) {
	tmp := t.TempDir()
	repo := filepath.Join(tmp, "repo")
	mustRun(t, tmp, "git", "init", repo)
	mustRun(t, repo, "git", "config", "user.email", "test@example.com")
	mustRun(t, repo, "git", "config", "user.name", "Test")
	writeFile(t, filepath.Join(repo, ".governance", "config.yaml"), "schemaVersion: 1\nsource:\n  repo: .\n  ref: \"HEAD\"\n  profile: \"backend-go-hex\"\npaths:\n  docsRoot: \".\"\n")
	writeFile(t, filepath.Join(repo, "Docs", "Plans", "feat", "old-plan.md"), "# plan\n")
	writeFile(t, filepath.Join(repo, "README.md"), "x\n")
	mustRun(t, repo, "git", "add", ".")
	mustRun(t, repo, "git", "commit", "-m", "init")

	mustRun(t, repo, "git", "checkout", "-b", "feat/old-plan")

	var out, errOut bytes.Buffer
	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	_ = os.Chdir(repo)
	code := Run([]string{"agent-gov", "preflight"}, &out, &errOut)
	if code != 1 {
		t.Fatalf("expected 1, got %d stderr=%s", code, errOut.String())
	}
	if !strings.Contains(errOut.String(), "existing plan") {
		t.Fatalf("expected plan-branch failure, got:\n%s", errOut.String())
	}
}

func TestPreflight_AllowsActivePlanBranch(t *testing.T) {
	tmp := t.TempDir()
	repo := filepath.Join(tmp, "repo")
	mustRun(t, tmp, "git", "init", repo)
	mustRun(t, repo, "git", "config", "user.email", "test@example.com")
	mustRun(t, repo, "git", "config", "user.name", "Test")
	writeFile(t, filepath.Join(repo, ".governance", "config.yaml"), "schemaVersion: 1\nsource:\n  repo: .\n  ref: \"HEAD\"\n  profile: \"backend-go-hex\"\npaths:\n  docsRoot: \".\"\n")
	writeFile(t, filepath.Join(repo, "Docs", "Plans", "feat", "active-plan.md"), strings.TrimSpace(`
---
branch: feat/active-plan
status: active
---

# plan
`)+"\n")
	writeFile(t, filepath.Join(repo, "README.md"), "x\n")
	mustRun(t, repo, "git", "add", ".")
	mustRun(t, repo, "git", "commit", "-m", "init")

	mustRun(t, repo, "git", "checkout", "-b", "feat/active-plan")

	var out, errOut bytes.Buffer
	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	_ = os.Chdir(repo)
	code := Run([]string{"agent-gov", "preflight"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("expected 0, got %d stderr=%s", code, errOut.String())
	}
	if strings.TrimSpace(out.String()) != "ok" {
		t.Fatalf("expected ok, got %q", out.String())
	}
}

func TestPreflight_FailsWhenRequiredPathMissing(t *testing.T) {
	tmp := t.TempDir()
	repo := filepath.Join(tmp, "repo")
	mustRun(t, tmp, "git", "init", repo)
	mustRun(t, repo, "git", "config", "user.email", "test@example.com")
	mustRun(t, repo, "git", "config", "user.name", "Test")
	writeFile(t, filepath.Join(repo, ".governance", "config.yaml"), "schemaVersion: 1\nsource:\n  repo: .\n  ref: \"HEAD\"\n  profile: \"backend-go-hex\"\npaths:\n  docsRoot: \".\"\n")
	writeFile(t, filepath.Join(repo, "README.md"), "x\n")
	mustRun(t, repo, "git", "add", ".")
	mustRun(t, repo, "git", "commit", "-m", "init")
	mustRun(t, repo, "git", "checkout", "-b", "feat/new-plan")

	var out, errOut bytes.Buffer
	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	_ = os.Chdir(repo)
	code := Run([]string{"agent-gov", "preflight", "--require", "nope.txt"}, &out, &errOut)
	if code != 1 {
		t.Fatalf("expected 1, got %d stderr=%s", code, errOut.String())
	}
	if !strings.Contains(errOut.String(), "required path missing") {
		t.Fatalf("expected required-path failure, got:\n%s", errOut.String())
	}
}

func TestBranchForPlanPath_UsesFrontmatterWhenPresent(t *testing.T) {
	tmp := t.TempDir()
	repo := filepath.Join(tmp, "repo")
	plan := filepath.Join(repo, "Docs", "Plans", "feat", "x.md")
	writeFile(t, plan, strings.TrimSpace(`
---
branch: feat/x
status: active
---

# plan
`)+"\n")

	got, err := branchForPlanPath(repo, "Docs/Plans/feat/x.md")
	if err != nil {
		t.Fatalf("branchForPlanPath: %v", err)
	}
	if got != "feat/x" {
		t.Fatalf("expected feat/x got %q", got)
	}
}

func TestPreflight_FailsWhenNoConfigFound(t *testing.T) {
	tmp := t.TempDir()
	var out, errOut bytes.Buffer
	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	_ = os.Chdir(tmp)
	code := Run([]string{"agent-gov", "preflight"}, &out, &errOut)
	if code != 1 {
		t.Fatalf("expected 1, got %d stderr=%s", code, errOut.String())
	}
	if !strings.Contains(errOut.String(), "could not find") {
		t.Fatalf("expected not-found failure, got:\n%s", errOut.String())
	}
}

func TestPreflight_FailsWhenDetachedHead(t *testing.T) {
	tmp := t.TempDir()
	repo := filepath.Join(tmp, "repo")
	mustRun(t, tmp, "git", "init", repo)
	mustRun(t, repo, "git", "config", "user.email", "test@example.com")
	mustRun(t, repo, "git", "config", "user.name", "Test")
	writeFile(t, filepath.Join(repo, ".governance", "config.yaml"), "schemaVersion: 1\nsource:\n  repo: .\n  ref: \"HEAD\"\n  profile: \"backend-go-hex\"\npaths:\n  docsRoot: \".\"\n")
	writeFile(t, filepath.Join(repo, "README.md"), "x\n")
	mustRun(t, repo, "git", "add", ".")
	mustRun(t, repo, "git", "commit", "-m", "init")

	sha := strings.TrimSpace(string(mustRunOut(t, repo, "git", "rev-parse", "HEAD")))
	mustRun(t, repo, "git", "checkout", sha)

	var out, errOut bytes.Buffer
	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	_ = os.Chdir(repo)
	code := Run([]string{"agent-gov", "preflight"}, &out, &errOut)
	if code != 1 {
		t.Fatalf("expected 1, got %d stderr=%s", code, errOut.String())
	}
	if !strings.Contains(errOut.String(), "detached HEAD") {
		t.Fatalf("expected detached failure, got:\n%s", errOut.String())
	}
}

func TestPreflight_ActivePlanFlag_AllowsBranchEvenWithoutStatusActive(t *testing.T) {
	tmp := t.TempDir()
	repo := filepath.Join(tmp, "repo")
	mustRun(t, tmp, "git", "init", repo)
	mustRun(t, repo, "git", "config", "user.email", "test@example.com")
	mustRun(t, repo, "git", "config", "user.name", "Test")
	writeFile(t, filepath.Join(repo, ".governance", "config.yaml"), "schemaVersion: 1\nsource:\n  repo: .\n  ref: \"HEAD\"\n  profile: \"backend-go-hex\"\npaths:\n  docsRoot: \".\"\n")
	planPath := filepath.Join(repo, "Docs", "Plans", "feat", "active.md")
	writeFile(t, planPath, "# plan\n")
	writeFile(t, filepath.Join(repo, "README.md"), "x\n")
	mustRun(t, repo, "git", "add", ".")
	mustRun(t, repo, "git", "commit", "-m", "init")
	mustRun(t, repo, "git", "checkout", "-b", "feat/active")

	var out, errOut bytes.Buffer
	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	_ = os.Chdir(repo)
	code := Run([]string{"agent-gov", "preflight", "--active-plan", "Docs/Plans/feat/active.md"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("expected 0, got %d stderr=%s", code, errOut.String())
	}
	if strings.TrimSpace(out.String()) != "ok" {
		t.Fatalf("expected ok, got %q", out.String())
	}
}

func TestPreflight_ErrorsOnMultipleActivePlans(t *testing.T) {
	tmp := t.TempDir()
	repo := filepath.Join(tmp, "repo")
	mustRun(t, tmp, "git", "init", repo)
	mustRun(t, repo, "git", "config", "user.email", "test@example.com")
	mustRun(t, repo, "git", "config", "user.name", "Test")
	writeFile(t, filepath.Join(repo, ".governance", "config.yaml"), "schemaVersion: 1\nsource:\n  repo: .\n  ref: \"HEAD\"\n  profile: \"backend-go-hex\"\npaths:\n  docsRoot: \".\"\n")
	writeFile(t, filepath.Join(repo, "Docs", "Plans", "feat", "a.md"), strings.TrimSpace(`
---
branch: feat/a
status: active
---

# plan
`)+"\n")
	writeFile(t, filepath.Join(repo, "Docs", "Plans", "feat", "b.md"), strings.TrimSpace(`
---
branch: feat/b
status: active
---

# plan
`)+"\n")
	writeFile(t, filepath.Join(repo, "README.md"), "x\n")
	mustRun(t, repo, "git", "add", ".")
	mustRun(t, repo, "git", "commit", "-m", "init")
	mustRun(t, repo, "git", "checkout", "-b", "feat/work")

	var out, errOut bytes.Buffer
	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	_ = os.Chdir(repo)
	code := Run([]string{"agent-gov", "preflight"}, &out, &errOut)
	if code != 2 {
		t.Fatalf("expected 2, got %d stderr=%s", code, errOut.String())
	}
	if !strings.Contains(errOut.String(), "multiple active plans") {
		t.Fatalf("expected multiple-active error, got:\n%s", errOut.String())
	}
}

func mustRunOut(t *testing.T, dir string, exe string, args ...string) []byte {
	t.Helper()
	cmd := execCommandContext(context.Background(), exe, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v failed: %v\n%s", exe, args, err, string(out))
	}
	return out
}
