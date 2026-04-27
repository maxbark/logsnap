package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
)

func writeRedactSnap(t *testing.T, dir string) string {
	t.Helper()
	snap := snapshot.New("redact-test", "unit")
	snap.AddEntry(snapshot.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "user action",
		ServiceID: "auth",
		Fields:    map[string]string{"email": "bob@example.com", "role": "admin"},
	})
	path := filepath.Join(dir, "redact_input.snap")
	data, _ := json.Marshal(snap)
	os.WriteFile(path, data, 0644)
	return path
}

func TestRedactCmd_PrintsToStdout(t *testing.T) {
	dir := t.TempDir()
	snapPath := writeRedactSnap(t, dir)

	root := newRootCmd()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"redact", snapPath, "--field", "email"})

	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if out == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRedactCmd_SavesOutput(t *testing.T) {
	dir := t.TempDir()
	snapPath := writeRedactSnap(t, dir)
	outPath := filepath.Join(dir, "redacted.snap")

	root := newRootCmd()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"redact", snapPath, "--field", "email", "--output", outPath})

	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(outPath); err != nil {
		t.Errorf("expected output file to exist: %v", err)
	}
}

func TestRedactCmd_NoFieldsOrPatterns(t *testing.T) {
	dir := t.TempDir()
	snapPath := writeRedactSnap(t, dir)

	root := newRootCmd()
	root.SetArgs([]string{"redact", snapPath})

	if err := root.Execute(); err == nil {
		t.Error("expected error when no fields or patterns specified")
	}
}
