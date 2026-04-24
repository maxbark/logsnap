package summary_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
	"github.com/yourorg/logsnap/internal/summary"
)

func makeSnap() *snapshot.Snapshot {
	snap := snapshot.New("test-snap", "v1")
	now := time.Now()
	snap.Entries = []snapshot.Entry{
		{Timestamp: now, Level: "info", ServiceID: "svc-a", Message: "started"},
		{Timestamp: now, Level: "error", ServiceID: "svc-a", Message: "failed"},
		{Timestamp: now, Level: "info", ServiceID: "svc-b", Message: "started"},
		{Timestamp: now, Level: "warn", ServiceID: "svc-b", Message: "slow query"},
		{Timestamp: now, Level: "info", ServiceID: "svc-a", Message: "started"}, // duplicate message
	}
	return snap
}

func TestCompute_TotalEntries(t *testing.T) {
	stats, err := summary.Compute(makeSnap())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.TotalEntries != 5 {
		t.Errorf("expected 5 entries, got %d", stats.TotalEntries)
	}
}

func TestCompute_ByLevel(t *testing.T) {
	stats, _ := summary.Compute(makeSnap())
	if stats.ByLevel["info"] != 3 {
		t.Errorf("expected 3 info entries, got %d", stats.ByLevel["info"])
	}
	if stats.ByLevel["error"] != 1 {
		t.Errorf("expected 1 error entry, got %d", stats.ByLevel["error"])
	}
	if stats.ByLevel["warn"] != 1 {
		t.Errorf("expected 1 warn entry, got %d", stats.ByLevel["warn"])
	}
}

func TestCompute_ByService(t *testing.T) {
	stats, _ := summary.Compute(makeSnap())
	if stats.ByService["svc-a"] != 3 {
		t.Errorf("expected 3 svc-a entries, got %d", stats.ByService["svc-a"])
	}
	if stats.ByService["svc-b"] != 2 {
		t.Errorf("expected 2 svc-b entries, got %d", stats.ByService["svc-b"])
	}
}

func TestCompute_UniqueMessages(t *testing.T) {
	stats, _ := summary.Compute(makeSnap())
	if stats.UniqueMessages != 3 {
		t.Errorf("expected 3 unique messages, got %d", stats.UniqueMessages)
	}
}

func TestCompute_NilSnapshot(t *testing.T) {
	_, err := summary.Compute(nil)
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}

func TestPrint_ContainsExpectedSections(t *testing.T) {
	stats, _ := summary.Compute(makeSnap())
	var buf bytes.Buffer
	summary.Print(&buf, stats)
	out := buf.String()
	for _, want := range []string{"Total entries", "By level", "By service", "svc-a", "info"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}
