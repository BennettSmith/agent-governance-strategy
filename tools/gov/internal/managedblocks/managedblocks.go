package managedblocks

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
)

type Block struct {
	ID string

	// BeginLineIdx and EndLineIdx are indices into the document's line slice.
	BeginLineIdx int
	EndLineIdx   int

	// Meta are key/value pairs parsed from the BEGIN marker.
	Meta map[string]string
}

// ReplaceOptions controls block updates.
type ReplaceOptions struct {
	// Prefix is the marker namespace, e.g. "GOV" for "<!-- GOV:BEGIN ... -->".
	Prefix string
	// BlockID is the managed block id to update.
	BlockID string
	// NewContent is the replacement content (without BEGIN/END marker lines).
	NewContent string
	// MetaUpdates are applied to the BEGIN marker (merged over existing meta).
	// The "id" field is always preserved as BlockID.
	MetaUpdates map[string]string
	// ComputeSHA256 controls whether sha256 is recomputed from NewContent
	// and written into the BEGIN marker. Defaults to true.
	ComputeSHA256 bool
}

// ReplaceBlock replaces a managed block's content and updates its BEGIN marker.
// Only the region between BEGIN and END markers (exclusive) is modified.
func ReplaceBlock(doc string, opts ReplaceOptions) (string, error) {
	if strings.TrimSpace(opts.Prefix) == "" {
		return "", errors.New("prefix is required")
	}
	if strings.TrimSpace(opts.BlockID) == "" {
		return "", errors.New("block id is required")
	}
	computeHash := opts.ComputeSHA256
	if opts.MetaUpdates == nil {
		opts.MetaUpdates = map[string]string{}
	}
	if opts.ComputeSHA256 == false && opts.MetaUpdates["sha256"] == "" {
		// Caller explicitly disabled hashing and didn't provide sha256; allow.
	} else if opts.ComputeSHA256 == false {
		// Caller supplied sha256; ok.
	} else {
		// Default: compute.
		computeHash = true
	}

	lines, trailingNewline := splitLines(doc)

	blocks, err := FindBlocks(lines, opts.Prefix)
	if err != nil {
		return "", err
	}

	var b *Block
	for i := range blocks {
		if blocks[i].ID == opts.BlockID {
			if b != nil {
				return "", fmt.Errorf("multiple blocks found with id %q", opts.BlockID)
			}
			b = &blocks[i]
		}
	}
	if b == nil {
		return "", fmt.Errorf("block id %q not found", opts.BlockID)
	}
	if b.BeginLineIdx+1 > b.EndLineIdx {
		return "", fmt.Errorf("invalid block indices for %q", opts.BlockID)
	}

	meta := map[string]string{}
	for k, v := range b.Meta {
		meta[k] = v
	}
	meta["id"] = opts.BlockID
	for k, v := range opts.MetaUpdates {
		if strings.TrimSpace(k) == "" {
			continue
		}
		meta[k] = v
	}
	if computeHash {
		meta["sha256"] = SHA256Hex(opts.NewContent)
	}

	beginLine := FormatBeginMarker(opts.Prefix, meta)
	lines[b.BeginLineIdx] = beginLine

	newContentLines, _ := splitLines(opts.NewContent)
	// Replace between begin+1 and end (exclusive).
	lines = splice(lines, b.BeginLineIdx+1, b.EndLineIdx, newContentLines)

	return joinLines(lines, trailingNewline), nil
}

// VerifyBlockSHA256 verifies that the sha256 in the BEGIN marker matches the block content.
func VerifyBlockSHA256(doc, prefix, blockID string) error {
	lines, _ := splitLines(doc)
	blocks, err := FindBlocks(lines, prefix)
	if err != nil {
		return err
	}
	for _, b := range blocks {
		if b.ID != blockID {
			continue
		}
		want := strings.TrimSpace(b.Meta["sha256"])
		if want == "" {
			return fmt.Errorf("block %q is missing sha256 in BEGIN marker", blockID)
		}
		content := strings.Join(lines[b.BeginLineIdx+1:b.EndLineIdx], "\n")
		got := SHA256Hex(content)
		if got != want {
			return fmt.Errorf("block %q sha256 mismatch: got %s want %s", blockID, got, want)
		}
		return nil
	}
	return fmt.Errorf("block id %q not found", blockID)
}

// FindBlocks finds all well-formed managed blocks in the document.
func FindBlocks(lines []string, prefix string) ([]Block, error) {
	var out []Block
	open := map[string]Block{} // id -> block

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if meta, ok := parseMarker(trimmed, prefix, "BEGIN"); ok {
			id := strings.TrimSpace(meta["id"])
			if id == "" {
				return nil, fmt.Errorf("BEGIN marker at line %d missing id", i+1)
			}
			if _, exists := open[id]; exists {
				return nil, fmt.Errorf("nested/duplicate BEGIN for id %q at line %d", id, i+1)
			}
			open[id] = Block{
				ID:           id,
				BeginLineIdx: i,
				Meta:         meta,
			}
			continue
		}
		if meta, ok := parseMarker(trimmed, prefix, "END"); ok {
			id := strings.TrimSpace(meta["id"])
			if id == "" {
				return nil, fmt.Errorf("END marker at line %d missing id", i+1)
			}
			b, exists := open[id]
			if !exists {
				return nil, fmt.Errorf("END without BEGIN for id %q at line %d", id, i+1)
			}
			delete(open, id)
			b.EndLineIdx = i
			out = append(out, b)
		}
	}
	if len(open) > 0 {
		var ids []string
		for id := range open {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		return nil, fmt.Errorf("unclosed blocks: %s", strings.Join(ids, ", "))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].BeginLineIdx < out[j].BeginLineIdx })
	return out, nil
}

func SHA256Hex(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}

// FormatBeginMarker produces a deterministic BEGIN marker line.
func FormatBeginMarker(prefix string, meta map[string]string) string {
	fields := canonicalizeMeta(meta)
	return "<!-- " + prefix + ":BEGIN " + strings.Join(fields, " ") + " -->"
}

func canonicalizeMeta(meta map[string]string) []string {
	// Stable ordering: id first, then commonly-used fields, then the rest sorted.
	orderedKeys := []string{"id", "version", "sha256", "sourceRepo", "sourceRef", "sourceCommit"}
	seen := map[string]bool{}

	var fields []string
	for _, k := range orderedKeys {
		v, ok := meta[k]
		if !ok || strings.TrimSpace(v) == "" {
			continue
		}
		fields = append(fields, k+"="+v)
		seen[k] = true
	}
	var rest []string
	for k, v := range meta {
		if seen[k] {
			continue
		}
		if strings.TrimSpace(k) == "" || strings.TrimSpace(v) == "" {
			continue
		}
		rest = append(rest, k+"="+v)
	}
	sort.Strings(rest)
	fields = append(fields, rest...)
	return fields
}

func parseMarker(trimmedLine, prefix, kind string) (map[string]string, bool) {
	if !strings.HasPrefix(trimmedLine, "<!--") || !strings.HasSuffix(trimmedLine, "-->") {
		return nil, false
	}
	body := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(trimmedLine, "<!--"), "-->"))
	if body == "" {
		return nil, false
	}
	parts := strings.Fields(body)
	if len(parts) == 0 {
		return nil, false
	}
	wantHead := prefix + ":" + kind
	if parts[0] != wantHead {
		return nil, false
	}
	meta := map[string]string{}
	for _, tok := range parts[1:] {
		k, v, ok := strings.Cut(tok, "=")
		if !ok {
			continue
		}
		meta[k] = v
	}
	return meta, true
}

func splitLines(s string) ([]string, bool) {
	trailingNewline := strings.HasSuffix(s, "\n")
	if trailingNewline {
		s = strings.TrimSuffix(s, "\n")
	}
	if s == "" {
		return []string{""}, trailingNewline
	}
	return strings.Split(s, "\n"), trailingNewline
}

func joinLines(lines []string, trailingNewline bool) string {
	out := strings.Join(lines, "\n")
	if trailingNewline {
		out += "\n"
	}
	return out
}

func splice(lines []string, start, end int, replacement []string) []string {
	if start < 0 {
		start = 0
	}
	if end < start {
		end = start
	}
	out := make([]string, 0, len(lines)-(end-start)+len(replacement))
	out = append(out, lines[:start]...)
	out = append(out, replacement...)
	out = append(out, lines[end:]...)
	return out
}

