// Package sample provides functionality for sampling log entries from a snapshot.
package sample

import (
	"fmt"
	"math/rand"

	"github.com/yourorg/logsnap/internal/snapshot"
)

// Options controls how sampling is performed.
type Options struct {
	// N is the maximum number of entries to return.
	N int
	// Seed is used to initialize the random source for reproducibility.
	// A zero value means a random seed is used.
	Seed int64
	// Deterministic, when true, uses the provided Seed for reproducibility.
	Deterministic bool
}

// DefaultOptions returns sensible defaults for sampling.
func DefaultOptions() Options {
	return Options{
		N: 10,
		Deterministic: false,
	}
}

// Apply returns a new snapshot containing at most N randomly sampled entries
// from the input snapshot. Order of selected entries is preserved.
func Apply(snap *snapshot.Snapshot, opts Options) (*snapshot.Snapshot, error) {
	if snap == nil {
		return nil, fmt.Errorf("sample: snapshot is nil")
	}
	if opts.N <= 0 {
		return nil, fmt.Errorf("sample: N must be greater than zero, got %d", opts.N)
	}

	entries := snap.Entries
	if opts.N >= len(entries) {
		// Return a copy of all entries.
		out := snapshot.New(snap.Label, snap.Source)
		for _, e := range entries {
			out.AddEntry(e)
		}
		return out, nil
	}

	var rng *rand.Rand
	if opts.Deterministic {
		//nolint:gosec
		rng = rand.New(rand.NewSource(opts.Seed))
	} else {
		//nolint:gosec
		rng = rand.New(rand.NewSource(rand.Int63()))
	}

	// Fisher-Yates partial shuffle to pick N indices.
	indices := make([]int, len(entries))
	for i := range indices {
		indices[i] = i
	}
	for i := 0; i < opts.N; i++ {
		j := i + rng.Intn(len(indices)-i)
		indices[i], indices[j] = indices[j], indices[i]
	}
	selected := indices[:opts.N]

	// Sort selected indices to preserve original order.
	for i := 1; i < len(selected); i++ {
		for j := i; j > 0 && selected[j] < selected[j-1]; j-- {
			selected[j], selected[j-1] = selected[j-1], selected[j]
		}
	}

	out := snapshot.New(snap.Label, snap.Source)
	for _, idx := range selected {
		out.AddEntry(entries[idx])
	}
	return out, nil
}
