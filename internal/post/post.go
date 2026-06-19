// Package post implements parsing of longform Markdown posts that carry a
// minimal frontmatter block, and the recency engine that groups posts by
// period relative to a reference time.
package post

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Post is a single parsed longform entry.
type Post struct {
	// Slug is a filesystem-safe identifier derived from the source filename.
	Slug string
	// Title is the human-readable headline (frontmatter "title").
	Title string
	// Date is the publication date (frontmatter "date", parsed from DateLayout).
	Date time.Time
	// Tags is the set of lowercase tags (frontmatter "tags").
	Tags []string
	// Summary is a one-line description (frontmatter "summary"), may be empty.
	Summary string
	// Body is the raw Markdown body following the frontmatter block.
	Body string
}

// DateLayout is the canonical date format accepted in frontmatter.
const DateLayout = "2006-01-02"

// HasTag reports whether the post carries the given tag (case-insensitive).
func (p Post) HasTag(tag string) bool {
	tag = strings.ToLower(strings.TrimSpace(tag))
	for _, t := range p.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// Parse parses a single post from its raw file contents. The slug should be
// supplied by the caller (typically derived from the filename). A valid post
// must carry a frontmatter block delimited by lines of exactly "---", and the
// frontmatter must contain a non-empty "title" and a "date" parseable as
// DateLayout.
func Parse(slug, raw string) (Post, error) {
	fm, body, err := splitFrontmatter(raw)
	if err != nil {
		return Post{}, err
	}

	p := Post{Slug: slug, Body: strings.TrimSpace(body)}

	fields, err := parseFrontmatter(fm)
	if err != nil {
		return Post{}, err
	}

	p.Title = strings.TrimSpace(fields["title"])
	if p.Title == "" {
		return Post{}, fmt.Errorf("missing required frontmatter field: title")
	}

	rawDate := strings.TrimSpace(fields["date"])
	if rawDate == "" {
		return Post{}, fmt.Errorf("missing required frontmatter field: date")
	}
	d, derr := time.Parse(DateLayout, rawDate)
	if derr != nil {
		return Post{}, fmt.Errorf("invalid date %q (want %s): %w", rawDate, DateLayout, derr)
	}
	p.Date = d

	p.Summary = strings.TrimSpace(fields["summary"])
	p.Tags = parseTags(fields["tags"])

	return p, nil
}

// splitFrontmatter separates the leading "---"-delimited frontmatter block from
// the body. The very first non-empty content must be the opening "---".
func splitFrontmatter(raw string) (fm, body string, err error) {
	// Normalize CRLF so the parser is robust on Windows-authored files.
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	trimmed := strings.TrimLeft(raw, "\n")
	if !strings.HasPrefix(trimmed, "---\n") && trimmed != "---" {
		return "", "", fmt.Errorf("missing frontmatter: file must start with a '---' line")
	}

	// Drop the opening delimiter.
	rest := strings.TrimPrefix(trimmed, "---\n")
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return "", "", fmt.Errorf("unterminated frontmatter: missing closing '---' line")
	}
	fm = rest[:idx]
	// After the closing delimiter, skip to end of that line.
	after := rest[idx+len("\n---"):]
	if nl := strings.IndexByte(after, '\n'); nl >= 0 {
		body = after[nl+1:]
	} else {
		body = ""
	}
	return fm, body, nil
}

// parseFrontmatter parses simple "key: value" lines. Blank lines and lines
// beginning with '#' are ignored. Duplicate keys take the last value.
func parseFrontmatter(fm string) (map[string]string, error) {
	fields := map[string]string{}
	for i, line := range strings.Split(fm, "\n") {
		s := strings.TrimSpace(line)
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}
		colon := strings.IndexByte(s, ':')
		if colon < 0 {
			return nil, fmt.Errorf("frontmatter line %d: expected 'key: value', got %q", i+1, s)
		}
		key := strings.ToLower(strings.TrimSpace(s[:colon]))
		val := strings.TrimSpace(s[colon+1:])
		val = strings.Trim(val, `"'`)
		if key == "" {
			return nil, fmt.Errorf("frontmatter line %d: empty key", i+1)
		}
		fields[key] = val
	}
	return fields, nil
}

// parseTags splits a tags value. Accepts a bracketed list ("[a, b]") or a plain
// comma/space separated list. Returns a deduplicated, lowercased, sorted slice.
func parseTags(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	raw = strings.TrimPrefix(raw, "[")
	raw = strings.TrimSuffix(raw, "]")
	fields := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t'
	})
	seen := map[string]bool{}
	var out []string
	for _, f := range fields {
		t := strings.ToLower(strings.Trim(strings.TrimSpace(f), `"'`))
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// SortNewestFirst orders posts by date descending; ties broken by title to
// keep output deterministic.
func SortNewestFirst(posts []Post) {
	sort.SliceStable(posts, func(i, j int) bool {
		if posts[i].Date.Equal(posts[j].Date) {
			return posts[i].Title < posts[j].Title
		}
		return posts[i].Date.After(posts[j].Date)
	})
}
