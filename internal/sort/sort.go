// Package sort provides utilities for ordering snapshot entries.
package sort

import (
	gosort "sort"
	"strings"

	"github.com/yourorg/logsnap/internal/snapshot"
)

// Field represents a sortable field on a log entry.
type Field string

const (
	FieldTimestamp Field = "timestamp"
	FieldLevel     Field = "level"
	FieldService   Field = "service"
	FieldMessage   Field = "message"
)

// Options controls how entries are sorted.
type Options struct {
	By        Field
	Descending bool
}

// DefaultOptions returns the default sort options (ascending by timestamp).
func DefaultOptions() Options {
	return Options{
		By:        FieldTimestamp,
		Descending: false,
	}
}

// Apply returns a new snapshot with entries sorted according to opts.
// The original snapshot is not modified.
func Apply(snap *snapshot.Snapshot, opts Options) (*snapshot.Snapshot, error) {
	if snap == nil {
		return nil, nil
	}

	out := &snapshot.Snapshot{
		ID:        snap.ID,
		Label:     snap.Label,
		CreatedAt: snap.CreatedAt,
		Source:    snap.Source,
		Tags:      snap.Tags,
	}

	entries := make([]snapshot.Entry, len(snap.Entries))
	copy(entries, snap.Entries)

	gosort.SliceStable(entries, func(i, j int) bool {
		less := compareEntries(entries[i], entries[j], opts.By)
		if opts.Descending {
			return !less
		}
		return less
	})

	out.Entries = entries
	return out, nil
}

func compareEntries(a, b snapshot.Entry, by Field) bool {
	switch by {
	case FieldTimestamp:
		return a.Timestamp.Before(b.Timestamp)
	case FieldLevel:
		return strings.ToLower(a.Level) < strings.ToLower(b.Level)
	case FieldService:
		return strings.ToLower(a.ServiceID) < strings.ToLower(b.ServiceID)
	case FieldMessage:
		return a.Message < b.Message
	default:
		return a.Timestamp.Before(b.Timestamp)
	}
}
