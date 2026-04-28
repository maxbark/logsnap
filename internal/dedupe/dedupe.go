// Package dedupe provides functionality for removing duplicate log entries
// from a snapshot based on configurable fields.
package dedupe

import (
	"fmt"
	"strings"

	"logsnap/internal/snapshot"
)

// Options controls how deduplication is performed.
type Options struct {
	// Fields to use as the deduplication key. Supported: "message", "level", "service_id".
	// Defaults to ["message"] if empty.
	Fields []string
	// KeepFirst retains the first occurrence when true; otherwise keeps the last.
	KeepFirst bool
}

// DefaultOptions returns sensible defaults for deduplication.
func DefaultOptions() Options {
	return Options{
		Fields:    []string{"message"},
		KeepFirst: true,
	}
}

// Apply removes duplicate entries from snap according to opts.
// Returns a new snapshot with duplicates removed and metadata preserved.
func Apply(snap *snapshot.Snapshot, opts Options) (*snapshot.Snapshot, error) {
	if snap == nil {
		return nil, fmt.Errorf("dedupe: snapshot is nil")
	}
	if len(opts.Fields) == 0 {
		opts.Fields = DefaultOptions().Fields
	}

	for _, f := range opts.Fields {
		switch f {
		case "message", "level", "service_id":
		default:
			return nil, fmt.Errorf("dedupe: unsupported field %q", f)
		}
	}

	out := snapshot.New(snap.Label, snap.Source)
	out.CreatedAt = snap.CreatedAt
	out.Tags = snap.Tags

	seen := make(map[string]int) // key -> index in out.Entries

	for _, entry := range snap.Entries {
		key := buildKey(entry, opts.Fields)
		if idx, exists := seen[key]; exists {
			if !opts.KeepFirst {
				// Replace existing entry with the latest occurrence.
				out.Entries[idx] = entry
			}
			continue
		}
		seen[key] = len(out.Entries)
		out.Entries = append(out.Entries, entry)
	}

	return out, nil
}

func buildKey(entry snapshot.Entry, fields []string) string {
	parts := make([]string, 0, len(fields))
	for _, f := range fields {
		switch f {
		case "message":
			parts = append(parts, entry.Message)
		case "level":
			parts = append(parts, entry.Level)
		case "service_id":
			parts = append(parts, entry.ServiceID)
		}
	}
	return strings.Join(parts, "|")
}
