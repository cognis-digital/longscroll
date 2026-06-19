---
title: Launching longscroll
date: 2026-06-19
tags: [release, intro]
summary: Why we built a tiny, dependency-free timeline generator.
---
# Launching longscroll

We wanted a way to publish a longform timeline of posts without dragging in a
heavyweight static-site framework. So we built **longscroll**: a single Go
binary that turns a folder of Markdown files into a self-contained website.

## Goals

- Zero third-party dependencies. Just the Go standard library.
- A *recency engine* that groups entries into Today, This week, This month and
  Older, so a reader can scan the timeline at a glance.
- Output you can host anywhere: plain HTML with inline CSS, no build step on the
  server side.

## How a post looks

Every post is a Markdown file with a small frontmatter block:

```
---
title: My Post
date: 2026-06-19
tags: [updates]
summary: One line.
---
Body goes here.
```

Run `longscroll build posts/ -o site/` and you get an `index.html` plus one page
per post. That is the whole idea. More to come.
