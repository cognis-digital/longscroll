package post

import (
	"sort"
	"time"
)

// Bucket is a recency period a post can fall into.
type Bucket int

const (
	// Today is the reference calendar day.
	Today Bucket = iota
	// ThisWeek is within the last 7 days but not today.
	ThisWeek
	// ThisMonth is within the last 31 days but not this week.
	ThisMonth
	// Older is everything earlier than this month.
	Older
)

// Label returns the human-readable name of the bucket.
func (b Bucket) Label() string {
	switch b {
	case Today:
		return "Today"
	case ThisWeek:
		return "This week"
	case ThisMonth:
		return "This month"
	default:
		return "Older"
	}
}

// BucketFor classifies a post date relative to the reference time ref. Dates in
// the future (after ref's day) are treated as Today. Comparison is done on
// calendar-day boundaries in ref's location.
func BucketFor(d, ref time.Time) Bucket {
	loc := ref.Location()
	refDay := truncDay(ref.In(loc))
	postDay := truncDay(d.In(loc))

	days := int(refDay.Sub(postDay).Hours() / 24)
	switch {
	case days <= 0:
		return Today
	case days < 7:
		return ThisWeek
	case days < 31:
		return ThisMonth
	default:
		return Older
	}
}

func truncDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// Group is a bucket together with the posts that fall into it.
type Group struct {
	Bucket Bucket
	Posts  []Post
}

// GroupByRecency buckets posts relative to ref and returns groups in canonical
// bucket order (Today, This week, This month, Older). Empty buckets are
// omitted. Posts within a group are sorted newest-first.
func GroupByRecency(posts []Post, ref time.Time) []Group {
	byBucket := map[Bucket][]Post{}
	for _, p := range posts {
		b := BucketFor(p.Date, ref)
		byBucket[b] = append(byBucket[b], p)
	}

	var groups []Group
	for _, b := range []Bucket{Today, ThisWeek, ThisMonth, Older} {
		ps := byBucket[b]
		if len(ps) == 0 {
			continue
		}
		SortNewestFirst(ps)
		groups = append(groups, Group{Bucket: b, Posts: ps})
	}
	return groups
}

// Filter returns the subset of posts matching the optional tag (case
// insensitive; "" means any) and published on or after since (zero time means
// no lower bound). The result is sorted newest-first.
func Filter(posts []Post, tag string, since time.Time) []Post {
	var out []Post
	for _, p := range posts {
		if tag != "" && !p.HasTag(tag) {
			continue
		}
		if !since.IsZero() && p.Date.Before(truncDay(since)) {
			continue
		}
		out = append(out, p)
	}
	SortNewestFirst(out)
	return out
}

// AllTags returns the sorted union of tags across posts.
func AllTags(posts []Post) []string {
	seen := map[string]bool{}
	var out []string
	for _, p := range posts {
		for _, t := range p.Tags {
			if !seen[t] {
				seen[t] = true
				out = append(out, t)
			}
		}
	}
	sort.Strings(out)
	return out
}
