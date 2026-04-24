package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single structured log entry captured in a snapshot.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
}

// Snapshot holds a collection of log entries along with metadata.
type Snapshot struct {
	ID          string    `json:"id"`
	CapturedAt  time.Time `json:"captured_at"`
	Deployment  string    `json:"deployment"`
	Entries     []Entry   `json:"entries"`
}

// New creates a new Snapshot with the given deployment label.
func New(deployment string) *Snapshot {
	return &Snapshot{
		ID:         fmt.Sprintf("%d", time.Now().UnixNano()),
		CapturedAt: time.Now().UTC(),
		Deployment: deployment,
		Entries:    []Entry{},
	}
}

// AddEntry appends a log entry to the snapshot.
func (s *Snapshot) AddEntry(level, message string, fields map[string]string) {
	s.Entries = append(s.Entries, Entry{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Message:   message,
		Fields:    fields,
	})
}

// Save writes the snapshot as JSON to the specified file path.
func (s *Snapshot) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: failed to marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshot: failed to write file %q: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: failed to read file %q: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: failed to unmarshal: %w", err)
	}
	return &s, nil
}
