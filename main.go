// longscroll is a longform timeline / posts static-site generator with a
// recency engine. It ingests a folder of Markdown posts carrying simple
// frontmatter (title, date, tags, summary), sorts them newest-first, and can
// emit a self-contained static HTML site, list posts grouped by recency, or
// validate a posts folder.
//
// Usage:
//
//	longscroll build <posts-dir> -o site/
//	longscroll list  <posts-dir> [--tag TAG] [--since YYYY-MM-DD]
//	longscroll validate <posts-dir>
package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/cognis-digital/longscroll/internal/post"
	"github.com/cognis-digital/longscroll/internal/site"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		usage()
		return 2
	}
	cmd := args[0]
	rest := args[1:]

	switch cmd {
	case "build":
		return cmdBuild(rest)
	case "list":
		return cmdList(rest)
	case "validate":
		return cmdValidate(rest)
	case "help", "-h", "--help":
		usage()
		return 0
	default:
		fmt.Fprintf(os.Stderr, "longscroll: unknown command %q\n\n", cmd)
		usage()
		return 2
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `longscroll - longform timeline static-site generator

Usage:
  longscroll build <posts-dir> -o <out-dir> [--title TITLE]
  longscroll list  <posts-dir> [--tag TAG] [--since YYYY-MM-DD]
  longscroll validate <posts-dir>

Each post is a Markdown file with a frontmatter block:
  ---
  title: My First Post
  date: 2026-06-19
  tags: [updates, intro]
  summary: A short one-line description.
  ---
  Body in Markdown...

License: COCL 1.0
`)
}

func cmdBuild(args []string) int {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	out := fs.String("o", "site", "output directory")
	title := fs.String("title", "longscroll", "site title shown on the index page")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	dir := fs.Arg(0)
	if dir == "" {
		fmt.Fprintln(os.Stderr, "build: missing <posts-dir>")
		return 2
	}

	res, err := post.LoadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "build: %v\n", err)
		return 1
	}
	for _, fe := range res.Errors {
		fmt.Fprintf(os.Stderr, "build: skipping %v\n", fe)
	}
	if len(res.Posts) == 0 {
		fmt.Fprintln(os.Stderr, "build: no valid posts found")
		return 1
	}

	br, err := site.Build(res.Posts, site.Options{OutDir: *out, SiteTitle: *title})
	if err != nil {
		fmt.Fprintf(os.Stderr, "build: %v\n", err)
		return 1
	}
	fmt.Printf("Built %d posts -> %s\n", len(br.PostPaths), br.IndexPath)
	return 0
}

func cmdList(args []string) int {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	tag := fs.String("tag", "", "only show posts carrying this tag")
	since := fs.String("since", "", "only show posts on or after this date (YYYY-MM-DD)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	dir := fs.Arg(0)
	if dir == "" {
		fmt.Fprintln(os.Stderr, "list: missing <posts-dir>")
		return 2
	}

	var sinceT time.Time
	if *since != "" {
		t, err := time.Parse(post.DateLayout, *since)
		if err != nil {
			fmt.Fprintf(os.Stderr, "list: invalid --since %q (want YYYY-MM-DD)\n", *since)
			return 2
		}
		sinceT = t
	}

	res, err := post.LoadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "list: %v\n", err)
		return 1
	}
	for _, fe := range res.Errors {
		fmt.Fprintf(os.Stderr, "list: skipping %v\n", fe)
	}

	posts := post.Filter(res.Posts, *tag, sinceT)
	if len(posts) == 0 {
		fmt.Println("(no matching posts)")
		return 0
	}

	now := time.Now()
	groups := post.GroupByRecency(posts, now)

	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	for _, g := range groups {
		fmt.Fprintf(w, "\n%s\n", g.Bucket.Label())
		for _, p := range g.Posts {
			tags := ""
			if len(p.Tags) > 0 {
				tags = "[" + joinTags(p.Tags) + "]"
			}
			fmt.Fprintf(w, "  %s\t%s\t%s\n", p.Date.Format(post.DateLayout), p.Title, tags)
		}
	}
	w.Flush()
	return 0
}

func joinTags(tags []string) string {
	out := ""
	for i, t := range tags {
		if i > 0 {
			out += ", "
		}
		out += t
	}
	return out
}

func cmdValidate(args []string) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	dir := fs.Arg(0)
	if dir == "" {
		fmt.Fprintln(os.Stderr, "validate: missing <posts-dir>")
		return 2
	}

	res, err := post.LoadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "validate: %v\n", err)
		return 1
	}

	if len(res.Errors) > 0 {
		for _, fe := range res.Errors {
			fmt.Fprintf(os.Stderr, "INVALID %v\n", fe)
		}
		fmt.Fprintf(os.Stderr, "validate: %d invalid, %d valid\n", len(res.Errors), len(res.Posts))
		return 1
	}

	if len(res.Posts) == 0 {
		fmt.Fprintln(os.Stderr, "validate: no posts found")
		return 1
	}

	for _, p := range res.Posts {
		fmt.Printf("OK %s (%s)\n", p.Slug, p.Date.Format(post.DateLayout))
	}
	fmt.Printf("validate: %d posts OK\n", len(res.Posts))
	return 0
}
