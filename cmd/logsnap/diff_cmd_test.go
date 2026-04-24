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

func writeDiffSnap(t *testing.T, snap *snapshot.Snapshot) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write snapshot: %v", err)
	}
	return path
}

func TestDiffCmd_ShowsAdded(t *testing.T) {
	snapA := snapshot.New("deploy-1")
	snapB := snapshot.New("deploy-2")
	snapB.AddEntry(snapshot.Entry{ServiceID: "svc-a", Level: "info", Message: "started"})

	pathA := writeDiffSnap(t, snapA)
	pathB := writeDiffSnap(t, snapB)

	cmd := newDiffCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{pathA, pathB})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ADDED") {
		t.Errorf("expected ADDED in output, got:\n%s", out)
	}
	if !strings.Contains(out, "svc-a") {
		t.Errorf("expected svc-a in output, got:\n%s", out)
	}
}

func TestDiffCmd_ShowsRemoved(t *testing.T) {
	snapA := snapshot.New("deploy-1")
	snapA.AddEntry(snapshot.Entry{ServiceID: "svc-b", Level: "error", Message: "crashed"})
	snapB := snapshot.New("deploy-2")

	pathA := writeDiffSnap(t, snapA)
	pathB := writeDiffSnap(t, snapB)

	cmd := newDiffCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{pathA, pathB})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "REMOVED") {
		t.Errorf("expected REMOVED in output, got:\n%s", out)
	}
}

func TestDiffCmd_WritesOutputFile(t *testing.T) {
	snapA := snapshot.New("deploy-1")
	snapB := snapshot.New("deploy-2")
	snapB.AddEntry(snapshot.Entry{ServiceID: "svc-c", Level: "warn", Message: "retrying"})

	pathA := writeDiffSnap(t, snapA)
	pathB := writeDiffSnap(t, snapB)
	outPath := filepath.Join(t.TempDir(), "diff_out.txt")

	cmd := newDiffCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetArgs([]string{pathA, pathB, "--output", outPath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("expected output file to exist: %v", err)
	}
	if !strings.Contains(string(data), "ADDED") {
		t.Errorf("expected ADDED in file output, got:\n%s", string(data))
	}
}
