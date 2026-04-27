// Package rename provides functionality for renaming (relabeling) snapshot metadata.
package rename

import (
	"errors"
	"strings"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
)

// Options configures the rename operation.
type Options struct {
	// Label sets a new human-readable label for the snapshot.
	Label string
	// Source sets a new source identifier for the snapshot.
	Source string
	// Tags replaces or merges additional metadata tags.
	Tags map[string]string
	// MergeTags controls whether Tags are merged with existing tags (true)
	// or replace them entirely (false).
	MergeTags bool
}

// Apply updates the metadata of snap according to opts and returns the
// modified snapshot. The original snapshot is not mutated.
func Apply(snap *snapshot.Snapshot, opts Options) (*snapshot.Snapshot, error) {
	if snap == nil {
		return nil, errors.New("rename: snapshot must not be nil")
	}

	label := strings.TrimSpace(opts.Label)
	source := strings.TrimSpace(opts.Source)

	if label == "" && source == "" && len(opts.Tags) == 0 {
		return nil, errors.New("rename: at least one of label, source, or tags must be provided")
	}

	out := snapshot.New(snap.Meta.Source, snap.Meta.Label)
	out.Meta.CreatedAt = snap.Meta.CreatedAt

	if label != "" {
		out.Meta.Label = label
	}
	if source != "" {
		out.Meta.Source = source
	}

	// Copy existing tags then optionally merge/replace.
	merged := make(map[string]string)
	if opts.MergeTags {
		for k, v := range snap.Meta.Tags {
			merged[k] = v
		}
	}
	for k, v := range opts.Tags {
		merged[k] = v
	}
	if len(merged) > 0 {
		out.Meta.Tags = merged
	} else if !opts.MergeTags {
		out.Meta.Tags = snap.Meta.Tags
	}

	out.Meta.UpdatedAt = time.Now().UTC()

	for _, e := range snap.Entries {
		out.AddEntry(e)
	}

	return out, nil
}
