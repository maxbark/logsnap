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
