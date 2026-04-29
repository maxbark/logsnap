// Package pivot provides functionality for grouping snapshot entries
// by a specified field and producing a summary table.
package pivot

import (
	"fmt"
	"sort"

	"github.com/user/logsnap/internal/snapshot"
)

// Result holds the pivoted data grouped by a field value.
type Result struct {
	// Field is the field name used to pivot.
	Field string
	// Groups maps each unique field value to its entries.
	Groups map[string][]*snapshot.Entry
	// Keys holds the sorted group keys for deterministic output.
	Keys []string
}

// Apply groups the entries in snap by the given field name.
// Supported fields: "level", "service_id", "message".
// Entries missing the field value are grouped under "(none)".
func Apply(snap *snapshot.Snapshot, field string) (*Result, error) {
	if snap == nil {
		return nil, fmt.Errorf("pivot: snapshot is nil")
	}
	if field == "" {
		return nil, fmt.Errorf("pivot: field must not be empty")
	}

	groups := make(map[string][]*snapshot.Entry)

	for i := range snap.Entries {
		e := &snap.Entries[i]
		key := extractField(e, field)
		groups[key] = append(groups[key], e)
	}

	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return &Result{
		Field:  field,
		Groups: groups,
		Keys:   keys,
	}, nil
}

func extractField(e *snapshot.Entry, field string) string {
	switch field {
	case "level":
		if e.Level == "" {
			return "(none)"
		}
		return e.Level
	case "service_id":
		if e.ServiceID == "" {
			return "(none)"
		}
		return e.ServiceID
	case "message":
		if e.Message == "" {
			return "(none)"
		}
		return e.Message
	default:
		if v, ok := e.Tags[field]; ok {
			return v
		}
		return "(none)"
	}
}
