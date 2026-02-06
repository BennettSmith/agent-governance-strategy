package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"agent-governance-strategy/tools/gov/internal/builder"
	"agent-governance-strategy/tools/gov/internal/config"
)

const defaultConfigPath = ".governance/config.yaml"

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) < 2 {
		printUsage(stderr)
		return 2
	}

	cmd := args[1]
	switch cmd {
	case "help", "-h", "--help":
		printUsage(stdout)
		return 0
	case "init", "sync", "verify", "build":
		return runSubcommand(cmd, args[2:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n\n", cmd)
		printUsage(stderr)
		return 2
	}
}

func runSubcommand(cmd string, subArgs []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
	fs.SetOutput(stderr)

	configPath := fs.String("config", defaultConfigPath, "path to .governance/config.yaml")
	outDir := fs.String("out", "", "output directory (build only)")

	if err := fs.Parse(subArgs); err != nil {
		// flag package already printed the error/usage.
		return 2
	}

	switch cmd {
	case "build":
		if strings.TrimSpace(*outDir) == "" {
			fmt.Fprintln(stderr, "--out is required for build")
			return 2
		}
		cfg, err := config.Load(*configPath)
		if err != nil {
			fmt.Fprintf(stderr, "config error: %v\n", err)
			return 2
		}
		cacheDir, err := cfg.CacheDir()
		if err != nil {
			fmt.Fprintf(stderr, "cache dir error: %v\n", err)
			return 2
		}
		sourceRepo := resolveRepoPathIfLocal(*configPath, cfg.Source.Repo)
		res, err := builder.Build(context.Background(), builder.BuildOptions{
			OutDir:         *outDir,
			DocsRoot:       cfg.Paths.DocsRoot,
			CacheDir:       cacheDir,
			SourceRepo:     sourceRepo,
			SourceRef:      cfg.Source.Ref,
			ProfileID:      cfg.Source.Profile,
			MarkerPrefix:   cfg.Sync.ManagedBlockPrefix,
			AddendaHeading: cfg.Sync.LocalAddendaHeading,
		})
		if err != nil {
			fmt.Fprintf(stderr, "build failed: %v\n", err)
			return 1
		}
		fmt.Fprintf(stdout, "built %d doc(s) and %d file(s) (sourceCommit=%s)\n", res.DocsWritten, res.ExtraFilesWritten, res.SourceCommit)
		return 0
	case "init":
		cfg, err := config.Load(*configPath)
		if err != nil {
			fmt.Fprintf(stderr, "config error: %v\n", err)
			return 2
		}
		cacheDir, err := cfg.CacheDir()
		if err != nil {
			fmt.Fprintf(stderr, "cache dir error: %v\n", err)
			return 2
		}
		sourceRepo := resolveRepoPathIfLocal(*configPath, cfg.Source.Repo)
		res, err := builder.Init(context.Background(), builder.InitOptions{
			RepoRoot:       ".",
			DocsRoot:       cfg.Paths.DocsRoot,
			CacheDir:       cacheDir,
			SourceRepo:     sourceRepo,
			SourceRef:      cfg.Source.Ref,
			ProfileID:      cfg.Source.Profile,
			MarkerPrefix:   cfg.Sync.ManagedBlockPrefix,
			AddendaHeading: cfg.Sync.LocalAddendaHeading,
		})
		if err != nil {
			fmt.Fprintf(stderr, "init failed: %v\n", err)
			return 1
		}
		fmt.Fprintf(stdout, "initialized %d doc(s) and %d file(s)\n", res.DocsWritten, res.ExtraFilesWritten)
		return 0
	case "sync":
		cfg, err := config.Load(*configPath)
		if err != nil {
			fmt.Fprintf(stderr, "config error: %v\n", err)
			return 2
		}
		cacheDir, err := cfg.CacheDir()
		if err != nil {
			fmt.Fprintf(stderr, "cache dir error: %v\n", err)
			return 2
		}
		sourceRepo := resolveRepoPathIfLocal(*configPath, cfg.Source.Repo)
		res, err := builder.Sync(context.Background(), builder.SyncOptions{
			RepoRoot:     ".",
			DocsRoot:     cfg.Paths.DocsRoot,
			CacheDir:     cacheDir,
			SourceRepo:   sourceRepo,
			SourceRef:    cfg.Source.Ref,
			ProfileID:    cfg.Source.Profile,
			MarkerPrefix: cfg.Sync.ManagedBlockPrefix,
		})
		if err != nil {
			fmt.Fprintf(stderr, "sync failed: %v\n", err)
			return 1
		}
		fmt.Fprintf(stdout, "synced %d doc(s)\n", res.DocsUpdated)
		return 0
	case "verify":
		cfg, err := config.Load(*configPath)
		if err != nil {
			fmt.Fprintf(stderr, "config error: %v\n", err)
			return 2
		}
		cacheDir, err := cfg.CacheDir()
		if err != nil {
			fmt.Fprintf(stderr, "cache dir error: %v\n", err)
			return 2
		}
		sourceRepo := resolveRepoPathIfLocal(*configPath, cfg.Source.Repo)
		res, err := builder.Verify(context.Background(), builder.VerifyOptions{
			RepoRoot:     ".",
			DocsRoot:     cfg.Paths.DocsRoot,
			CacheDir:     cacheDir,
			SourceRepo:   sourceRepo,
			SourceRef:    cfg.Source.Ref,
			ProfileID:    cfg.Source.Profile,
			MarkerPrefix: cfg.Sync.ManagedBlockPrefix,
		})
		if err != nil {
			fmt.Fprintf(stderr, "verify failed: %v\n", err)
			return 1
		}
		if res.OK {
			fmt.Fprintln(stdout, "ok")
			return 0
		}
		fmt.Fprintf(stderr, "verification failed: %d issue(s)\n", len(res.Issues))
		for _, issue := range res.Issues {
			fmt.Fprintf(stderr, "- %s\n", issue)
		}
		return 1
	default:
		fmt.Fprintf(stderr, "internal error: unhandled command %s\n", cmd)
		return 1
	}
}

// If the configured repo looks like a local path (relative) and exists on disk
// relative to the config file directory, resolve it to an absolute path.
// Remote URLs will be left unchanged.
func resolveRepoPathIfLocal(configPath, repo string) string {
	repo = strings.TrimSpace(repo)
	if repo == "" {
		return repo
	}
	// Try resolving relative to the config file directory.
	base := filepath.Dir(configPath)
	if !filepath.IsAbs(repo) {
		candidate := filepath.Clean(filepath.Join(base, repo))
		if _, err := os.Stat(candidate); err == nil {
			if abs, err := filepath.Abs(candidate); err == nil {
				return abs
			}
			return candidate
		}
	} else {
		if _, err := os.Stat(repo); err == nil {
			return repo
		}
	}
	// Fallback: if it exists as-is (relative to current working directory), accept it.
	if _, err := os.Stat(repo); err == nil {
		if abs, err := filepath.Abs(repo); err == nil {
			return abs
		}
		return repo
	}
	return repo
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "agent-gov <command> [--config PATH] [options]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  init     Initialize governance docs in this repo")
	fmt.Fprintln(w, "  sync     Update managed governance blocks in-place")
	fmt.Fprintln(w, "  verify   Verify managed governance blocks match expected content")
	fmt.Fprintln(w, "  build    Assemble governance bundle into an output folder")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Global options:")
	fmt.Fprintf(w, "  --config PATH   Path to config (default %s)\n", defaultConfigPath)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Build options:")
	fmt.Fprintln(w, "  --out DIR       Output directory (required)")
}
