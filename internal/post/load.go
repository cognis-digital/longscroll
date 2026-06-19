package post

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// LoadResult holds successfully parsed posts plus per-file errors gathered
// during a directory load.
type LoadResult struct {
	Posts  []Post
	Errors []FileError
}

// FileError associates a parse failure with its source file.
type FileError struct {
	Path string
	Err  error
}

func (e FileError) Error() string {
	return fmt.Sprintf("%s: %v", e.Path, e.Err)
}

// slugFromName derives a slug from a markdown filename, stripping the extension
// and any leading "YYYY-MM-DD-" date prefix for readability.
func slugFromName(name string) string {
	base := strings.TrimSuffix(name, filepath.Ext(name))
	// Strip a leading ISO date prefix like 2026-06-19- if present.
	if len(base) > 11 && base[4] == '-' && base[7] == '-' && base[10] == '-' {
		if _, err := parseDatePrefix(base[:10]); err == nil {
			base = base[11:]
		}
	}
	return base
}

func parseDatePrefix(s string) (struct{}, error) {
	for i, r := range s {
		if i == 4 || i == 7 {
			if r != '-' {
				return struct{}{}, fmt.Errorf("bad")
			}
			continue
		}
		if r < '0' || r > '9' {
			return struct{}{}, fmt.Errorf("bad")
		}
	}
	return struct{}{}, nil
}

// LoadDir reads every *.md / *.markdown file directly inside dir, parses each,
// and returns parsed posts (sorted newest-first) together with any per-file
// errors. A missing directory is reported as a returned error.
func LoadDir(dir string) (LoadResult, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return LoadResult{}, fmt.Errorf("reading posts dir: %w", err)
	}

	var res LoadResult
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext == ".md" || ext == ".markdown" {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		path := filepath.Join(dir, name)
		raw, rerr := os.ReadFile(path)
		if rerr != nil {
			res.Errors = append(res.Errors, FileError{Path: path, Err: rerr})
			continue
		}
		p, perr := Parse(slugFromName(name), string(raw))
		if perr != nil {
			res.Errors = append(res.Errors, FileError{Path: path, Err: perr})
			continue
		}
		res.Posts = append(res.Posts, p)
	}

	SortNewestFirst(res.Posts)
	return res, nil
}
