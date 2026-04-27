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

func writeSampleSnap(t *testing.T, n int) string {
	t.Helper()
	snap := snapshot.New("sample-test", "src")
	for i := 0; i < n; i++ {
		snap.AddEntry(snapshot.Entry{
			Timestamp: time.Now(),
			Level:     "info",
			Message:   fmt.Sprintf("log line %d", i),
			ServiceID: "svc",
		})
	}
	f, err := os.CreateTemp(t.TempDir(), "snap-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(snap); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

import "fmt"

func TestSampleCmd_PrintsToStdout(t *testing.T) {
	snapPath := writeSampleSnap(t, 20)
	root := newRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"sample", "--count", "5", "--deterministic", "--seed", "42", snapPath})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
	if len(lines) != 5 {
		t.Errorf("expected 5 output lines, got %d", len(lines))
	}
}

func TestSampleCmd_SavesOutput(t *testing.T) {
	snapPath := writeSampleSnap(t, 30)
	outPath := filepath.Join(t.TempDir(), "sampled.json")
	root := newRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"sample", "--count", "8", "--output", outPath, snapPath})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(outPath); err != nil {
		t.Errorf("expected output file to exist: %v", err)
	}
}

func TestSampleCmd_MissingSnapshot(t *testing.T) {
	root := newRootCmd()
	root.SetArgs([]string{"sample", "--count", "5", "/no/such/file.json"})
	if err := root.Execute(); err == nil {
		t.Error("expected error for missing snapshot")
	}
}
