package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SchemaVersion int `yaml:"schemaVersion"`

	Source SourceConfig `yaml:"source"`
	Paths  PathsConfig  `yaml:"paths"`
	Sync   SyncConfig   `yaml:"sync"`
}

type SourceConfig struct {
	Repo    string `yaml:"repo"`
	Ref     string `yaml:"ref"`
	Profile string `yaml:"profile"`
}

type PathsConfig struct {
	DocsRoot string `yaml:"docsRoot"`
	CacheDir string `yaml:"cacheDir"`
}

type SyncConfig struct {
	ManagedBlockPrefix  string `yaml:"managedBlockPrefix"`
	LocalAddendaHeading string `yaml:"localAddendaHeading"`
}

func Load(path string) (Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return Config{}, err
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	cfg = cfg.WithDefaults()
	return cfg, nil
}

func (c Config) Validate() error {
	var problems []string
	if c.SchemaVersion != 1 {
		problems = append(problems, "schemaVersion must be 1")
	}
	if strings.TrimSpace(c.Source.Repo) == "" {
		problems = append(problems, "source.repo is required")
	}
	if strings.TrimSpace(c.Source.Ref) == "" {
		problems = append(problems, "source.ref is required")
	}
	if strings.TrimSpace(c.Source.Profile) == "" {
		problems = append(problems, "source.profile is required")
	}
	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "; "))
	}
	return nil
}

func (c Config) WithDefaults() Config {
	if strings.TrimSpace(c.Paths.DocsRoot) == "" {
		c.Paths.DocsRoot = "."
	}
	if strings.TrimSpace(c.Sync.ManagedBlockPrefix) == "" {
		c.Sync.ManagedBlockPrefix = "GOV"
	}
	if strings.TrimSpace(c.Sync.LocalAddendaHeading) == "" {
		c.Sync.LocalAddendaHeading = "Local Addenda (project-owned)"
	}
	return c
}

func (c Config) CacheDir() (string, error) {
	if strings.TrimSpace(c.Paths.CacheDir) != "" {
		return expandHome(c.Paths.CacheDir), nil
	}
	base, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("user cache dir: %w", err)
	}
	return filepath.Join(base, "govbuilder"), nil
}

func expandHome(p string) string {
	if p == "~" {
		if h, err := os.UserHomeDir(); err == nil {
			return h
		}
		return p
	}
	if strings.HasPrefix(p, "~/") {
		if h, err := os.UserHomeDir(); err == nil {
			return filepath.Join(h, strings.TrimPrefix(p, "~/"))
		}
	}
	return p
}
