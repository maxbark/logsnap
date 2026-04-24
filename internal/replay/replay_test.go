package replay_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/replay"
	"github.com/yourorg/logsnap/internal/snapshot"
)

func makeReplaySnap() *snapshot.Snapshot {
	snap := snapshot.New("test-service", "v1.0.0")
	snap.AddEntry(snapshot.Entry{
		Level:     "info",
		ServiceID: "svc-a",
		Message:   "started",
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	snap.AddEntry(snapshot.Entry{
		Level:     "error",
		ServiceID: "svc-b",
		Message:   "connection refused",
		Timestamp: time.Date(2024, 1, 1, 0, 0, 1, 0, time.UTC),
	})
	return snap
}

func TestRun_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := replay.DefaultOptions()
	if err := replay.Run(&buf, makeReplaySnap(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[info]") || !strings.Contains(out, "started") {
		t.Errorf("expected info entry in output, got: %s", out)
	}
	if !strings.Contains(out, "[error]") || !strings.Contains(out, "connection refused") {
		t.Errorf("expected error entry in output, got: %s", out)
	}
}

func TestRun_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := replay.Options{Format: "json"}
	if err := replay.Run(&buf, makeReplaySnap(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"level":"info"`) {
		t.Errorf("expected JSON level field, got: %s", out)
	}
}

func TestRun_LogfmtFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := replay.Options{Format: "logfmt"}
	if err := replay.Run(&buf, makeReplaySnap(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "level=info") {
		t.Errorf("expected logfmt level field, got: %s", out)
	}
}

func TestRun_NilSnapshot(t *testing.T) {
	var buf bytes.Buffer
	err := replay.Run(&buf, nil, replay.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestRun_EmptySnapshot(t *testing.T) {
	var buf bytes.Buffer
	snap := snapshot.New("svc", "v0")
	if err := replay.Run(&buf, snap, replay.DefaultOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for empty snapshot")
	}
}
