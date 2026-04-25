// Package annotate provides functionality for adding human-readable
// annotations (notes) to snapshot entries.
package annotate

import (
	"errors"
	"fmt"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
)

// Options controls annotation behaviour.
type Options struct {
	// EntryIndex is the zero-based index of the entry to annotate.
	// Use -1 to annotate all entries.
	EntryIndex int

	// Note is the text to attach to the entry/entries.
	Note string

	// Author is an optional label stored alongside the note.
	Author string

	// Overwrite replaces an existing annotation instead of appending.
	Overwrite bool
}

// Apply attaches notes to entries in snap according to opts.
// It returns a new snapshot with the annotations applied.
func Apply(snap *snapshot.Snapshot, opts Options) (*snapshot.Snapshot, error) {
	if snap == nil {
		return nil, errors.New("annotate: snapshot must not be nil")
	}
	if opts.Note == "" {
		return nil, errors.New("annotate: note must not be empty")
	}

	entries := snap.Entries
	if opts.EntryIndex != -1 {
		if opts.EntryIndex < 0 || opts.EntryIndex >= len(entries) {
			return nil, fmt.Errorf("annotate: index %d out of range (snapshot has %d entries)",
				opts.EntryIndex, len(entries))
		}
	}

	out := snapshot.New(snap.Label)
	out.CreatedAt = snap.CreatedAt

	for i, e := range entries {
		if opts.EntryIndex == -1 || opts.EntryIndex == i {
			e = annotateEntry(e, opts)
		}
		out.AddEntry(e)
	}
	return out, nil
}

func annotateEntry(e snapshot.Entry, opts Options) snapshot.Entry {
	if e.Fields == nil {
		e.Fields = map[string]string{}
	}

	note := opts.Note
	if opts.Author != "" {
		note = fmt.Sprintf("%s: %s", opts.Author, note)
	}
	note = fmt.Sprintf("[%s] %s", time.Now().UTC().Format(time.RFC3339), note)

	const key = "_annotation"
	if existing, ok := e.Fields[key]; ok && !opts.Overwrite {
		e.Fields[key] = existing + " | " + note
	} else {
		e.Fields[key] = note
	}
	return e
}
