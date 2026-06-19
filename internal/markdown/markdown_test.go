package markdown

import (
	"strings"
	"testing"
)

func TestHeadings(t *testing.T) {
	got := Render("# Title\n## Sub")
	if !strings.Contains(got, "<h1>Title</h1>") {
		t.Errorf("missing h1 in %q", got)
	}
	if !strings.Contains(got, "<h2>Sub</h2>") {
		t.Errorf("missing h2 in %q", got)
	}
}

func TestParagraph(t *testing.T) {
	got := Render("Line one\nline two\n\nSecond para")
	if !strings.Contains(got, "<p>Line one line two</p>") {
		t.Errorf("paragraph join failed: %q", got)
	}
	if strings.Count(got, "<p>") != 2 {
		t.Errorf("expected 2 paragraphs, got %q", got)
	}
}

func TestUnorderedList(t *testing.T) {
	got := Render("- one\n- two")
	if !strings.Contains(got, "<ul>") || !strings.Contains(got, "<li>one</li>") {
		t.Errorf("unordered list failed: %q", got)
	}
}

func TestOrderedList(t *testing.T) {
	got := Render("1. first\n2. second")
	if !strings.Contains(got, "<ol>") || !strings.Contains(got, "<li>first</li>") {
		t.Errorf("ordered list failed: %q", got)
	}
}

func TestCodeFence(t *testing.T) {
	got := Render("```\nfmt.Println(\"hi\")\n```")
	if !strings.Contains(got, "<pre><code>") {
		t.Errorf("code fence failed: %q", got)
	}
	if !strings.Contains(got, "&#34;hi&#34;") && !strings.Contains(got, "&quot;hi&quot;") {
		t.Errorf("code not escaped: %q", got)
	}
}

func TestInlineLinkAndEmphasis(t *testing.T) {
	got := Render("See [the site](https://example.com) for **bold** and *italic*.")
	if !strings.Contains(got, `<a href="https://example.com"`) {
		t.Errorf("link failed: %q", got)
	}
	if !strings.Contains(got, "<strong>bold</strong>") {
		t.Errorf("bold failed: %q", got)
	}
	if !strings.Contains(got, "<em>italic</em>") {
		t.Errorf("italic failed: %q", got)
	}
}

func TestInlineCodeNotEmphasized(t *testing.T) {
	got := Render("Use `a*b*c` literally.")
	if !strings.Contains(got, "<code>a*b*c</code>") {
		t.Errorf("inline code span failed: %q", got)
	}
	if strings.Contains(got, "<em>") {
		t.Errorf("emphasis should not apply inside code span: %q", got)
	}
}

func TestEscaping(t *testing.T) {
	got := Render("A < B & C > D")
	if strings.Contains(got, "<B") {
		t.Errorf("html not escaped: %q", got)
	}
	if !strings.Contains(got, "&lt;") || !strings.Contains(got, "&amp;") {
		t.Errorf("expected escaped entities: %q", got)
	}
}
