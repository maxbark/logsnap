package ingest_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/ingest"
	"github.com/yourorg/logsnap/internal/snapshot"
)

func newSnap() *snapshot.Snapshot {
	return snapshot.New("test-service", "v1.0.0")
}

func TestFromReader_JSON(t *testing.T) {
	input := `{"level":"info","msg":"started","time":"2024-01-01T00:00:00Z"}
{"level":"error","message":"failed","time":"2024-01-01T00:01:00Z"}`
	snap := newSnap()
	err := ingest.FromReader(strings.NewReader(input), snap, ingest.Options{
		ServiceID: "svc-a",
		Format:    "json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap.Entries))
	}
	if snap.Entries[0].Level != "info" {
		t.Errorf("expected level info, got %q", snap.Entries[0].Level)
	}
	if snap.Entries[1].Message != "failed" {
		t.Errorf("expected message 'failed', got %q", snap.Entries[1].Message)
	}
	if snap.Entries[0].ServiceID != "svc-a" {
		t.Errorf("expected serviceID svc-a, got %q", snap.Entries[0].ServiceID)
	}
}

func TestFromReader_Logfmt(t *testing.T) {
	input := `level=info msg="server started" time=2024-01-01T00:00:00Z
level=warn msg=retrying`
	snap := newSnap()
	err := ingest.FromReader(strings.NewReader(input), snap, ingest.Options{
		ServiceID: "svc-b",
		Format:    "logfmt",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap.Entries))
	}
	if snap.Entries[0].Message != "server started" {
		t.Errorf("unexpected message: %q", snap.Entries[0].Message)
	}
	expected := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if !snap.Entries[0].Timestamp.Equal(expected) {
		t.Errorf("unexpected timestamp: %v", snap.Entries[0].Timestamp)
	}
}

func TestFromReader_SkipsBlankLines(t *testing.T) {
	input := `{"level":"info","msg":"a"}

{"level":"info","msg":"b"}`
	snap := newSnap()
	_ = ingest.FromReader(strings.NewReader(input), snap, ingest.Options{Format: "json"})
	if len(snap.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(snap.Entries))
	}
}

func TestFromReader_InvalidJSON(t *testing.T) {
	input := `not-json`
	snap := newSnap()
	err := ingest.FromReader(strings.NewReader(input), snap, ingest.Options{Format: "json"})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestFromReader_UnsupportedFormat(t *testing.T) {
	snap := newSnap()
	err := ingest.FromReader(strings.NewReader("anything"), snap, ingest.Options{Format: "xml"})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
