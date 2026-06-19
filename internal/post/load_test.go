package post

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func TestLoadDir(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "2026-06-19-newest.md", "---\ntitle: Newest\ndate: 2026-06-19\n---\nbody")
	writeFile(t, dir, "older.md", "---\ntitle: Older\ndate: 2026-01-01\n---\nbody")
	writeFile(t, dir, "broken.md", "no frontmatter here")
	writeFile(t, dir, "ignored.txt", "not markdown")

	res, err := LoadDir(dir)
	if err != nil {
		t.Fatalf("LoadDir error: %v", err)
	}
	if len(res.Posts) != 2 {
		t.Errorf("got %d posts, want 2", len(res.Posts))
	}
	if len(res.Errors) != 1 {
		t.Errorf("got %d errors, want 1", len(res.Errors))
	}
	// Newest-first.
	if res.Posts[0].Title != "Newest" {
		t.Errorf("first post = %q, want Newest", res.Posts[0].Title)
	}
	// Date prefix stripped from slug.
	if res.Posts[0].Slug != "newest" {
		t.Errorf("slug = %q, want newest", res.Posts[0].Slug)
	}
}

func TestLoadDirMissing(t *testing.T) {
	if _, err := LoadDir(filepath.Join(t.TempDir(), "does-not-exist")); err == nil {
		t.Error("expected error for missing directory")
	}
}
