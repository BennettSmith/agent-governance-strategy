package builder

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitSyncVerify_PreservesLocalAddenda(t *testing.T) {
	ctx := context.Background()
	tmp := t.TempDir()

	// governance source repo fixture
	srcRepo := filepath.Join(tmp, "govsrc")
	cache := filepath.Join(tmp, "cache")
	target := filepath.Join(tmp, "target")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	// v0.0.1 fragments
	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "NonNegotiables.Core.md"), "COREv1\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "NonNegotiables.Profile.md"), "PROFILEv1\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "Constitution.Core.md"), "CONSTv1\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "Constitution.Profile.md"), "C-PROFILEv1\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "Architecture.Profile.md"), "ARCHv1\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: backend-go-hex
documents:
  - output: Non-Negotiables.md
    fragments:
      - ../../Core/NonNegotiables.Core.md
      - ./NonNegotiables.Profile.md
  - output: Constitution.md
    fragments:
      - ../../Core/Constitution.Core.md
      - ./Constitution.Profile.md
  - output: Architecture.md
    fragments:
      - ./Architecture.Profile.md
`)+"\n")

	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "v1")
	mustRun(t, srcRepo, "git", "tag", "v0.0.1")

	// init into target
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	_, err := Init(ctx, InitOptions{
		RepoRoot: target,
		DocsRoot: ".",
		CacheDir: cache,
		SourceRepo: srcRepo,
		SourceRef: "v0.0.1",
		ProfileID: "backend-go-hex",
		MarkerPrefix: "GOV",
		AddendaHeading: "Local Addenda (project-owned)",
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	// add local addenda
	nonnegPath := filepath.Join(target, "Non-Negotiables.md")
	b, err := os.ReadFile(nonnegPath)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	withAddenda := string(b) + "\n- local override\n"
	if err := os.WriteFile(nonnegPath, []byte(withAddenda), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	// v0.0.2 change source fragments
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "NonNegotiables.Profile.md"), "PROFILEv2\n")
	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "v2")
	mustRun(t, srcRepo, "git", "tag", "v0.0.2")

	// sync to v0.0.2 (managed block should update, addenda preserved)
	_, err = Sync(ctx, SyncOptions{
		RepoRoot: target,
		DocsRoot: ".",
		CacheDir: cache,
		SourceRepo: srcRepo,
		SourceRef: "v0.0.2",
		ProfileID: "backend-go-hex",
		MarkerPrefix: "GOV",
	})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}

	updated, err := os.ReadFile(nonnegPath)
	if err != nil {
		t.Fatalf("read updated: %v", err)
	}
	s := string(updated)
	if !strings.Contains(s, "PROFILEv2") {
		t.Fatalf("expected managed content updated, got:\n%s", s)
	}
	if !strings.Contains(s, "- local override") {
		t.Fatalf("expected local addenda preserved, got:\n%s", s)
	}

	// verify should pass
	vr, err := Verify(ctx, VerifyOptions{
		RepoRoot: target,
		DocsRoot: ".",
		CacheDir: cache,
		SourceRepo: srcRepo,
		SourceRef: "v0.0.2",
		ProfileID: "backend-go-hex",
		MarkerPrefix: "GOV",
	})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if !vr.OK {
		t.Fatalf("expected ok, issues=%v", vr.Issues)
	}

	// Tamper inside the managed block and ensure verify reports an issue.
	tampered := strings.Replace(s, "PROFILEv2", "PROFILE_TAMPERED", 1)
	if err := os.WriteFile(nonnegPath, []byte(tampered), 0o644); err != nil {
		t.Fatalf("write tampered: %v", err)
	}
	vr2, err := Verify(ctx, VerifyOptions{
		RepoRoot: target,
		DocsRoot: ".",
		CacheDir: cache,
		SourceRepo: srcRepo,
		SourceRef: "v0.0.2",
		ProfileID: "backend-go-hex",
		MarkerPrefix: "GOV",
	})
	if err != nil {
		t.Fatalf("Verify tampered: %v", err)
	}
	if vr2.OK {
		t.Fatalf("expected not ok after tamper")
	}
}

func TestSync_ErrorsWhenTargetDocMissing(t *testing.T) {
	ctx := context.Background()
	tmp := t.TempDir()

	srcRepo := filepath.Join(tmp, "govsrc")
	cache := filepath.Join(tmp, "cache")
	target := filepath.Join(tmp, "target")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "NonNegotiables.Core.md"), "CORE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "NonNegotiables.Profile.md"), "PROFILE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: backend-go-hex
documents:
  - output: Non-Negotiables.md
    fragments:
      - ../../Core/NonNegotiables.Core.md
      - ./NonNegotiables.Profile.md
`)+"\n")
	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "v1")
	mustRun(t, srcRepo, "git", "tag", "v0.0.1")

	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	// No init; doc is missing.
	_, err := Sync(ctx, SyncOptions{
		RepoRoot: target,
		DocsRoot: ".",
		CacheDir: cache,
		SourceRepo: srcRepo,
		SourceRef: "v0.0.1",
		ProfileID: "backend-go-hex",
		MarkerPrefix: "GOV",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestVerify_ReportsIssueWhenShaMissing(t *testing.T) {
	ctx := context.Background()
	tmp := t.TempDir()

	srcRepo := filepath.Join(tmp, "govsrc")
	cache := filepath.Join(tmp, "cache")
	target := filepath.Join(tmp, "target")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "NonNegotiables.Core.md"), "CORE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "NonNegotiables.Profile.md"), "PROFILE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: backend-go-hex
documents:
  - output: Non-Negotiables.md
    fragments:
      - ../../Core/NonNegotiables.Core.md
      - ./NonNegotiables.Profile.md
`)+"\n")
	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "v1")
	mustRun(t, srcRepo, "git", "tag", "v0.0.1")

	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	// Write a doc with a managed block missing sha256.
	doc := strings.Join([]string{
		"<!-- GOV:BEGIN id=doc-non-negotiables -->",
		"CORE",
		"<!-- GOV:END id=doc-non-negotiables -->",
	}, "\n")
	if err := os.WriteFile(filepath.Join(target, "Non-Negotiables.md"), []byte(doc), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	vr, err := Verify(ctx, VerifyOptions{
		RepoRoot: target,
		DocsRoot: "",
		CacheDir: cache,
		SourceRepo: srcRepo,
		SourceRef: "v0.0.1",
		ProfileID: "backend-go-hex",
		MarkerPrefix: "",
	})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if vr.OK {
		t.Fatalf("expected not ok")
	}
}

func TestSync_ErrorsWhenManagedBlockMissingAndUsesDefaults(t *testing.T) {
	ctx := context.Background()
	tmp := t.TempDir()

	srcRepo := filepath.Join(tmp, "govsrc")
	cache := filepath.Join(tmp, "cache")
	target := filepath.Join(tmp, "target")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")

	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "NonNegotiables.Core.md"), "CORE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "NonNegotiables.Profile.md"), "PROFILE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: backend-go-hex
documents:
  - output: Non-Negotiables.md
    fragments:
      - ../../Core/NonNegotiables.Core.md
      - ./NonNegotiables.Profile.md
`)+"\n")
	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "v1")
	mustRun(t, srcRepo, "git", "tag", "v0.0.1")

	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	_, err := Init(ctx, InitOptions{
		RepoRoot: target,
		DocsRoot: ".",
		CacheDir: cache,
		SourceRepo: srcRepo,
		SourceRef: "v0.0.1",
		ProfileID: "backend-go-hex",
		MarkerPrefix: "GOV",
		AddendaHeading: "Local Addenda (project-owned)",
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	// Break the managed block markers.
	if err := os.WriteFile(filepath.Join(target, "Non-Negotiables.md"), []byte("no markers\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	// Use defaults by running relative to target.
	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	if err := os.Chdir(target); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	_, err = Sync(ctx, SyncOptions{
		RepoRoot:     "",
		DocsRoot:     "",
		CacheDir:     cache,
		SourceRepo:   srcRepo,
		SourceRef:    "v0.0.1",
		ProfileID:    "backend-go-hex",
		MarkerPrefix: "",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestVerify_ReportsMissingDocIssue(t *testing.T) {
	ctx := context.Background()
	tmp := t.TempDir()

	srcRepo := filepath.Join(tmp, "govsrc")
	cache := filepath.Join(tmp, "cache")
	target := filepath.Join(tmp, "target")

	mustRun(t, tmp, "git", "init", srcRepo)
	mustRun(t, srcRepo, "git", "config", "user.email", "test@example.com")
	mustRun(t, srcRepo, "git", "config", "user.name", "Test")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Core", "NonNegotiables.Core.md"), "CORE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "NonNegotiables.Profile.md"), "PROFILE\n")
	writeFile(t, filepath.Join(srcRepo, "Governance", "Profiles", "backend-go-hex", "profile.yaml"), strings.TrimSpace(`
schemaVersion: 1
id: backend-go-hex
documents:
  - output: Non-Negotiables.md
    fragments:
      - ../../Core/NonNegotiables.Core.md
      - ./NonNegotiables.Profile.md
`)+"\n")
	mustRun(t, srcRepo, "git", "add", ".")
	mustRun(t, srcRepo, "git", "commit", "-m", "v1")
	mustRun(t, srcRepo, "git", "tag", "v0.0.1")

	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	vr, err := Verify(ctx, VerifyOptions{
		RepoRoot: target,
		DocsRoot: ".",
		CacheDir: cache,
		SourceRepo: srcRepo,
		SourceRef: "v0.0.1",
		ProfileID: "backend-go-hex",
		MarkerPrefix: "GOV",
	})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if vr.OK {
		t.Fatalf("expected not ok")
	}
	if len(vr.Issues) == 0 {
		t.Fatalf("expected issues")
	}
}

