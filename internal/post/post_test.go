package post

import (
	"testing"
	"time"
)

func TestParseValid(t *testing.T) {
	raw := `---
title: Hello World
date: 2026-06-19
tags: [intro, Updates]
summary: A short note.
---
# Body

Some text.
`
	p, err := Parse("hello", raw)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if p.Title != "Hello World" {
		t.Errorf("title = %q, want %q", p.Title, "Hello World")
	}
	if got := p.Date.Format(DateLayout); got != "2026-06-19" {
		t.Errorf("date = %q, want 2026-06-19", got)
	}
	if p.Summary != "A short note." {
		t.Errorf("summary = %q", p.Summary)
	}
	// Tags are lowercased, deduped, sorted.
	if len(p.Tags) != 2 || p.Tags[0] != "intro" || p.Tags[1] != "updates" {
		t.Errorf("tags = %v, want [intro updates]", p.Tags)
	}
	if p.Body == "" {
		t.Error("body should not be empty")
	}
}

func TestParseQuotedAndPlainTags(t *testing.T) {
	raw := `---
title: "Quoted Title"
date: 2026-01-02
tags: alpha beta alpha
---
body`
	p, err := Parse("q", raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Title != "Quoted Title" {
		t.Errorf("title = %q", p.Title)
	}
	if len(p.Tags) != 2 {
		t.Errorf("tags = %v, want 2 deduped", p.Tags)
	}
}

func TestParseErrors(t *testing.T) {
	cases := map[string]string{
		"no frontmatter": "just text with no frontmatter\n",
		"unterminated":   "---\ntitle: X\ndate: 2026-01-01\nbody without close",
		"missing title":  "---\ndate: 2026-01-01\n---\nbody",
		"missing date":   "---\ntitle: X\n---\nbody",
		"bad date":       "---\ntitle: X\ndate: June 1 2026\n---\nbody",
		"bad fm line":    "---\ntitle X\ndate: 2026-01-01\n---\nbody",
	}
	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := Parse("s", raw); err == nil {
				t.Errorf("expected error for %q, got nil", name)
			}
		})
	}
}

func TestParseCRLF(t *testing.T) {
	raw := "---\r\ntitle: Win\r\ndate: 2026-03-03\r\n---\r\nBody line\r\n"
	p, err := Parse("win", raw)
	if err != nil {
		t.Fatalf("CRLF parse error: %v", err)
	}
	if p.Title != "Win" {
		t.Errorf("title = %q", p.Title)
	}
}

func TestSortNewestFirst(t *testing.T) {
	mk := func(d, title string) Post {
		dt, _ := time.Parse(DateLayout, d)
		return Post{Date: dt, Title: title}
	}
	posts := []Post{
		mk("2026-01-01", "old"),
		mk("2026-06-01", "new"),
		mk("2026-03-01", "mid"),
		mk("2026-03-01", "amid"), // tie -> title order
	}
	SortNewestFirst(posts)
	want := []string{"new", "amid", "mid", "old"}
	for i, w := range want {
		if posts[i].Title != w {
			t.Errorf("pos %d = %q, want %q", i, posts[i].Title, w)
		}
	}
}

func TestHasTag(t *testing.T) {
	p := Post{Tags: []string{"defense", "osint"}}
	if !p.HasTag("OSINT") {
		t.Error("HasTag should be case-insensitive")
	}
	if p.HasTag("missing") {
		t.Error("HasTag returned true for absent tag")
	}
}
