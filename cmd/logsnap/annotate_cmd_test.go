package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
)

func writeAnnotateSnap(t *testing.T, label string) string {
	t.Helper()
	s := snapshot.New(label)
	s.AddEntry(snapshot.Entry{
		Timestamp: time.Now().UTC(),
		Level:     "warn",
		Message:   "disk usage high",
		ServiceID: "storage",
		Fields:    map[string]string{},
	})
	s.AddEntry(snapshot.Entry{
		Timestamp: time.Now().UTC(),
		Level:     "error",
		Message:   "connection refused",
		ServiceID: "api",
		Fields:    map[string]string{},
	})
	p := filepath.Join(t.TempDir(), "snap.json")
	if err := s.Save(p); err != nil {
		t.Fatalf("setup: %v", err)
	}
	return p
}

func TestAnnotateCmd_AnnotatesAllEntries(t *testing.T) {
	p := writeAnnotateSnap(t, "prod")

	root := newRootCmd()
	root.SetArgs([]string{"annotate", p, "--note", "reviewed"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(p)
	if !strings.Contains(string(data), "reviewed") {
		t.Error("expected annotation in saved snapshot")
	}
}

func TestAnnotateCmd_SingleIndex(t *testing.T) {
	p := writeAnnotateSnap(t, "prod")
	out := filepath.Join(filepath.Dir(p), "out.json")

	root := newRootCmd()
	root.SetArgs([]string{"annotate", p, "--note", "only first", "--index", "0", "--output", out})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	var s snapshot.Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !strings.Contains(s.Entries[0].Fields["_annotation"], "only first") {
		t.Error("entry 0 should be annotated")
	}
	if _, ok := s.Entries[1].Fields["_annotation"]; ok {
		t.Error("entry 1 should not be annotated")
	}
}

func TestAnnotateCmd_MissingNote(t *testing.T) {
	p := writeAnnotateSnap(t, "prod")
	root := newRootCmd()
	root.SetArgs([]string{"annotate", p})
	if err := root.Execute(); err == nil {
		t.Error("expected error when --note is missing")
	}
}
