package post

import (
	"testing"
	"time"
)

func date(s string) time.Time {
	t, _ := time.Parse(DateLayout, s)
	return t
}

func TestBucketFor(t *testing.T) {
	ref := date("2026-06-19")
	cases := []struct {
		day  string
		want Bucket
	}{
		{"2026-06-19", Today},
		{"2026-06-20", Today}, // future treated as today
		{"2026-06-18", ThisWeek},
		{"2026-06-13", ThisWeek},  // 6 days
		{"2026-06-12", ThisMonth}, // 7 days
		{"2026-05-25", ThisMonth}, // 25 days
		{"2026-05-19", Older},     // 31 days
		{"2025-01-01", Older},
	}
	for _, c := range cases {
		if got := BucketFor(date(c.day), ref); got != c.want {
			t.Errorf("BucketFor(%s) = %v, want %v", c.day, got.Label(), c.want.Label())
		}
	}
}

func TestGroupByRecency(t *testing.T) {
	ref := date("2026-06-19")
	posts := []Post{
		{Title: "today", Date: date("2026-06-19")},
		{Title: "week", Date: date("2026-06-16")},
		{Title: "month", Date: date("2026-06-01")},
		{Title: "ancient", Date: date("2024-01-01")},
	}
	groups := GroupByRecency(posts, ref)
	if len(groups) != 4 {
		t.Fatalf("got %d groups, want 4", len(groups))
	}
	wantOrder := []Bucket{Today, ThisWeek, ThisMonth, Older}
	for i, g := range groups {
		if g.Bucket != wantOrder[i] {
			t.Errorf("group %d = %v, want %v", i, g.Bucket.Label(), wantOrder[i].Label())
		}
	}
}

func TestGroupByRecencyOmitsEmpty(t *testing.T) {
	ref := date("2026-06-19")
	posts := []Post{
		{Title: "today", Date: date("2026-06-19")},
		{Title: "ancient", Date: date("2024-01-01")},
	}
	groups := GroupByRecency(posts, ref)
	if len(groups) != 2 {
		t.Fatalf("got %d groups, want 2 (empty buckets omitted)", len(groups))
	}
	if groups[0].Bucket != Today || groups[1].Bucket != Older {
		t.Errorf("unexpected buckets: %v, %v", groups[0].Bucket.Label(), groups[1].Bucket.Label())
	}
}

func TestFilter(t *testing.T) {
	posts := []Post{
		{Title: "a", Date: date("2026-06-19"), Tags: []string{"x"}},
		{Title: "b", Date: date("2026-06-10"), Tags: []string{"y"}},
		{Title: "c", Date: date("2026-05-01"), Tags: []string{"x", "y"}},
	}

	byTag := Filter(posts, "x", time.Time{})
	if len(byTag) != 2 {
		t.Errorf("tag filter: got %d, want 2", len(byTag))
	}

	bySince := Filter(posts, "", date("2026-06-05"))
	if len(bySince) != 2 {
		t.Errorf("since filter: got %d, want 2", len(bySince))
	}

	both := Filter(posts, "y", date("2026-06-05"))
	if len(both) != 1 || both[0].Title != "b" {
		t.Errorf("combined filter: got %v, want [b]", both)
	}

	// Result is newest-first.
	all := Filter(posts, "", time.Time{})
	if all[0].Title != "a" || all[2].Title != "c" {
		t.Errorf("filter not sorted newest-first: %v", all)
	}
}

func TestAllTags(t *testing.T) {
	posts := []Post{
		{Tags: []string{"b", "a"}},
		{Tags: []string{"a", "c"}},
	}
	tags := AllTags(posts)
	want := []string{"a", "b", "c"}
	if len(tags) != 3 {
		t.Fatalf("got %v, want %v", tags, want)
	}
	for i, w := range want {
		if tags[i] != w {
			t.Errorf("tag %d = %q, want %q", i, tags[i], w)
		}
	}
}
