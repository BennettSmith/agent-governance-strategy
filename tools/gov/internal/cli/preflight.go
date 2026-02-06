package cli

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type stringSliceFlag []string

func (s *stringSliceFlag) String() string { return strings.Join(*s, ",") }
func (s *stringSliceFlag) Set(v string) error {
	*s = append(*s, v)
	return nil
}

type planFrontmatter struct {
	Branch string `yaml:"branch"`
	Status string `yaml:"status"`
}

func runPreflight(subArgs []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("preflight", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var require stringSliceFlag
	fs.Var(&require, "require", "required path relative to repo root (repeatable)")
	activePlan := fs.String("active-plan", "", "path to active plan file (optional)")

	if err := fs.Parse(subArgs); err != nil {
		return 2
	}

	cfgPath, ok, err := findNearestConfig(".")
	if err != nil {
		fmt.Fprintf(stderr, "preflight error: find config: %v\n", err)
		return 2
	}
	if !ok {
		fmt.Fprintf(stderr, "preflight failed: could not find %s upward from current directory\n", defaultConfigPath)
		return 1
	}

	repoRoot := repoRootForConfig(cfgPath)
	branch, err := gitCurrentBranch(repoRoot)
	if err != nil {
		fmt.Fprintf(stderr, "preflight error: git branch: %v\n", err)
		return 2
	}
	if branch == "HEAD" || strings.TrimSpace(branch) == "" {
		fmt.Fprintln(stderr, "preflight failed: detached HEAD (create/switch to a feature branch)")
		return 1
	}
	if branch == "main" {
		fmt.Fprintln(stderr, "preflight failed: on main (create/switch to a feature branch)")
		return 1
	}

	plansDir := filepath.Join(repoRoot, "Docs", "Plans")
	activeBranch := ""
	if strings.TrimSpace(*activePlan) != "" {
		ab, err := branchForPlanPath(repoRoot, *activePlan)
		if err != nil {
			fmt.Fprintf(stderr, "preflight error: active plan: %v\n", err)
			return 2
		}
		activeBranch = ab
	} else {
		ab, err := findActiveBranchFromPlans(plansDir)
		if err != nil {
			fmt.Fprintf(stderr, "preflight error: active plan scan: %v\n", err)
			return 2
		}
		activeBranch = ab
	}

	knownBranches, err := listPlannedBranches(plansDir)
	if err != nil {
		fmt.Fprintf(stderr, "preflight error: plan scan: %v\n", err)
		return 2
	}
	if activeBranch != "" {
		delete(knownBranches, activeBranch)
	}
	if knownBranches[branch] {
		if activeBranch != "" {
			fmt.Fprintf(stderr, "preflight failed: branch %q belongs to a different plan (active plan: %q)\n", branch, activeBranch)
		} else {
			fmt.Fprintf(stderr, "preflight failed: branch %q appears to belong to an existing plan\n", branch)
		}
		return 1
	}

	for _, rel := range require {
		rel = strings.TrimSpace(rel)
		if rel == "" {
			continue
		}
		p := filepath.Join(repoRoot, rel)
		if _, err := os.Stat(p); err != nil {
			fmt.Fprintf(stderr, "preflight failed: required path missing: %s (%v)\n", rel, err)
			return 1
		}
	}

	fmt.Fprintln(stdout, "ok")
	return 0
}

func gitCurrentBranch(dir string) (string, error) {
	ctx := context.Background()
	cmd := execCommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git rev-parse: %v (%s)", err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

func listPlannedBranches(plansDir string) (map[string]bool, error) {
	branches := make(map[string]bool)
	entries, err := walkPlanFiles(plansDir)
	if err != nil {
		return nil, err
	}
	for _, p := range entries {
		b, _, err := planMetadataForPath(plansDir, p)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(b) != "" {
			branches[b] = true
		}
	}
	return branches, nil
}

func findActiveBranchFromPlans(plansDir string) (string, error) {
	entries, err := walkPlanFiles(plansDir)
	if err != nil {
		return "", err
	}
	active := ""
	for _, p := range entries {
		b, status, err := planMetadataForPath(plansDir, p)
		if err != nil {
			return "", err
		}
		if strings.EqualFold(strings.TrimSpace(status), "active") && strings.TrimSpace(b) != "" {
			if active != "" && active != b {
				return "", fmt.Errorf("multiple active plans detected (%q and %q)", active, b)
			}
			active = b
		}
	}
	return active, nil
}

func branchForPlanPath(repoRoot, planPath string) (string, error) {
	abs := planPath
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(repoRoot, planPath)
	}
	plansDir := filepath.Join(repoRoot, "Docs", "Plans")
	b, _, err := planMetadataForPath(plansDir, abs)
	return b, err
}

func walkPlanFiles(plansDir string) ([]string, error) {
	var out []string
	info, err := os.Stat(plansDir)
	if err != nil || !info.IsDir() {
		// No plans directory is treated as "no known branches" (preflight still checks branch != main).
		return nil, nil
	}
	err = filepath.WalkDir(plansDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}
		// Skip templates so they don't become fake branches.
		if strings.Contains(strings.ToLower(d.Name()), "template") {
			return nil
		}
		out = append(out, path)
		return nil
	})
	return out, err
}

func planMetadataForPath(plansDir, planPath string) (branch string, status string, err error) {
	b, st, ok, err := readPlanFrontmatter(planPath)
	if err != nil {
		return "", "", err
	}
	if ok && strings.TrimSpace(b) != "" {
		return strings.TrimSpace(b), strings.TrimSpace(st), nil
	}
	// Fallback: derive from relative path under Docs/Plans.
	rel, err := filepath.Rel(plansDir, planPath)
	if err != nil {
		return "", "", err
	}
	rel = strings.TrimSuffix(rel, filepath.Ext(rel))
	rel = filepath.ToSlash(rel)
	return rel, strings.TrimSpace(st), nil
}

func readPlanFrontmatter(planPath string) (branch string, status string, ok bool, err error) {
	b, err := os.ReadFile(planPath)
	if err != nil {
		return "", "", false, err
	}
	// Frontmatter must be at the very start of the file.
	if !bytes.HasPrefix(b, []byte("---\n")) && !bytes.HasPrefix(b, []byte("---\r\n")) {
		return "", "", false, nil
	}
	// Find the next line that is exactly "---".
	parts := bytes.SplitN(b, []byte("\n---"), 2)
	if len(parts) < 2 {
		return "", "", false, nil
	}
	// parts[0] begins with "---\n", so drop the leading fence.
	fm := bytes.TrimPrefix(parts[0], []byte("---\n"))
	fm = bytes.TrimPrefix(fm, []byte("---\r\n"))
	var meta planFrontmatter
	if err := yaml.Unmarshal(fm, &meta); err != nil {
		return "", "", false, nil
	}
	return meta.Branch, meta.Status, true, nil
}

