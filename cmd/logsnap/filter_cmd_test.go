package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/logsnap/internal/snapshot"
)

func writeTempSnap(t *testing.T) string {
	t.Helper()
	s := snapshot.New("snap-cmd", "svc-x")
	s.AddEntry(snapshot.Entry{ID: "a", Level: "info", ServiceID: "svc-x", Message: "hello", Fields: map[string]string{}})
	s.AddEntry(snapshot.Entry{ID: "b", Level: "error", ServiceID: "svc-x", Message: "boom", Fields: map[string]string{}})
	dir := t.TempDir()
	p := filepath.Join(dir, "snap.json")
	if err := s.Save(p); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestFilterCmd_PrintsToStdout(t *testing.T) {
	input := writeTempSnap(t)
	cmd := newFilterCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--input", input, "--level", "error"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "boom") {
		t.Fatalf("expected 'boom' in output, got: %s", buf.String())
	}
}

func TestFilterCmd_SavesOutput(t *testing.T) {
	input := writeTempSnap(t)
	output := filepath.Join(t.TempDir(), "out.json")
	cmd := newFilterCmd()
	cmd.SetArgs([]string{"--input", input, "--output", output, "--level", "info"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(output); err != nil {
		t.Fatalf("output file not created: %v", err)
	}
}
