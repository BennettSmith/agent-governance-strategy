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

// Build metadata. These are intended to be set at build time via:
// -ldflags "-X agent-governance-strategy/tools/gov/internal/cli.Version=... -X ...Commit=... -X ...Date=..."
var (
	Version = "dev"
	Commit  = ""
	Date    = ""
)

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) < 2 {
		printUsage(stderr)
		return 2
	}

	cmd := args[1]
	switch cmd {
	case "version", "-v", "--version":
		printVersion(stdout)
		return 0
	case "help", "-h", "--help":
		printUsage(stdout)
		return 0
	case "preflight":
		return runPreflight(args[2:], stdout, stderr)
	case "bootstrap":
		return runBootstrap(args[2:], stdout, stderr)
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

	resolvedConfigPath, autoDiscovered, err := resolveConfigPath(*configPath, subArgs)
	if err != nil {
		fmt.Fprintf(stderr, "config discovery error: %v\n", err)
		return 2
	}
	if autoDiscovered {
		fmt.Fprintf(stderr, "using config: %s\n", resolvedConfigPath)
	}

	switch cmd {
	case "build":
		if strings.TrimSpace(*outDir) == "" {
			fmt.Fprintln(stderr, "--out is required for build")
			return 2
		}
		cfg, err := config.Load(resolvedConfigPath)
		if err != nil {
			fmt.Fprintf(stderr, "config error: %v\n", err)
			return 2
		}
		cacheDir, err := cfg.CacheDir()
		if err != nil {
			fmt.Fprintf(stderr, "cache dir error: %v\n", err)
			return 2
		}
		sourceRepo := resolveRepoPathIfLocal(resolvedConfigPath, cfg.Source.Repo)
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
		cfg, err := config.Load(resolvedConfigPath)
		if err != nil {
			fmt.Fprintf(stderr, "config error: %v\n", err)
			return 2
		}
		cacheDir, err := cfg.CacheDir()
		if err != nil {
			fmt.Fprintf(stderr, "cache dir error: %v\n", err)
			return 2
		}
		sourceRepo := resolveRepoPathIfLocal(resolvedConfigPath, cfg.Source.Repo)
		repoRoot := repoRootForConfig(resolvedConfigPath)
		res, err := builder.Init(context.Background(), builder.InitOptions{
			RepoRoot:       repoRoot,
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
		cfg, err := config.Load(resolvedConfigPath)
		if err != nil {
			fmt.Fprintf(stderr, "config error: %v\n", err)
			return 2
		}
		cacheDir, err := cfg.CacheDir()
		if err != nil {
			fmt.Fprintf(stderr, "cache dir error: %v\n", err)
			return 2
		}
		sourceRepo := resolveRepoPathIfLocal(resolvedConfigPath, cfg.Source.Repo)
		repoRoot := repoRootForConfig(resolvedConfigPath)
		res, err := builder.Sync(context.Background(), builder.SyncOptions{
			RepoRoot:     repoRoot,
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
		cfg, err := config.Load(resolvedConfigPath)
		if err != nil {
			fmt.Fprintf(stderr, "config error: %v\n", err)
			return 2
		}
		cacheDir, err := cfg.CacheDir()
		if err != nil {
			fmt.Fprintf(stderr, "cache dir error: %v\n", err)
			return 2
		}
		sourceRepo := resolveRepoPathIfLocal(resolvedConfigPath, cfg.Source.Repo)
		repoRoot := repoRootForConfig(resolvedConfigPath)
		res, err := builder.Verify(context.Background(), builder.VerifyOptions{
			RepoRoot:     repoRoot,
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

func resolveConfigPath(configPath string, args []string) (string, bool, error) {
	if configFlagProvided(args) {
		return configPath, false, nil
	}
	found, ok, err := findNearestConfig(".")
	if err != nil {
		return "", false, err
	}
	if !ok {
		return configPath, false, nil
	}
	// Only log when the config was discovered above the current working directory.
	cwd, err := os.Getwd()
	if err == nil {
		localDefault := filepath.Clean(filepath.Join(cwd, defaultConfigPath))
		if filepath.Clean(found) == localDefault {
			return found, false, nil
		}
	}
	return found, true, nil
}

func configFlagProvided(args []string) bool {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "-config" || a == "--config" {
			return true
		}
		if strings.HasPrefix(a, "-config=") || strings.HasPrefix(a, "--config=") {
			return true
		}
	}
	return false
}

func findNearestConfig(startDir string) (string, bool, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", false, err
	}
	for {
		candidate := filepath.Join(dir, defaultConfigPath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false, nil
		}
		dir = parent
	}
}

func repoRootForConfig(configPath string) string {
	// Conventional layout: <repoRoot>/.governance/config.yaml
	cfgDir := filepath.Dir(configPath)
	if filepath.Base(cfgDir) == ".governance" {
		if abs, err := filepath.Abs(filepath.Dir(cfgDir)); err == nil {
			return abs
		}
		return filepath.Dir(cfgDir)
	}
	// If config lives elsewhere, treat its directory as the root.
	if abs, err := filepath.Abs(cfgDir); err == nil {
		return abs
	}
	return cfgDir
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
	fmt.Fprintln(w, "  version  Print the agent-gov version")
	fmt.Fprintln(w, "  preflight Run branch/baseline sanity checks")
	fmt.Fprintln(w, "  bootstrap Create .governance/config.yaml (optionally run init)")
	fmt.Fprintln(w, "  init     Initialize governance docs in this repo")
	fmt.Fprintln(w, "  sync     Update managed governance blocks in-place")
	fmt.Fprintln(w, "  verify   Verify managed governance blocks match expected content")
	fmt.Fprintln(w, "  build    Assemble governance bundle into an output folder")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Global options:")
	fmt.Fprintf(w, "  --config PATH   Path to config (default %s; auto-discovers upward when omitted)\n", defaultConfigPath)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Build options:")
	fmt.Fprintln(w, "  --out DIR       Output directory (required)")
}

func printVersion(w io.Writer) {
	v := strings.TrimSpace(Version)
	if v == "" {
		v = "dev"
	}
	c := strings.TrimSpace(Commit)
	if c != "" {
		fmt.Fprintf(w, "%s (%s)\n", v, c)
		return
	}
	fmt.Fprintln(w, v)
}
