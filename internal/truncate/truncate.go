// Package truncate provides utilities for limiting snapshot size by entry count or time range.
package truncate

import (
	"fmt"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
)

// Options controls how truncation is applied to a snapshot.
type Options struct {
	// MaxEntries limits the snapshot to the most recent N entries. 0 means no limit.
	MaxEntries int
	// Since discards entries older than this time. Zero value means no lower bound.
	Since time.Time
	// Until discards entries newer than this time. Zero value means no upper bound.
	Until time.Time
}

// Apply returns a new snapshot containing only the entries that satisfy the
// truncation options. The snapshot metadata (ID, service, timestamp) is
// preserved from the source snapshot.
func Apply(snap *snapshot.Snapshot, opts Options) (*snapshot.Snapshot, error) {
	if snap == nil {
		return nil, fmt.Errorf("truncate: snapshot must not be nil")
	}

	if opts.MaxEntries < 0 {
		return nil, fmt.Errorf("truncate: MaxEntries must be non-negative, got %d", opts.MaxEntries)
	}

	out := snapshot.New(snap.ServiceID)
	out.ID = snap.ID
	out.CapturedAt = snap.CapturedAt

	filtered := make([]snapshot.Entry, 0, len(snap.Entries))
	for _, e := range snap.Entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if !opts.Until.IsZero() && e.Timestamp.After(opts.Until) {
			continue
		}
		filtered = append(filtered, e)
	}

	// Apply MaxEntries by keeping the tail (most recent entries).
	if opts.MaxEntries > 0 && len(filtered) > opts.MaxEntries {
		filtered = filtered[len(filtered)-opts.MaxEntries:]
	}

	for _, e := range filtered {
		out.AddEntry(e)
	}

	return out, nil
}
