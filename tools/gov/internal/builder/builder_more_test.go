package builder

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuild_DefaultsDocsRootMarkerPrefixAndAddendaHeading(t *testing.T) {
	ctx := context.Background()
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	outDir := filepath.Join(tmp, "out")
	cache := filepath.Join(tmp, "cache")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "NonNegotiables.Core.md"), "CORE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "p", "NonNegotiables.Profile.md"), "PROFILE\n")
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
	mustRun(t, srcRepo, "git", "commit", "-m", "gov")
	mustRun(t, srcRepo, "git", "tag", "v0.0.1")

	_, err := Build(ctx, BuildOptions{
		OutDir:         outDir,
		DocsRoot:       "", // default to "."
		CacheDir:       cache,
		SourceRepo:     srcRepo,
		SourceRef:      "v0.0.1",
		ProfileID:      "p",
		MarkerPrefix:   "", // default to "GOV"
		AddendaHeading: "", // default heading
	})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	b, err := os.ReadFile(filepath.Join(outDir, "Non-Negotiables.md"))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	s := string(b)
	if !strings.Contains(s, "<!-- GOV:BEGIN") {
		t.Fatalf("expected default GOV marker prefix, got:\n%s", s)
	}
	if !strings.Contains(s, "## Local Addenda (project-owned)") {
		t.Fatalf("expected default addenda heading, got:\n%s", s)
	}
}

func TestBuild_ErrorsWhenTemplateSourceMissing(t *testing.T) {
	ctx := context.Background()
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	outDir := filepath.Join(tmp, "out")
	cache := filepath.Join(tmp, "cache")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "NonNegotiables.Core.md"), "CORE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "p", "NonNegotiables.Profile.md"), "PROFILE\n")
	// Template source path is intentionally missing.
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "p", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: p
documents:
  - output: Non-Negotiables.md
    fragments:
      - ../../Core/NonNegotiables.Core.md
      - ./NonNegotiables.Profile.md
templates:
  - source: ./Templates/Nope.md
    output: Docs/Nope.md
`)+"\n")

	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "gov")
	mustRun(t, srcRepo, "git", "tag", "v0.0.1")

	_, err := Build(ctx, BuildOptions{
		OutDir:     outDir,
		DocsRoot:   ".",
		CacheDir:   cache,
		SourceRepo: srcRepo,
		SourceRef:  "v0.0.1",
		ProfileID:  "p",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}
