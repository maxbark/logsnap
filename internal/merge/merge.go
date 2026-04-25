// Package merge provides functionality for combining multiple snapshots into one.
package merge

import (
	"fmt"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
)

// Options controls how snapshots are merged.
type Options struct {
	// DeduplicateByMessage removes duplicate entries with the same message and level.
	DeduplicateByMessage bool
	// Label is used as the merged snapshot's deployment label.
	Label string
}

// DefaultOptions returns sensible defaults for merge operations.
func DefaultOptions() Options {
	return Options{
		DeduplicateByMessage: false,
		Label:                "merged",
	}
}

// Merge combines two or more snapshots into a single snapshot.
// Entries are appended in the order the snapshots are provided.
// If opts.DeduplicateByMessage is true, only the first occurrence of each
// (level, message) pair is kept.
func Merge(opts Options, snaps ...*snapshot.Snapshot) (*snapshot.Snapshot, error) {
	if len(snaps) < 2 {
		return nil, fmt.Errorf("merge requires at least 2 snapshots, got %d", len(snaps))
	}

	for i, s := range snaps {
		if s == nil {
			return nil, fmt.Errorf("snapshot at index %d is nil", i)
		}
	}

	label := opts.Label
	if label == "" {
		label = "merged"
	}

	out := snapshot.New(label, time.Now())

	seen := make(map[string]bool)

	for _, s := range snaps {
		for _, entry := range s.Entries {
			if opts.DeduplicateByMessage {
				key := entry.Level + "|" + entry.Message
				if seen[key] {
					continue
				}
				seen[key] = true
			}
			out.AddEntry(entry)
		}
	}

	return out, nil
}
