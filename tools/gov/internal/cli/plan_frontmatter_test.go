package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadPlanFrontmatter_NoFrontmatter(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "p.md")
	if err := os.WriteFile(p, []byte("# hi\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	_, _, ok, err := readPlanFrontmatter(p)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false")
	}
}

func TestReadPlanFrontmatter_MissingClosingFence(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "p.md")
	if err := os.WriteFile(p, []byte("---\nbranch: feat/x\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	_, _, ok, err := readPlanFrontmatter(p)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false")
	}
}

func TestReadPlanFrontmatter_InvalidYAML_IsIgnored(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "p.md")
	if err := os.WriteFile(p, []byte("---\n: :\n---\n# x\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	_, _, ok, err := readPlanFrontmatter(p)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false for invalid yaml")
	}
}

func TestWalkPlanFiles_MissingPlansDir(t *testing.T) {
	tmp := t.TempDir()
	paths, err := walkPlanFiles(filepath.Join(tmp, "Docs", "Plans"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(paths) != 0 {
		t.Fatalf("expected empty, got %v", paths)
	}
}

func TestWalkPlanFiles_SkipsTemplatesAndNonMarkdown(t *testing.T) {
	tmp := t.TempDir()
	plansDir := filepath.Join(tmp, "Docs", "Plans")
	if err := os.MkdirAll(filepath.Join(plansDir, "feat"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(plansDir, "Plan.Template.md"), []byte("x\n"), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	if err := os.WriteFile(filepath.Join(plansDir, "feat", "notes.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatalf("write txt: %v", err)
	}
	plan := filepath.Join(plansDir, "feat", "real-plan.md")
	if err := os.WriteFile(plan, []byte("# plan\n"), 0o644); err != nil {
		t.Fatalf("write plan: %v", err)
	}

	paths, err := walkPlanFiles(plansDir)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(paths) != 1 || filepath.Clean(paths[0]) != filepath.Clean(plan) {
		t.Fatalf("expected only %q, got %v", plan, paths)
	}
}
