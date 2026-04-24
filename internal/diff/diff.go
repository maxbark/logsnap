package diff

import (
	"fmt"
	"strings"

	"github.com/logsnap/internal/snapshot"
)

// Result holds the comparison between two snapshots.
type Result struct {
	Added   []snapshot.Entry
	Removed []snapshot.Entry
	Changed []Change
}

// Change represents a log entry that exists in both snapshots but differs.
type Change struct {
	Key  string
	From snapshot.Entry
	To   snapshot.Entry
}

// Compare diffs two snapshots and returns a Result describing the differences.
func Compare(base, current *snapshot.Snapshot) *Result {
	result := &Result{}

	baseMap := indexEntries(base.Entries)
	currentMap := indexEntries(current.Entries)

	for key, baseEntry := range baseMap {
		currentEntry, exists := currentMap[key]
		if !exists {
			result.Removed = append(result.Removed, baseEntry)
			continue
		}
		if baseEntry.Message != currentEntry.Message || baseEntry.Level != currentEntry.Level {
			result.Changed = append(result.Changed, Change{
				Key:  key,
				From: baseEntry,
				To:   currentEntry,
			})
		}
	}

	for key, currentEntry := range currentMap {
		if _, exists := baseMap[key]; !exists {
			result.Added = append(result.Added, currentEntry)
		}
	}

	return result
}

// Summary returns a human-readable summary of the diff result.
func (r *Result) Summary() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Diff Summary: +%d added, -%d removed, ~%d changed\n",
		len(r.Added), len(r.Removed), len(r.Changed))

	for _, e := range r.Added {
		fmt.Fprintf(&sb, "  [+] [%s] %s\n", e.Level, e.Message)
	}
	for _, e := range r.Removed {
		fmt.Fprintf(&sb, "  [-] [%s] %s\n", e.Level, e.Message)
	}
	for _, c := range r.Changed {
		fmt.Fprintf(&sb, "  [~] %s: [%s] %q -> [%s] %q\n",
			c.Key, c.From.Level, c.From.Message, c.To.Level, c.To.Message)
	}

	return sb.String()
}

func indexEntries(entries []snapshot.Entry) map[string]snapshot.Entry {
	m := make(map[string]snapshot.Entry, len(entries))
	for _, e := range entries {
		key := fmt.Sprintf("%s::%s", e.Level, e.Source)
		m[key] = e
	}
	return m
}
