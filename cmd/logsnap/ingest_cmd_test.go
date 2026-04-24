package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/logsnap/internal/snapshot"
)

func TestIngestCmd_FromStdin(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "out.snap")

	input := `{"level":"info","msg":"hello"}
{"level":"error","msg":"oops"}`

	old := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	_, _ = w.WriteString(input)
	w.Close()
	defer func() { os.Stdin = old }()

	cmd := newIngestCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--output", outPath, "--service", "svc-x", "--version", "v2", "--format", "json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Ingested 2 entries") {
		t.Errorf("unexpected output: %q", buf.String())
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("snapshot not written: %v", err)
	}
	var snap snapshot.Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		t.Fatalf("invalid snapshot JSON: %v", err)
	}
	if len(snap.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(snap.Entries))
	}
}

func TestIngestCmd_FromFile(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "app.log")
	outPath := filepath.Join(tmpDir, "out.snap")

	lines := "level=info msg=started\nlevel=warn msg=slow\n"
	if err := os.WriteFile(logFile, []byte(lines), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := newIngestCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--file", logFile, "--output", outPath, "--format", "logfmt"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Ingested 2 entries") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestIngestCmd_InvalidFormat(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "out.snap")

	r, w, _ := os.Pipe()
	old := os.Stdin
	os.Stdin = r
	_, _ = w.WriteString("some line\n")
	w.Close()
	defer func() { os.Stdin = old }()

	cmd := newIngestCmd()
	cmd.SetArgs([]string{"--output", outPath, "--format", "xml"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
