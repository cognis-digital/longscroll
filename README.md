# longscroll

A longform timeline / posts **static-site generator** with a built-in recency
engine. Point it at a folder of Markdown posts and it produces a self-contained
HTML site: a chronological timeline index plus one page per post.

No databases, no servers, no JavaScript framework — a single Go binary and the
standard library. The output is plain HTML with inline CSS that you can host
anywhere.


<!-- cognis:example:start -->
## 🔎 Example output

**Sample result format** _(illustrative values — run on your own data for real findings):_

```
{
  "id": "1234567890",
  "name": "John Doe",
  "email": "john.doe@example.com",
  "phone_numbers": [
    {"type": "mobile", "number": "+1 555-1234"},
    {"type": "work", "number": "+1 555-5678"}
  ]
}
```

<!-- cognis:example:end -->

## What it does

- **Ingests** a folder of Markdown posts carrying simple frontmatter
  (`title`, `date`, `tags`, `summary`).
- **Sorts** them newest-first.
- **Buckets** them with a recency engine: *Today*, *This week*, *This month*,
  *Older* (empty buckets are omitted).
- **Emits** a static site, or lists/validates posts from the command line.

## Install

```
go install github.com/cognis-digital/longscroll@latest
```

Or build from source:

```
git clone https://github.com/cognis-digital/longscroll
cd longscroll
go build ./...
```

Requires Go 1.22 or newer.

## Post format

Each post is a Markdown file beginning with a `---`-delimited frontmatter block:

```
---
title: My First Post
date: 2026-06-19
tags: [updates, intro]
summary: A short one-line description.
---
# Body heading

Your post body in Markdown...
```

| Field     | Required | Notes                                            |
|-----------|----------|--------------------------------------------------|
| `title`   | yes      | Non-empty headline.                              |
| `date`    | yes      | `YYYY-MM-DD`. Invalid dates fail validation.     |
| `tags`    | no       | `[a, b]` or plain `a b c`; lowercased and deduped. |
| `summary` | no       | One-line description shown on the index.         |

Filenames may carry a leading `YYYY-MM-DD-` prefix; it is stripped from the
generated slug for cleaner URLs.

### Supported Markdown

ATX headings (`#`–`####`), paragraphs, unordered (`-`/`*`) and ordered (`1.`)
lists, fenced code blocks (```` ``` ````), and inline `code`, **bold**,
*italic*, and [links](https://example.com). All text is HTML-escaped before
markup is applied, so the output is safe to embed directly.

## Usage

### Build a site

```
longscroll build examples/posts -o site --title "My Timeline"
```

Writes `site/index.html` and `site/posts/<slug>.html`.

| Flag      | Default      | Meaning                          |
|-----------|--------------|----------------------------------|
| `-o`      | `site`       | Output directory.                |
| `--title` | `longscroll` | Site title on the index page.    |

### List posts

```
longscroll list examples/posts
longscroll list examples/posts --tag design
longscroll list examples/posts --since 2026-06-01
```

Prints a table grouped into recency buckets, newest-first. `--tag` and
`--since` compose.

### Validate posts

```
longscroll validate examples/posts
```

Checks that every post has a present title and a valid `YYYY-MM-DD` date. Exits
non-zero if any file is invalid — handy in CI.

## Examples

The [`examples/posts`](examples/posts) folder contains three authored posts you
can build immediately:

```
longscroll build examples/posts -o site
```

## Development

```
go build ./...
go test ./...
go vet ./...
```

The codebase is organized into small internal packages:

- `internal/post` — frontmatter parsing, directory loading, the recency engine.
- `internal/markdown` — the dependency-free Markdown-to-HTML renderer.
- `internal/site` — HTML templating and site output.

## License

License: COCL 1.0

## Maintainer

Cognis Digital
