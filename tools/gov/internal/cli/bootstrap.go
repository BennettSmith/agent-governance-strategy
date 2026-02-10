package cli

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"agent-governance-strategy/tools/gov/internal/builder"
	"agent-governance-strategy/tools/gov/internal/config"
	"agent-governance-strategy/tools/gov/internal/profile"
	"agent-governance-strategy/tools/gov/internal/source"
	"gopkg.in/yaml.v3"
)

type bootstrapConfigYAML struct {
	SchemaVersion int `yaml:"schemaVersion"`
	Source        struct {
		Repo    string `yaml:"repo"`
		Ref     string `yaml:"ref"`
		Profile string `yaml:"profile"`
	} `yaml:"source"`
	Paths struct {
		DocsRoot string `yaml:"docsRoot"`
		CacheDir string `yaml:"cacheDir,omitempty"`
	} `yaml:"paths"`
	Sync struct {
		ManagedBlockPrefix  string `yaml:"managedBlockPrefix"`
		LocalAddendaHeading string `yaml:"localAddendaHeading"`
	} `yaml:"sync"`
}

type profileInfo struct {
	ID          string
	Description string
}

func runBootstrap(subArgs []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("bootstrap", flag.ContinueOnError)
	fs.SetOutput(stderr)

	configPath := fs.String("config", defaultConfigPath, "path to write .governance/config.yaml")
	sourceRepo := fs.String("source-repo", "", "governance source repo URL/path (required in --non-interactive)")
	sourceRef := fs.String("source-ref", "", "governance source ref/tag/SHA (required in --non-interactive)")
	profileID := fs.String("profile", "", "governance profile id (required in --non-interactive)")
	docsRoot := fs.String("docs-root", ".", "docs root inside target repo")
	cacheDirFlag := fs.String("cache-dir", "", "cache dir for source fetch (optional)")

	managedBlockPrefix := fs.String("managed-block-prefix", "GOV", "managed block marker prefix")
	localAddendaHeading := fs.String("local-addenda-heading", "Local Addenda (project-owned)", "local addenda heading")

	nonInteractive := fs.Bool("non-interactive", false, "disable prompts; require flags")
	force := fs.Bool("force", false, "overwrite existing config if present")
	printOnly := fs.Bool("print", false, "print generated YAML to stdout instead of writing a file")
	runInit := fs.Bool("run-init", false, "run init after writing config")

	listProfiles := fs.Bool("list-profiles", false, "list profiles available in the selected source repo/ref and exit")

	if err := fs.Parse(subArgs); err != nil {
		return 2
	}

	fetchCacheDir := strings.TrimSpace(*cacheDirFlag)
	if fetchCacheDir == "" {
		cd, err := defaultCacheDir()
		if err != nil {
			fmt.Fprintf(stderr, "bootstrap error: cache dir: %v\n", err)
			return 2
		}
		fetchCacheDir = cd
	}

	interactive := isTTY(os.Stdin) && !*nonInteractive
	if interactive {
		if err := promptBootstrapValues(os.Stdin, stdout, stderr, fetchCacheDir, sourceRepo, sourceRef, profileID, docsRoot); err != nil {
			fmt.Fprintf(stderr, "bootstrap failed: %v\n", err)
			return 1
		}
	}

	// In non-interactive mode (and as a general safety net), require the essentials.
	if strings.TrimSpace(*sourceRepo) == "" || strings.TrimSpace(*sourceRef) == "" || strings.TrimSpace(*profileID) == "" {
		if *nonInteractive {
			fmt.Fprintln(stderr, "bootstrap error: --source-repo, --source-ref, and --profile are required in --non-interactive mode")
			return 2
		}
		fmt.Fprintln(stderr, "bootstrap error: missing required values (provide --source-repo, --source-ref, --profile or run interactively)")
		return 2
	}

	// Fetch source once for listing/validation (also ensures profile exists).
	src, err := source.Fetch(context.Background(), source.FetchOptions{
		RepoURL:  strings.TrimSpace(*sourceRepo),
		Ref:      strings.TrimSpace(*sourceRef),
		CacheDir: fetchCacheDir,
	})
	if err != nil {
		fmt.Fprintf(stderr, "bootstrap failed: fetch source: %v\n", err)
		return 1
	}

	if *listProfiles {
		infos, err := listProfilesFromCheckout(src.CheckoutDir)
		if err != nil {
			fmt.Fprintf(stderr, "bootstrap failed: list profiles: %v\n", err)
			return 1
		}
		for _, pi := range infos {
			desc := strings.TrimSpace(pi.Description)
			if desc != "" {
				fmt.Fprintf(stdout, "%s\t%s\n", pi.ID, desc)
			} else {
				fmt.Fprintf(stdout, "%s\n", pi.ID)
			}
		}
		return 0
	}

	manifestPath := filepath.Join(src.CheckoutDir, "Governance", "Profiles", strings.TrimSpace(*profileID), "profile.yaml")
	if _, err := os.Stat(manifestPath); err != nil {
		fmt.Fprintf(stderr, "bootstrap failed: profile %q not found at %s (%v)\n", strings.TrimSpace(*profileID), manifestPath, err)
		return 1
	}

	y := bootstrapConfigYAML{SchemaVersion: 1}
	y.Source.Repo = strings.TrimSpace(*sourceRepo)
	y.Source.Ref = strings.TrimSpace(*sourceRef)
	y.Source.Profile = strings.TrimSpace(*profileID)
	y.Paths.DocsRoot = strings.TrimSpace(*docsRoot)
	y.Paths.CacheDir = strings.TrimSpace(*cacheDirFlag) // only written if non-empty due to omitempty
	y.Sync.ManagedBlockPrefix = strings.TrimSpace(*managedBlockPrefix)
	y.Sync.LocalAddendaHeading = strings.TrimSpace(*localAddendaHeading)

	out, err := yaml.Marshal(y)
	if err != nil {
		fmt.Fprintf(stderr, "bootstrap failed: marshal yaml: %v\n", err)
		return 1
	}
	if len(out) == 0 || out[len(out)-1] != '\n' {
		out = append(out, '\n')
	}

	resolvedCfgPath := *configPath
	if !filepath.IsAbs(resolvedCfgPath) {
		if abs, err := filepath.Abs(resolvedCfgPath); err == nil {
			resolvedCfgPath = abs
		}
	}

	if *printOnly {
		_, _ = stdout.Write(out)
		return 0
	}

	if err := writeConfigFile(resolvedCfgPath, out, *force); err != nil {
		fmt.Fprintf(stderr, "bootstrap failed: write config: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "wrote config: %s\n", resolvedCfgPath)

	if *runInit {
		cfg, err := config.Load(resolvedCfgPath)
		if err != nil {
			fmt.Fprintf(stderr, "bootstrap failed: config error: %v\n", err)
			return 2
		}
		cacheDir, err := cfg.CacheDir()
		if err != nil {
			fmt.Fprintf(stderr, "bootstrap failed: cache dir error: %v\n", err)
			return 2
		}
		sourceRepoResolved := resolveRepoPathIfLocal(resolvedCfgPath, cfg.Source.Repo)
		repoRoot := repoRootForConfig(resolvedCfgPath)
		_, err = builder.Init(context.Background(), builder.InitOptions{
			RepoRoot:       repoRoot,
			DocsRoot:       cfg.Paths.DocsRoot,
			CacheDir:       cacheDir,
			SourceRepo:     sourceRepoResolved,
			SourceRef:      cfg.Source.Ref,
			ProfileID:      cfg.Source.Profile,
			MarkerPrefix:   cfg.Sync.ManagedBlockPrefix,
			AddendaHeading: cfg.Sync.LocalAddendaHeading,
		})
		if err != nil {
			fmt.Fprintf(stderr, "bootstrap failed: init failed: %v\n", err)
			return 1
		}
		fmt.Fprintln(stdout, "initialized governance docs")
	}

	return 0
}

func writeConfigFile(path string, content []byte, force bool) error {
	if _, err := os.Stat(path); err == nil && !force {
		return fmt.Errorf("config already exists (use --force to overwrite): %s", path)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}

func defaultCacheDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "govbuilder"), nil
}

func isTTY(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func listProfilesFromCheckout(checkoutDir string) ([]profileInfo, error) {
	root := filepath.Join(checkoutDir, "Governance", "Profiles")
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var infos []profileInfo
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		manifestPath := filepath.Join(root, e.Name(), "profile.yaml")
		if _, err := os.Stat(manifestPath); err != nil {
			continue
		}
		m, err := profile.LoadManifest(manifestPath)
		if err != nil {
			return nil, err
		}
		infos = append(infos, profileInfo{ID: m.ID, Description: m.Description})
	}
	sort.Slice(infos, func(i, j int) bool { return infos[i].ID < infos[j].ID })
	return infos, nil
}

func promptBootstrapValues(stdin io.Reader, stdout, stderr io.Writer, cacheDir string, sourceRepo, sourceRef, profileID, docsRoot *string) error {
	r := bufio.NewReader(stdin)

	if strings.TrimSpace(*sourceRepo) == "" {
		fmt.Fprint(stdout, "source repo (url/path): ")
		v, _ := r.ReadString('\n')
		*sourceRepo = strings.TrimSpace(v)
	}
	if strings.TrimSpace(*sourceRepo) == "" {
		return fmt.Errorf("source repo is required")
	}

	if strings.TrimSpace(*sourceRef) == "" {
		fmt.Fprint(stdout, "source ref (tag/sha, default HEAD): ")
		v, _ := r.ReadString('\n')
		v = strings.TrimSpace(v)
		if v == "" {
			v = "HEAD"
		}
		*sourceRef = v
	}

	// Profile selection via listing.
	if strings.TrimSpace(*profileID) == "" {
		src, err := source.Fetch(context.Background(), source.FetchOptions{
			RepoURL:  strings.TrimSpace(*sourceRepo),
			Ref:      strings.TrimSpace(*sourceRef),
			CacheDir: cacheDir,
		})
		if err != nil {
			return fmt.Errorf("fetch source: %w", err)
		}
		infos, err := listProfilesFromCheckout(src.CheckoutDir)
		if err != nil {
			return err
		}
		if len(infos) == 0 {
			return fmt.Errorf("no profiles found in source")
		}
		fmt.Fprintln(stderr, "available profiles:")
		for i, pi := range infos {
			desc := strings.TrimSpace(pi.Description)
			if desc != "" {
				fmt.Fprintf(stderr, "  %d) %s — %s\n", i+1, pi.ID, desc)
			} else {
				fmt.Fprintf(stderr, "  %d) %s\n", i+1, pi.ID)
			}
		}
		fmt.Fprint(stdout, "choose profile (number): ")
		v, _ := r.ReadString('\n')
		v = strings.TrimSpace(v)
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > len(infos) {
			return fmt.Errorf("invalid selection")
		}
		*profileID = infos[n-1].ID
	}

	if strings.TrimSpace(*docsRoot) == "" {
		*docsRoot = "."
	}
	return nil
}

