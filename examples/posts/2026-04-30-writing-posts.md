---
title: Writing Posts in Markdown
date: 2026-04-30
tags: [guide, notes]
summary: The Markdown subset longscroll understands, and how it stays safe.
---
# Writing Posts in Markdown

longscroll ships its own small Markdown renderer. It is deliberately limited to
the elements a longform post actually needs, which keeps the tool dependency
free and the output predictable.

## Supported elements

- Headings with `#` through `####`
- Paragraphs separated by a blank line
- Unordered lists with `-` or `*`
- Ordered lists with `1.`
- Fenced code blocks using triple backticks
- Inline `code`, **bold**, *italic*, and [links](https://example.com)

## Safety first

All text is HTML-escaped before any markup is applied, so a stray `<` or `&` in
your prose renders literally rather than breaking the page. Inline code spans
are extracted before emphasis is applied, which means a snippet like
`rate * time` keeps its asterisk instead of turning into accidental italics.

That is the entire contract. If you stick to these basics, what you write is
exactly what your readers see.
