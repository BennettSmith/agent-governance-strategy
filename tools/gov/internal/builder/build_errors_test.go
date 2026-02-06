package builder

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuild_ErrorsWhenProfileManifestMissing(t *testing.T) {
	ctx := context.Background()
	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "govsrc")
	outDir := filepath.Join(tmp, "out")
	cache := filepath.Join(tmp, "cache")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(srcRepo, "README.md"), []byte("x\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	mustRun(t, srcRepo, "git", "add", "README.md")
	mustRun(t, srcRepo, "git", "commit", "-m", "init")
	mustRun(t, srcRepo, "git", "tag", "v0.0.1")

	_, err := Build(ctx, BuildOptions{
		OutDir:       outDir,
		DocsRoot:     ".",
		CacheDir:     cache,
		SourceRepo:   srcRepo,
		SourceRef:    "v0.0.1",
		ProfileID:    "does-not-exist",
		MarkerPrefix: "GOV",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "profile.yaml") {
		t.Fatalf("expected manifest error, got %v", err)
	}
}
