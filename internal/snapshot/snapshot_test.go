package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logsnap/internal/snapshot"
)

func TestNew(t *testing.T) {
	s := snapshot.New("production-v1")
	if s == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if s.Deployment != "production-v1" {
		t.Errorf("expected deployment %q, got %q", "production-v1", s.Deployment)
	}
	if s.ID == "" {
		t.Error("expected non-empty snapshot ID")
	}
	if len(s.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(s.Entries))
	}
}

func TestAddEntry(t *testing.T) {
	s := snapshot.New("staging")
	s.AddEntry("info", "service started", map[string]string{"service": "api"})
	s.AddEntry("error", "connection refused", nil)

	if len(s.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(s.Entries))
	}
	if s.Entries[0].Level != "info" {
		t.Errorf("expected level %q, got %q", "info", s.Entries[0].Level)
	}
	if s.Entries[1].Message != "connection refused" {
		t.Errorf("expected message %q, got %q", "connection refused", s.Entries[1].Message)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := snapshot.New("canary")
	orig.AddEntry("warn", "high latency", map[string]string{"region": "us-east-1"})

	if err := orig.Save(path); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.ID != orig.ID {
		t.Errorf("ID mismatch: got %q, want %q", loaded.ID, orig.ID)
	}
	if loaded.Deployment != orig.Deployment {
		t.Errorf("Deployment mismatch: got %q, want %q", loaded.Deployment, orig.Deployment)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Fields["region"] != "us-east-1" {
		t.Errorf("field mismatch: got %q", loaded.Entries[0].Fields["region"])
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestSaveInvalidPath(t *testing.T) {
	s := snapshot.New("test")
	err := s.Save("/nonexistent_dir/snap.json")
	if err == nil {
		t.Error("expected error for invalid save path, got nil")
	}
	_ = os.Remove("/nonexistent_dir/snap.json")
}
