package source

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFetch_ErrorsWhenCacheDirIsFile(t *testing.T) {
	tmp := t.TempDir()
	cacheFile := filepath.Join(tmp, "cachefile")
	if err := os.WriteFile(cacheFile, []byte("x\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	_, err := Fetch(context.Background(), FetchOptions{
		RepoURL:  "does-not-matter",
		Ref:      "HEAD",
		CacheDir: cacheFile,
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestFetch_ErrorsOnBadRepoURL(t *testing.T) {
	tmp := t.TempDir()
	cache := filepath.Join(tmp, "cache")
	_, err := Fetch(context.Background(), FetchOptions{
		RepoURL:  filepath.Join(tmp, "missing-repo"),
		Ref:      "HEAD",
		CacheDir: cache,
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "git clone") {
		t.Fatalf("expected git clone error, got %v", err)
	}
}

func TestSanitizeRef_EmptyBecomesRef(t *testing.T) {
	if got := sanitizeRef("-"); got != "ref" {
		t.Fatalf("expected ref, got %q", got)
	}
}

