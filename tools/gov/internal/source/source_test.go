package source

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFetch_ClonesAndResolvesCommit_FromLocalRepo(t *testing.T) {
	ctx := context.Background()

	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "src")
	cache := filepath.Join(tmp, "cache")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(srcRepo, "README.md"), []byte("hello\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	mustRun(t, srcRepo, "git", "add", "README.md")
	mustRun(t, srcRepo, "git", "commit", "-m", "init")
	mustRun(t, srcRepo, "git", "tag", "v0.0.1")

	res, err := Fetch(ctx, FetchOptions{
		RepoURL:  srcRepo,
		Ref:      "v0.0.1",
		CacheDir: cache,
	})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if res.CheckoutDir == "" || res.SourceCommit == "" {
		t.Fatalf("expected checkout and commit, got %+v", res)
	}
	if _, err := os.Stat(filepath.Join(res.CheckoutDir, "README.md")); err != nil {
		t.Fatalf("expected checked out file: %v", err)
	}

	// Second fetch should reuse cache and still resolve.
	res2, err := Fetch(ctx, FetchOptions{
		RepoURL:  srcRepo,
		Ref:      "v0.0.1",
		CacheDir: cache,
	})
	if err != nil {
		t.Fatalf("Fetch 2: %v", err)
	}
	if res2.SourceCommit != res.SourceCommit {
		t.Fatalf("expected same commit, got %s vs %s", res2.SourceCommit, res.SourceCommit)
	}
}

func TestFetch_ValidatesRequiredFields(t *testing.T) {
	_, err := Fetch(context.Background(), FetchOptions{})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestSanitizeRef(t *testing.T) {
	got := sanitizeRef("release/v1.2.3")
	if strings.ContainsAny(got, "/\\") {
		t.Fatalf("expected sanitized ref, got %q", got)
	}
}

func TestFetch_ErrorsOnBadRef(t *testing.T) {
	ctx := context.Background()

	tmp := t.TempDir()
	srcRepo := filepath.Join(tmp, "src")
	cache := filepath.Join(tmp, "cache")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(srcRepo, "README.md"), []byte("hello\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	mustRun(t, srcRepo, "git", "add", "README.md")
	mustRun(t, srcRepo, "git", "commit", "-m", "init")

	_, err := Fetch(ctx, FetchOptions{RepoURL: srcRepo, Ref: "does-not-exist", CacheDir: cache})
	if err == nil {
		t.Fatalf("expected error")
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
