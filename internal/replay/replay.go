// Package replay provides functionality for replaying log snapshots
// to stdout or a writer, with optional delay simulation between entries.
package replay

import (
	"fmt"
	"io"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
)

// Options configures replay behaviour.
type Options struct {
	// Delay between each log entry (0 means no delay).
	Delay time.Duration
	// Format controls how each entry is printed: "text", "json", or "logfmt".
	Format string
}

// DefaultOptions returns sensible defaults for replay.
func DefaultOptions() Options {
	return Options{
		Delay:  0,
		Format: "text",
	}
}

// Run replays all entries in snap to w, respecting opts.
func Run(w io.Writer, snap *snapshot.Snapshot, opts Options) error {
	if snap == nil {
		return fmt.Errorf("replay: snapshot must not be nil")
	}

	format := opts.Format
	if format == "" {
		format = "text"
	}

	for i, entry := range snap.Entries {
		var line string
		switch format {
		case "json":
			line = formatJSON(entry)
		case "logfmt":
			line = formatLogfmt(entry)
		default:
			line = formatText(entry)
		}

		if _, err := fmt.Fprintln(w, line); err != nil {
			return fmt.Errorf("replay: write error at entry %d: %w", i, err)
		}

		if opts.Delay > 0 && i < len(snap.Entries)-1 {
			time.Sleep(opts.Delay)
		}
	}
	return nil
}

// RunFiltered replays only the entries in snap that satisfy the given predicate.
// It otherwise behaves identically to Run.
func RunFiltered(w io.Writer, snap *snapshot.Snapshot, opts Options, keep func(snapshot.Entry) bool) error {
	if snap == nil {
		return fmt.Errorf("replay: snapshot must not be nil")
	}
	if keep == nil {
		return Run(w, snap, opts)
	}

	filtered := &snapshot.Snapshot{}
	for _, entry := range snap.Entries {
		if keep(entry) {
			filtered.Entries = append(filtered.Entries, entry)
		}
	}
	return Run(w, filtered, opts)
}

func formatText(e snapshot.Entry) string {
	return fmt.Sprintf("[%s] %s %s: %s", e.Level, e.Timestamp.Format(time.RFC3339), e.ServiceID, e.Message)
}

func formatJSON(e snapshot.Entry) string {
	return fmt.Sprintf(`{"level":%q,"service":%q,"msg":%q,"ts":%q}`,
		e.Level, e.ServiceID, e.Message, e.Timestamp.Format(time.RFC3339))
}

func formatLogfmt(e snapshot.Entry) string {
	return fmt.Sprintf("level=%s service=%s msg=%q ts=%s",
		e.Level, e.ServiceID, e.Message, e.Timestamp.Format(time.RFC3339))
}
