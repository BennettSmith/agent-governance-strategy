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
	CheckoutDir  string
	SourceRepo   string
	SourceRef    string
	SourceCommit string
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

	// Fetch tags (best-effort) and resolve the ref to a commit SHA.
	_ = runGit(ctx, checkoutDir, "fetch", "--tags", "--prune")
	commit, err := resolveRemoteRef(ctx, opts.RepoURL, opts.Ref)
	if err != nil {
		return ResolvedSource{}, err
	}

	// Checkout the resolved commit for determinism.
	if err := runGit(ctx, checkoutDir, "checkout", "--force", commit); err != nil {
		return ResolvedSource{}, fmt.Errorf("git checkout %s: %w", commit, err)
	}
	if err := runGit(ctx, checkoutDir, "reset", "--hard", "HEAD"); err != nil {
		return ResolvedSource{}, fmt.Errorf("git reset: %w", err)
	}
	checkedOut, err := gitOutput(ctx, checkoutDir, "rev-parse", "HEAD")
	if err != nil {
		return ResolvedSource{}, fmt.Errorf("git rev-parse: %w", err)
	}

	return ResolvedSource{
		CheckoutDir:  checkoutDir,
		SourceRepo:   opts.RepoURL,
		SourceRef:    opts.Ref,
		SourceCommit: strings.TrimSpace(checkedOut),
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

func resolveRemoteRef(ctx context.Context, repoURL, ref string) (string, error) {
	out, err := execGit(ctx, "", "ls-remote", repoURL, ref, ref+"^{}")
	if err != nil {
		return "", fmt.Errorf("git ls-remote %s %s: %w", repoURL, ref, err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		return "", fmt.Errorf("ref %q not found in %s", ref, repoURL)
	}
	// Prefer peeled annotated tag (ends with ^{}).
	var commit string
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		sha := fields[0]
		name := fields[1]
		if strings.HasSuffix(name, "^{}") {
			commit = sha
			break
		}
		if commit == "" {
			commit = sha
		}
	}
	if commit == "" {
		return "", fmt.Errorf("could not resolve ref %q in %s", ref, repoURL)
	}
	return commit, nil
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
