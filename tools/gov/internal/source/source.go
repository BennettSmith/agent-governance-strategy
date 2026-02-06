package source

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type ResolvedSource struct {
	CheckoutDir   string
	SourceRepo    string
	SourceRef     string
	SourceCommit  string
}

type FetchOptions struct {
	RepoURL  string
	Ref      string
	CacheDir string
}

func Fetch(ctx context.Context, opts FetchOptions) (ResolvedSource, error) {
	if strings.TrimSpace(opts.RepoURL) == "" {
		return ResolvedSource{}, errors.New("repo url is required")
	}
	if strings.TrimSpace(opts.Ref) == "" {
		return ResolvedSource{}, errors.New("ref is required")
	}
	if strings.TrimSpace(opts.CacheDir) == "" {
		return ResolvedSource{}, errors.New("cache dir is required")
	}

	repoKey := shortHash(opts.RepoURL)
	refKey := sanitizeRef(opts.Ref)
	checkoutDir := filepath.Join(opts.CacheDir, "sources", repoKey, refKey)

	if err := os.MkdirAll(filepath.Dir(checkoutDir), 0o755); err != nil {
		return ResolvedSource{}, fmt.Errorf("create cache parent: %w", err)
	}

	if _, err := os.Stat(filepath.Join(checkoutDir, ".git")); err != nil {
		// Clone.
		if err := runGit(ctx, "", "clone", "--no-checkout", opts.RepoURL, checkoutDir); err != nil {
			return ResolvedSource{}, fmt.Errorf("git clone: %w", err)
		}
	}

	// Fetch tags (best-effort) and checkout pinned ref.
	_ = runGit(ctx, checkoutDir, "fetch", "--tags", "--prune")
	if err := runGit(ctx, checkoutDir, "checkout", "--force", opts.Ref); err != nil {
		return ResolvedSource{}, fmt.Errorf("git checkout %s: %w", opts.Ref, err)
	}
	if err := runGit(ctx, checkoutDir, "reset", "--hard", "HEAD"); err != nil {
		return ResolvedSource{}, fmt.Errorf("git reset: %w", err)
	}
	commit, err := gitOutput(ctx, checkoutDir, "rev-parse", "HEAD")
	if err != nil {
		return ResolvedSource{}, fmt.Errorf("git rev-parse: %w", err)
	}

	return ResolvedSource{
		CheckoutDir:  checkoutDir,
		SourceRepo:   opts.RepoURL,
		SourceRef:    opts.Ref,
		SourceCommit: strings.TrimSpace(commit),
	}, nil
}

func runGit(ctx context.Context, dir string, args ...string) error {
	_, err := execGit(ctx, dir, args...)
	return err
}

func gitOutput(ctx context.Context, dir string, args ...string) (string, error) {
	out, err := execGit(ctx, dir, args...)
	return out, err
}

func execGit(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	b, err := cmd.CombinedOutput()
	out := string(b)
	if err != nil {
		return out, fmt.Errorf("%w: %s", err, strings.TrimSpace(out))
	}
	return out, nil
}

func shortHash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])[:12]
}

var refSanitizeRe = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func sanitizeRef(ref string) string {
	ref = strings.TrimSpace(ref)
	ref = strings.ReplaceAll(ref, "/", "_")
	ref = strings.ReplaceAll(ref, "\\", "_")
	ref = refSanitizeRe.ReplaceAllString(ref, "_")
	ref = strings.Trim(ref, "._-")
	if ref == "" {
		return "ref"
	}
	return ref
}

