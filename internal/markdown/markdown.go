// Package markdown implements a small, dependency-free Markdown-to-HTML
// renderer covering the subset longscroll needs: ATX headings, paragraphs,
// unordered and ordered lists, fenced code blocks, plus inline emphasis,
// inline code and links. All text is HTML-escaped before output, so the result
// is safe to embed directly.
package markdown

import (
	"html"
	"regexp"
	"strings"
)

var (
	linkRe   = regexp.MustCompile(`\[([^\]]+)\]\(([^)\s]+)\)`)
	boldRe   = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	italicRe = regexp.MustCompile(`\*([^*]+)\*`)
)

// Render converts Markdown source into an HTML fragment.
func Render(src string) string {
	src = strings.ReplaceAll(src, "\r\n", "\n")
	lines := strings.Split(src, "\n")

	var out strings.Builder
	i := 0
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		switch {
		case trimmed == "":
			i++

		case strings.HasPrefix(trimmed, "```"):
			i = renderFence(&out, lines, i)

		case isHeading(trimmed):
			renderHeading(&out, trimmed)
			i++

		case isUnordered(trimmed):
			i = renderList(&out, lines, i, false)

		case isOrdered(trimmed):
			i = renderList(&out, lines, i, true)

		default:
			i = renderParagraph(&out, lines, i)
		}
	}
	return out.String()
}

func isHeading(s string) bool {
	return strings.HasPrefix(s, "# ") ||
		strings.HasPrefix(s, "## ") ||
		strings.HasPrefix(s, "### ") ||
		strings.HasPrefix(s, "#### ")
}

func isUnordered(s string) bool {
	return strings.HasPrefix(s, "- ") || strings.HasPrefix(s, "* ")
}

var orderedRe = regexp.MustCompile(`^\d+\.\s+`)

func isOrdered(s string) bool {
	return orderedRe.MatchString(s)
}

func renderHeading(out *strings.Builder, s string) {
	level := 0
	for level < len(s) && s[level] == '#' {
		level++
	}
	text := strings.TrimSpace(s[level:])
	out.WriteString("<h")
	out.WriteByte(byte('0' + level))
	out.WriteByte('>')
	out.WriteString(inline(text))
	out.WriteString("</h")
	out.WriteByte(byte('0' + level))
	out.WriteString(">\n")
}

// renderFence emits a code block; start points at the opening ``` line.
// Returns the index just past the closing fence.
func renderFence(out *strings.Builder, lines []string, start int) int {
	i := start + 1
	var body []string
	for i < len(lines) {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "```") {
			i++
			break
		}
		body = append(body, lines[i])
		i++
	}
	out.WriteString("<pre><code>")
	out.WriteString(html.EscapeString(strings.Join(body, "\n")))
	out.WriteString("</code></pre>\n")
	return i
}

// renderList emits a <ul> or <ol>; start points at the first item.
func renderList(out *strings.Builder, lines []string, start int, ordered bool) int {
	tag := "ul"
	if ordered {
		tag = "ol"
	}
	out.WriteString("<")
	out.WriteString(tag)
	out.WriteString(">\n")

	i := start
	for i < len(lines) {
		t := strings.TrimSpace(lines[i])
		var item string
		if ordered && isOrdered(t) {
			item = orderedRe.ReplaceAllString(t, "")
		} else if !ordered && isUnordered(t) {
			item = strings.TrimSpace(t[2:])
		} else {
			break
		}
		out.WriteString("<li>")
		out.WriteString(inline(item))
		out.WriteString("</li>\n")
		i++
	}

	out.WriteString("</")
	out.WriteString(tag)
	out.WriteString(">\n")
	return i
}

// renderParagraph consumes consecutive non-blank, non-structural lines.
func renderParagraph(out *strings.Builder, lines []string, start int) int {
	i := start
	var buf []string
	for i < len(lines) {
		t := strings.TrimSpace(lines[i])
		if t == "" || isHeading(t) || isUnordered(t) || isOrdered(t) ||
			strings.HasPrefix(t, "```") {
			break
		}
		buf = append(buf, t)
		i++
	}
	if len(buf) > 0 {
		out.WriteString("<p>")
		out.WriteString(inline(strings.Join(buf, " ")))
		out.WriteString("</p>\n")
	}
	return i
}

// inline escapes text then applies inline markup. Inline code spans are
// extracted first so their contents are not treated as emphasis or links.
func inline(s string) string {
	type span struct{ html string }
	var spans []span
	var b strings.Builder

	// Replace `code` spans with placeholders.
	for {
		open := strings.IndexByte(s, '`')
		if open < 0 {
			b.WriteString(s)
			break
		}
		close := strings.IndexByte(s[open+1:], '`')
		if close < 0 {
			b.WriteString(s)
			break
		}
		close += open + 1
		b.WriteString(s[:open])
		code := s[open+1 : close]
		b.WriteString("\x00")
		b.WriteByte(byte('0' + len(spans)%10))
		spans = append(spans, span{html: "<code>" + html.EscapeString(code) + "</code>"})
		s = s[close+1:]
	}

	text := html.EscapeString(b.String())

	// Links: [text](url) — escape both parts.
	text = linkRe.ReplaceAllStringFunc(text, func(m string) string {
		parts := linkRe.FindStringSubmatch(m)
		label := parts[1]
		url := parts[2]
		return `<a href="` + url + `" rel="noopener noreferrer">` + label + `</a>`
	})

	text = boldRe.ReplaceAllString(text, "<strong>$1</strong>")
	text = italicRe.ReplaceAllString(text, "<em>$1</em>")

	// Restore code spans.
	for idx, sp := range spans {
		ph := "\x00" + string(byte('0'+idx%10))
		text = strings.Replace(text, ph, sp.html, 1)
	}
	return text
}
