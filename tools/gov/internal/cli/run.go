package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
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
		res, err := builder.Build(context.Background(), builder.BuildOptions{
			OutDir:         *outDir,
			DocsRoot:       cfg.Paths.DocsRoot,
			CacheDir:       cacheDir,
			SourceRepo:     cfg.Source.Repo,
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
		fmt.Fprintf(stdout, "TODO: init (config=%s)\n", *configPath)
		return 0
	case "sync":
		fmt.Fprintf(stdout, "TODO: sync (config=%s)\n", *configPath)
		return 0
	case "verify":
		fmt.Fprintf(stdout, "TODO: verify (config=%s)\n", *configPath)
		return 0
	default:
		fmt.Fprintf(stderr, "internal error: unhandled command %s\n", cmd)
		return 1
	}
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

