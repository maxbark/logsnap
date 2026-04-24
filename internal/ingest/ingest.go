package ingest

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
)

// Options controls how log lines are parsed during ingestion.
type Options struct {
	ServiceID string
	Format    string // "json" or "logfmt"
}

// FromReader reads log lines from r and appends parsed entries into snap.
func FromReader(r io.Reader, snap *snapshot.Snapshot, opts Options) error {
	scanner := bufio.NewScanner(r)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		entry, err := parseLine(line, opts)
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNum, err)
		}
		snap.AddEntry(entry)
	}
	return scanner.Err()
}

func parseLine(line string, opts Options) (snapshot.Entry, error) {
	switch opts.Format {
	case "json":
		return parseJSON(line, opts.ServiceID)
	case "logfmt":
		return parseLogfmt(line, opts.ServiceID)
	default:
		return snapshot.Entry{}, fmt.Errorf("unsupported format: %q", opts.Format)
	}
}

func parseJSON(line, serviceID string) (snapshot.Entry, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return snapshot.Entry{}, fmt.Errorf("invalid JSON: %w", err)
	}
	entry := snapshot.Entry{
		ServiceID: serviceID,
		Timestamp: time.Now().UTC(),
	}
	if v, ok := raw["level"].(string); ok {
		entry.Level = v
	}
	if v, ok := raw["msg"].(string); ok {
		entry.Message = v
	} else if v, ok := raw["message"].(string); ok {
		entry.Message = v
	}
	if v, ok := raw["time"].(string); ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			entry.Timestamp = t
		}
	}
	return entry, nil
}

func parseLogfmt(line, serviceID string) (snapshot.Entry, error) {
	fields := map[string]string{}
	for _, pair := range strings.Fields(line) {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			fields[parts[0]] = strings.Trim(parts[1], `"`)
		}
	}
	entry := snapshot.Entry{
		ServiceID: serviceID,
		Level:     fields["level"],
		Message:   fields["msg"],
		Timestamp: time.Now().UTC(),
	}
	if v, ok := fields["time"]; ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			entry.Timestamp = t
		}
	}
	return entry, nil
}
