// Package summary provides aggregation statistics over a snapshot's log entries.
package summary

import (
	"fmt"
	"io"
	"sort"

	"github.com/yourorg/logsnap/internal/snapshot"
)

// Stats holds aggregated counts derived from a snapshot.
type Stats struct {
	TotalEntries  int
	ByLevel       map[string]int
	ByService     map[string]int
	UniqueMessages int
}

// Compute derives Stats from the given snapshot.
func Compute(snap *snapshot.Snapshot) (*Stats, error) {
	if snap == nil {
		return nil, fmt.Errorf("summary: snapshot must not be nil")
	}

	stats := &Stats{
		ByLevel:   make(map[string]int),
		ByService: make(map[string]int),
	}

	messages := make(map[string]struct{})
	for _, e := range snap.Entries {
		stats.TotalEntries++
		if e.Level != "" {
			stats.ByLevel[e.Level]++
		}
		if e.ServiceID != "" {
			stats.ByService[e.ServiceID]++
		}
		if e.Message != "" {
			messages[e.Message] = struct{}{}
		}
	}
	stats.UniqueMessages = len(messages)
	return stats, nil
}

// Print writes a human-readable summary to w.
func Print(w io.Writer, stats *Stats) {
	fmt.Fprintf(w, "Total entries : %d\n", stats.TotalEntries)
	fmt.Fprintf(w, "Unique messages: %d\n", stats.UniqueMessages)

	if len(stats.ByLevel) > 0 {
		fmt.Fprintln(w, "By level:")
		for _, lvl := range sortedKeys(stats.ByLevel) {
			fmt.Fprintf(w, "  %-10s %d\n", lvl, stats.ByLevel[lvl])
		}
	}

	if len(stats.ByService) > 0 {
		fmt.Fprintln(w, "By service:")
		for _, svc := range sortedKeys(stats.ByService) {
			fmt.Fprintf(w, "  %-20s %d\n", svc, stats.ByService[svc])
		}
	}
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
