package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	SchemaVersion int      `yaml:"schemaVersion"`
	ID            string   `yaml:"id"`
	Description   string   `yaml:"description"`
	Extends       []string `yaml:"extends"`

	Documents []DocumentSpec `yaml:"documents"`
	Templates []FileSpec     `yaml:"templates"`
	Playbooks []FileSpec     `yaml:"playbooks"`
}

type DocumentSpec struct {
	Output    string   `yaml:"output"`
	Fragments []string `yaml:"fragments"`
}

type FileSpec struct {
	Source string `yaml:"source"`
	Output string `yaml:"output"`
}

func LoadManifest(path string) (Manifest, error) {
	baseDir := filepath.Dir(path)
	raw, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, err
	}
	var m Manifest
	if err := yaml.Unmarshal(raw, &m); err != nil {
		return Manifest{}, err
	}
	if m.SchemaVersion != 1 {
		return Manifest{}, fmt.Errorf("profile schemaVersion must be 1: %s", path)
	}
	if strings.TrimSpace(m.ID) == "" {
		return Manifest{}, fmt.Errorf("profile id is required: %s", path)
	}

	// Merge any base manifests first.
	for _, ext := range m.Extends {
		extPath := filepath.Clean(filepath.Join(baseDir, ext))
		base, err := LoadManifest(extPath)
		if err != nil {
			return Manifest{}, err
		}
		m = merge(base, m)
	}
	return m, nil
}

func merge(base, overlay Manifest) Manifest {
	out := base

	// Overlay identity fields.
	out.SchemaVersion = overlay.SchemaVersion
	out.ID = overlay.ID
	if strings.TrimSpace(overlay.Description) != "" {
		out.Description = overlay.Description
	}

	// Append specs.
	out.Documents = append(out.Documents, overlay.Documents...)
	out.Templates = append(out.Templates, overlay.Templates...)
	out.Playbooks = append(out.Playbooks, overlay.Playbooks...)

	return out
}

