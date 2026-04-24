package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/logsnap/internal/snapshot"
)

// Format represents the supported export formats.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
	FormatText Format = "text"
)

// Write serializes the snapshot entries to the given writer in the specified format.
func Write(snap *snapshot.Snapshot, format Format, w io.Writer) error {
	switch format {
	case FormatJSON:
		return writeJSON(snap, w)
	case FormatCSV:
		return writeCSV(snap, w)
	case FormatText:
		return writeText(snap, w)
	default:
		return fmt.Errorf("unsupported export format: %q", format)
	}
}

func writeJSON(snap *snapshot.Snapshot, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}

func writeCSV(snap *snapshot.Snapshot, w io.Writer) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"timestamp", "level", "service_id", "message"}); err != nil {
		return err
	}
	for _, e := range snap.Entries {
		row := []string{
			e.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			e.Level,
			e.ServiceID,
			e.Message,
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func writeText(snap *snapshot.Snapshot, w io.Writer) error {
	var sb strings.Builder
	for _, e := range snap.Entries {
		sb.WriteString(fmt.Sprintf("[%s] %s %s: %s\n",
			e.Level,
			e.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			e.ServiceID,
			e.Message,
		))
	}
	_, err := fmt.Fprint(w, sb.String())
	return err
}
