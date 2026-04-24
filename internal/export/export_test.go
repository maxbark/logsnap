package export_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/export"
	"github.com/yourorg/logsnap/internal/snapshot"
)

var fixedTime = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

func makeExportSnap() *snapshot.Snapshot {
	snap := snapshot.New("export-test")
	snap.AddEntry(snapshot.Entry{
		Timestamp: fixedTime,
		Level:     "INFO",
		ServiceID: "svc-a",
		Message:   "started",
	})
	snap.AddEntry(snapshot.Entry{
		Timestamp: fixedTime,
		Level:     "ERROR",
		ServiceID: "svc-b",
		Message:   "failed to connect",
	})
	return snap
}

func TestWriteJSON(t *testing.T) {
	snap := makeExportSnap()
	var buf bytes.Buffer
	if err := export.Write(snap, export.FormatJSON, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}

func TestWriteCSV(t *testing.T) {
	snap := makeExportSnap()
	var buf bytes.Buffer
	if err := export.Write(snap, export.FormatCSV, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (header + 2 entries), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "timestamp") {
		t.Errorf("expected CSV header, got: %s", lines[0])
	}
}

func TestWriteText(t *testing.T) {
	snap := makeExportSnap()
	var buf bytes.Buffer
	if err := export.Write(snap, export.FormatText, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[INFO]") || !strings.Contains(out, "[ERROR]") {
		t.Errorf("expected level tags in text output, got: %s", out)
	}
}

func TestWriteUnsupportedFormat(t *testing.T) {
	snap := makeExportSnap()
	var buf bytes.Buffer
	err := export.Write(snap, export.Format("xml"), &buf)
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}
