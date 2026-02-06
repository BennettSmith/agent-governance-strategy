package builder

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuild_EmitsManagedDocAndTemplates_FromLocalTaggedSource(t *testing.T) {
	ctx := context.Background()
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	outDir := filepath.Join(tmp, "out")
	cache := filepath.Join(tmp, "cache")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	// Minimal governance tree.
	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "NonNegotiables.Core.md"), "CORE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "NonNegotiables.Profile.md"), "PROFILE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "Constitution.Profile.md"), "C-PROFILE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "Architecture.Profile.md"), "A-PROFILE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Templates", "Decisions", "MADR.Template.md"), "MADR\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Templates", "Plans", "Plan.Template.md"), "PLAN\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "Playbooks", "Go-Packaging.md"), "PLAYBOOK\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "Playbooks", "Hexagonal-Ports-And-Adapters.md"), "HEX\n")

	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "Constitution.Core.md"), "CORE-CONST\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: backend-go-hex
description: test
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
templates:
  - source: ../../Templates/Decisions/MADR.Template.md
    output: Docs/Decisions/MADR.Template.md
  - source: ../../Templates/Plans/Plan.Template.md
    output: Docs/Plans/Plan.Template.md
playbooks:
  - source: ./Playbooks/Go-Packaging.md
    output: Docs/Playbooks/Go-Packaging.md
`)+"\n")

	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "gov")
	mustRun(t, srcRepo, "git", "tag", "v0.0.1")

	res, err := Build(ctx, BuildOptions{
		OutDir:         outDir,
		DocsRoot:       ".",
		CacheDir:       cache,
		SourceRepo:     srcRepo,
		SourceRef:      "v0.0.1",
		ProfileID:      "backend-go-hex",
		MarkerPrefix:   "GOV",
		AddendaHeading: "Local Addenda (project-owned)",
	})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if res.DocsWritten != 3 {
		t.Fatalf("expected 3 docs, got %d", res.DocsWritten)
	}
	if res.ExtraFilesWritten == 0 {
		t.Fatalf("expected extra files written")
	}
	nonneg, err := os.ReadFile(filepath.Join(outDir, "Non-Negotiables.md"))
	if err != nil {
		t.Fatalf("read out: %v", err)
	}
	s := string(nonneg)
	if !strings.Contains(s, "<!-- GOV:BEGIN") || !strings.Contains(s, "sourceRef=v0.0.1") {
		t.Fatalf("expected managed markers with audit fields, got:\n%s", s)
	}
	if !strings.Contains(s, "CORE") || !strings.Contains(s, "PROFILE") {
		t.Fatalf("expected assembled fragments, got:\n%s", s)
	}
	if !strings.Contains(s, "## Local Addenda (project-owned)") {
		t.Fatalf("expected addenda section, got:\n%s", s)
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
	cmd := execCommand(exe, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v failed: %v\n%s", exe, args, err, string(out))
	}
}
