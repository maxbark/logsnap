// Package count provides entry counting grouped by a snapshot field.
package count

import (
	"fmt"
	"sort"

	"logsnap/internal/snapshot"
)

// Result holds the count output for a snapshot.
type Result struct {
	Field   string
	Counts  map[string]int
	Total   int
}

// ValidFields are the snapshot entry fields supported for grouping.
var ValidFields = []string{"level", "service_id", "message"}

// Apply counts entries in snap grouped by the given field.
// Supported fields: "level", "service_id", "message".
func Apply(snap *snapshot.Snapshot, field string) (*Result, error) {
	if snap == nil {
		return nil, fmt.Errorf("count: snapshot is nil")
	}

	switch field {
	case "level", "service_id", "message":
	default:
		return nil, fmt.Errorf("count: unsupported field %q; valid fields: %v", field, ValidFields)
	}

	counts := make(map[string]int)
	for _, e := range snap.Entries {
		var key string
		switch field {
		case "level":
			key = e.Level
		case "service_id":
			key = e.ServiceID
		case "message":
			key = e.Message
		}
		if key == "" {
			key = "(none)"
		}
		counts[key]++
	}

	return &Result{
		Field:  field,
		Counts: counts,
		Total:  len(snap.Entries),
	}, nil
}

// SortedKeys returns the keys of r.Counts in descending count order.
func (r *Result) SortedKeys() []string {
	keys := make([]string, 0, len(r.Counts))
	for k := range r.Counts {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if r.Counts[keys[i]] != r.Counts[keys[j]] {
			return r.Counts[keys[i]] > r.Counts[keys[j]]
		}
		return keys[i] < keys[j]
	})
	return keys
}
