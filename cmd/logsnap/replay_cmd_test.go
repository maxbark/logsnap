package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
)

func writeReplaySnap(t *testing.T, dir string) string {
	t.Helper()
	snap := snapshot.New("replay-svc", "v2.0.0")
	snap.AddEntry(snapshot.Entry{
		Level:     "info",
		ServiceID: "replay-svc",
		Message:   "boot complete",
		Timestamp: time.Now(),
	})
	path := filepath.Join(dir, "replay.snap")
	data, _ := json.Marshal(snap)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("failed to write snap: %v", err)
	}
	return path
}

func TestReplayCmd_TextOutput(t *testing.T) {
	dir := t.TempDir()
	snapPath := writeReplaySnap(t, dir)

	cmd := newReplayCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{snapPath, "--format", "text"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "boot complete") {
		t.Errorf("expected 'boot complete' in output, got: %s", buf.String())
	}
}

func TestReplayCmd_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	snapPath := writeReplaySnap(t, dir)

	cmd := newReplayCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{snapPath, "--format", "json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"level":"info"`) {
		t.Errorf("expected JSON output, got: %s", buf.String())
	}
}

func TestReplayCmd_MissingSnapshot(t *testing.T) {
	cmd := newReplayCmd()
	cmd.SetArgs([]string{"/nonexistent/path.snap"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing snapshot file")
	}
}
