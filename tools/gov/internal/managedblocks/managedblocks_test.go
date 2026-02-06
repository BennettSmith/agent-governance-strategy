package managedblocks

import (
	"strings"
	"testing"
)

func TestFindBlocks_FindsWellFormedBlocks(t *testing.T) {
	doc := strings.Join([]string{
		"before",
		"<!-- GOV:BEGIN id=core version=v1 sha256=abc -->",
		"line1",
		"line2",
		"<!-- GOV:END id=core -->",
		"between",
		"<!-- GOV:BEGIN id=profile sha256=def -->",
		"",
		"<!-- GOV:END id=profile -->",
		"after",
	}, "\n")
	lines, _ := splitLines(doc)

	blocks, err := FindBlocks(lines, "GOV")
	if err != nil {
		t.Fatalf("FindBlocks: %v", err)
	}
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}
	if blocks[0].ID != "core" || blocks[1].ID != "profile" {
		t.Fatalf("unexpected ids: %+v", blocks)
	}
}

func TestFindBlocks_RejectsUnclosedBlock(t *testing.T) {
	lines, _ := splitLines("<!-- GOV:BEGIN id=x -->\nhi\n")
	_, err := FindBlocks(lines, "GOV")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "unclosed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReplaceBlock_ReplacesOnlyManagedRegionAndUpdatesSHA(t *testing.T) {
	doc := strings.Join([]string{
		"top",
		"<!-- GOV:BEGIN id=core version=v1 sha256=oldhash -->",
		"OLD",
		"<!-- GOV:END id=core -->",
		"",
		"## Local Addenda (project-owned)",
		"",
		"- keep me",
	}, "\n")

	out, err := ReplaceBlock(doc, ReplaceOptions{
		Prefix:     "GOV",
		BlockID:    "core",
		NewContent: "NEW1\nNEW2",
		MetaUpdates: map[string]string{
			"version":      "v2",
			"sourceRepo":   "git@example/repo",
			"sourceRef":    "v2.0.0",
			"sourceCommit": "deadbeef",
		},
	})
	if err != nil {
		t.Fatalf("ReplaceBlock: %v", err)
	}

	if !strings.Contains(out, "NEW1\nNEW2") {
		t.Fatalf("expected new content, got:\n%s", out)
	}
	if strings.Contains(out, "OLD") {
		t.Fatalf("expected old content removed, got:\n%s", out)
	}
	if !strings.Contains(out, "## Local Addenda (project-owned)\n\n- keep me") {
		t.Fatalf("expected local addenda preserved, got:\n%s", out)
	}
	if !strings.Contains(out, "version=v2") {
		t.Fatalf("expected version updated in marker, got:\n%s", out)
	}
	if !strings.Contains(out, "sha256=") {
		t.Fatalf("expected sha256 in marker, got:\n%s", out)
	}
}

func TestVerifyBlockSHA256_DetectsMismatch(t *testing.T) {
	doc := strings.Join([]string{
		"<!-- GOV:BEGIN id=core sha256=0000 -->",
		"content",
		"<!-- GOV:END id=core -->",
	}, "\n")
	err := VerifyBlockSHA256(doc, "GOV", "core")
	if err == nil {
		t.Fatalf("expected mismatch error")
	}
	if !strings.Contains(err.Error(), "mismatch") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVerifyBlockSHA256_PassesWhenCorrect(t *testing.T) {
	content := "line1\nline2"
	doc := strings.Join([]string{
		"<!-- GOV:BEGIN id=core sha256=" + SHA256Hex(content) + " -->",
		"line1",
		"line2",
		"<!-- GOV:END id=core -->",
	}, "\n")
	if err := VerifyBlockSHA256(doc, "GOV", "core"); err != nil {
		t.Fatalf("expected ok, got %v", err)
	}
}

