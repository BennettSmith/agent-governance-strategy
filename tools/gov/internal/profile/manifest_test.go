package profile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadManifest_ExtendsAndNormalizesPaths(t *testing.T) {
	tmp := t.TempDir()

	// base profile layout
	baseDir := filepath.Join(tmp, "Governance", "Profiles", "mobile-clean")
	iosDir := filepath.Join(tmp, "Governance", "Profiles", "mobile-clean-ios")
	coreDir := filepath.Join(tmp, "Governance", "Core")
	tplDir := filepath.Join(tmp, "Governance", "Templates", "UseCases")
	playbooksDir := filepath.Join(iosDir, "Playbooks")

	mkdirAll(t, coreDir, tplDir, baseDir, iosDir, playbooksDir)
	writeFile(t, filepath.Join(coreDir, "NonNegotiables.Core.md"), "core\n")
	writeFile(t, filepath.Join(baseDir, "NonNegotiables.Profile.md"), "profile\n")
	writeFile(t, filepath.Join(tplDir, "UseCase.Template.md"), "tpl\n")
	writeFile(t, filepath.Join(playbooksDir, "iOS-Packaging.md"), "playbook\n")

	baseManifestPath := filepath.Join(baseDir, "profile.yaml")
	writeFile(t, baseManifestPath, strings.TrimSpace(`
schemaVersion: 1
id: mobile-clean
documents:
  - output: Non-Negotiables.md
    fragments:
      - ../../Core/NonNegotiables.Core.md
      - ./NonNegotiables.Profile.md
templates:
  - source: ../../Templates/UseCases/UseCase.Template.md
    output: Docs/UseCases/UseCase.Template.md
`)+"\n")

	iosManifestPath := filepath.Join(iosDir, "profile.yaml")
	writeFile(t, iosManifestPath, strings.TrimSpace(`
schemaVersion: 1
id: mobile-clean-ios
extends:
  - ../mobile-clean/profile.yaml
playbooks:
  - source: ./Playbooks/iOS-Packaging.md
    output: Docs/Playbooks/iOS-Packaging.md
`)+"\n")

	m, err := LoadManifest(iosManifestPath)
	if err != nil {
		t.Fatalf("LoadManifest: %v", err)
	}
	if m.ID != "mobile-clean-ios" {
		t.Fatalf("expected overlay id, got %q", m.ID)
	}
	if len(m.Documents) != 1 {
		t.Fatalf("expected inherited documents, got %d", len(m.Documents))
	}
	if len(m.Documents[0].Fragments) != 2 {
		t.Fatalf("expected 2 fragments, got %d", len(m.Documents[0].Fragments))
	}
	for _, p := range m.Documents[0].Fragments {
		if !filepath.IsAbs(p) {
			t.Fatalf("expected absolute fragment path, got %q", p)
		}
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("expected fragment to exist: %q (%v)", p, err)
		}
	}
	if len(m.Templates) != 1 || !filepath.IsAbs(m.Templates[0].Source) {
		t.Fatalf("expected normalized template source, got %+v", m.Templates)
	}
	if len(m.Playbooks) != 1 || !filepath.IsAbs(m.Playbooks[0].Source) {
		t.Fatalf("expected normalized playbook source, got %+v", m.Playbooks)
	}
}

func TestLoadManifest_RejectsMissingID(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "profile.yaml")
	writeFile(t, path, "schemaVersion: 1\n")
	_, err := LoadManifest(path)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestLoadManifest_RejectsWrongSchemaVersion(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "profile.yaml")
	writeFile(t, path, "schemaVersion: 2\nid: x\n")
	_, err := LoadManifest(path)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func mkdirAll(t *testing.T, dirs ...string) {
	t.Helper()
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", d, err)
		}
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
