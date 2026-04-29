// Package watch provides live-tail functionality for snapshot files,
// re-emitting new log entries as they are appended.
package watch

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
)

// Options configures the Watch run.
type Options struct {
	// PollInterval is how often the file is checked for new entries.
	PollInterval time.Duration
	// Format controls how entries are printed: "text", "json", or "logfmt".
	Format string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		PollInterval: 500 * time.Millisecond,
		Format:       "text",
	}
}

// Run watches a snapshot file and writes new entries to out until ctx is cancelled.
func Run(ctx context.Context, path string, out *os.File, opts Options) error {
	if opts.PollInterval == 0 {
		opts.PollInterval = DefaultOptions().PollInterval
	}
	if opts.Format == "" {
		opts.Format = DefaultOptions().Format
	}

	seen := 0
	ticker := time.NewTicker(opts.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			snap, err := snapshot.Load(path)
			if err != nil {
				continue // file may not exist yet
			}
			entries := snap.Entries
			if len(entries) <= seen {
				continue
			}
			for _, e := range entries[seen:] {
				switch opts.Format {
				case "json":
					b, _ := json.Marshal(e)
					_, _ = out.Write(append(b, '\n'))
				case "logfmt":
					_, _ = out.WriteString(formatLogfmt(e) + "\n")
				default:
					_, _ = out.WriteString(formatText(e) + "\n")
				}
			}
			seen = len(entries)
		}
	}
}

func formatText(e snapshot.Entry) string {
	return e.Timestamp.Format(time.RFC3339) + " [" + e.Level + "] " + e.ServiceID + ": " + e.Message
}

func formatLogfmt(e snapshot.Entry) string {
	return "ts=" + e.Timestamp.Format(time.RFC3339) +
		" level=" + e.Level +
		" service=" + e.ServiceID +
		" msg=\"" + e.Message + "\""
}
