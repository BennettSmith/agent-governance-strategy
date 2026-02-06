package builder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"agent-governance-strategy/tools/gov/internal/managedblocks"
	"agent-governance-strategy/tools/gov/internal/profile"
	"agent-governance-strategy/tools/gov/internal/source"
)

type InitOptions struct {
	RepoRoot string
	DocsRoot string

	CacheDir   string
	SourceRepo string
	SourceRef  string
	ProfileID  string

	MarkerPrefix   string
	AddendaHeading string
}

type InitResult struct {
	DocsWritten       int
	ExtraFilesWritten int
}

func Init(ctx context.Context, opts InitOptions) (InitResult, error) {
	outDir := filepath.Clean(opts.RepoRoot)
	res, err := Build(ctx, BuildOptions{
		OutDir:         outDir,
		DocsRoot:       opts.DocsRoot,
		CacheDir:       opts.CacheDir,
		SourceRepo:     opts.SourceRepo,
		SourceRef:      opts.SourceRef,
		ProfileID:      opts.ProfileID,
		MarkerPrefix:   opts.MarkerPrefix,
		AddendaHeading: opts.AddendaHeading,
	})
	return InitResult{DocsWritten: res.DocsWritten, ExtraFilesWritten: res.ExtraFilesWritten}, err
}

type SyncOptions struct {
	RepoRoot string
	DocsRoot string

	CacheDir   string
	SourceRepo string
	SourceRef  string
	ProfileID  string

	MarkerPrefix string
}

type SyncResult struct {
	DocsUpdated int
}

func Sync(ctx context.Context, opts SyncOptions) (SyncResult, error) {
	if strings.TrimSpace(opts.RepoRoot) == "" {
		opts.RepoRoot = "."
	}
	if strings.TrimSpace(opts.DocsRoot) == "" {
		opts.DocsRoot = "."
	}
	if strings.TrimSpace(opts.MarkerPrefix) == "" {
		opts.MarkerPrefix = "GOV"
	}

	src, err := source.Fetch(ctx, source.FetchOptions{
		RepoURL:  opts.SourceRepo,
		Ref:      opts.SourceRef,
		CacheDir: opts.CacheDir,
	})
	if err != nil {
		return SyncResult{}, err
	}
	manifestPath := filepath.Join(src.CheckoutDir, "Governance", "Profiles", opts.ProfileID, "profile.yaml")
	m, err := profile.LoadManifest(manifestPath)
	if err != nil {
		return SyncResult{}, err
	}

	targetBase := filepath.Clean(filepath.Join(opts.RepoRoot, opts.DocsRoot))
	updated := 0
	for _, doc := range m.Documents {
		targetPath := filepath.Join(targetBase, doc.Output)
		existing, err := os.ReadFile(targetPath)
		if err != nil {
			return SyncResult{}, fmt.Errorf("read target doc %s: %w", targetPath, err)
		}

		newContent, err := assembleFragments(doc.Fragments)
		if err != nil {
			return SyncResult{}, fmt.Errorf("assemble %s: %w", doc.Output, err)
		}
		blockID := managedBlockIDForDoc(doc.Output)
		out, err := managedblocks.ReplaceBlock(string(existing), managedblocks.ReplaceOptions{
			Prefix:     opts.MarkerPrefix,
			BlockID:    blockID,
			NewContent: newContent,
			MetaUpdates: map[string]string{
				"version":      src.SourceRef,
				"sourceRepo":   src.SourceRepo,
				"sourceRef":    src.SourceRef,
				"sourceCommit": src.SourceCommit,
			},
		})
		if err != nil {
			return SyncResult{}, fmt.Errorf("update %s: %w", targetPath, err)
		}
		if err := os.WriteFile(targetPath, []byte(out), 0o644); err != nil {
			return SyncResult{}, fmt.Errorf("write %s: %w", targetPath, err)
		}
		updated++
	}
	return SyncResult{DocsUpdated: updated}, nil
}

type VerifyOptions struct {
	RepoRoot string
	DocsRoot string

	CacheDir   string
	SourceRepo string
	SourceRef  string
	ProfileID  string

	MarkerPrefix string
}

type VerifyResult struct {
	OK     bool
	Issues []string
}

func Verify(ctx context.Context, opts VerifyOptions) (VerifyResult, error) {
	if strings.TrimSpace(opts.RepoRoot) == "" {
		opts.RepoRoot = "."
	}
	if strings.TrimSpace(opts.DocsRoot) == "" {
		opts.DocsRoot = "."
	}
	if strings.TrimSpace(opts.MarkerPrefix) == "" {
		opts.MarkerPrefix = "GOV"
	}

	src, err := source.Fetch(ctx, source.FetchOptions{
		RepoURL:  opts.SourceRepo,
		Ref:      opts.SourceRef,
		CacheDir: opts.CacheDir,
	})
	if err != nil {
		return VerifyResult{}, err
	}
	manifestPath := filepath.Join(src.CheckoutDir, "Governance", "Profiles", opts.ProfileID, "profile.yaml")
	m, err := profile.LoadManifest(manifestPath)
	if err != nil {
		return VerifyResult{}, err
	}

	targetBase := filepath.Clean(filepath.Join(opts.RepoRoot, opts.DocsRoot))
	var issues []string
	for _, doc := range m.Documents {
		targetPath := filepath.Join(targetBase, doc.Output)
		existing, err := os.ReadFile(targetPath)
		if err != nil {
			issues = append(issues, fmt.Sprintf("%s: missing or unreadable (%v)", doc.Output, err))
			continue
		}
		blockID := managedBlockIDForDoc(doc.Output)
		if err := managedblocks.VerifyBlockSHA256(string(existing), opts.MarkerPrefix, blockID); err != nil {
			issues = append(issues, fmt.Sprintf("%s: %v", doc.Output, err))
		}
	}
	return VerifyResult{OK: len(issues) == 0, Issues: issues}, nil
}
