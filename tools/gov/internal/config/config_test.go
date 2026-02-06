package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad_ValidatesAndAppliesDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(strings.TrimSpace(`
schemaVersion: 1
source:
  repo: "/tmp/gov"
  ref: "v1.2.3"
  profile: "mobile-clean-ios"
`)), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.Paths.DocsRoot != "." {
		t.Fatalf("DocsRoot default: got %q", cfg.Paths.DocsRoot)
	}
	if cfg.Sync.ManagedBlockPrefix != "GOV" {
		t.Fatalf("ManagedBlockPrefix default: got %q", cfg.Sync.ManagedBlockPrefix)
	}
	if cfg.Sync.LocalAddendaHeading == "" {
		t.Fatalf("LocalAddendaHeading default missing")
	}
}

func TestLoad_RejectsMissingRequiredFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("schemaVersion: 1\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatalf("expected error")
	}
	if got := err.Error(); !strings.Contains(got, "source.repo") {
		t.Fatalf("expected source.repo error, got %q", got)
	}
}

