---
title: The Recency Engine
date: 2026-06-12
tags: [design, notes]
summary: How posts get bucketed into Today, This week, This month, and Older.
---
# The Recency Engine

A timeline is only useful if a reader can immediately tell what is fresh. The
recency engine sorts everything newest-first and then drops each entry into one
of four buckets relative to *now*:

1. **Today** — published on the current calendar day (or dated in the future).
2. **This week** — within the last seven days, but not today.
3. **This month** — within the last 31 days, but not this week.
4. **Older** — everything before that.

Empty buckets are simply omitted, so a quiet month does not leave hollow
headings on the page.

## Filtering

Beyond grouping, you can narrow the view. The `list` command accepts `--tag` to
restrict to a single tag and `--since` to set a date floor. Both compose, so
`longscroll list posts/ --tag design --since 2026-06-01` shows only design notes
from June onward.

The engine lives in its own package and is covered by unit tests for bucket
boundaries, sort order, and filtering — because off-by-one date math is exactly
the kind of bug that hides in plain sight.
