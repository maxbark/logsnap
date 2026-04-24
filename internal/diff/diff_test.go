package diff_test

import (
	"strings"
	"testing"

	"github.com/logsnap/internal/diff"
	"github.com/logsnap/internal/snapshot"
)

func makeSnapshot(entries []snapshot.Entry) *snapshot.Snapshot {
	s := snapshot.New("test")
	for _, e := range entries {
		s.AddEntry(e)
	}
	return s
}

func TestCompare_Added(t *testing.T) {
	base := makeSnapshot(nil)
	current := makeSnapshot([]snapshot.Entry{
		{Level: "INFO", Message: "service started", Source: "main"},
	})

	result := diff.Compare(base, current)

	if len(result.Added) != 1 {
		t.Fatalf("expected 1 added entry, got %d", len(result.Added))
	}
	if result.Added[0].Message != "service started" {
		t.Errorf("unexpected added message: %s", result.Added[0].Message)
	}
}

func TestCompare_Removed(t *testing.T) {
	base := makeSnapshot([]snapshot.Entry{
		{Level: "WARN", Message: "disk usage high", Source: "monitor"},
	})
	current := makeSnapshot(nil)

	result := diff.Compare(base, current)

	if len(result.Removed) != 1 {
		t.Fatalf("expected 1 removed entry, got %d", len(result.Removed))
	}
}

func TestCompare_Changed(t *testing.T) {
	base := makeSnapshot([]snapshot.Entry{
		{Level: "INFO", Message: "old message", Source: "api"},
	})
	current := makeSnapshot([]snapshot.Entry{
		{Level: "INFO", Message: "new message", Source: "api"},
	})

	result := diff.Compare(base, current)

	if len(result.Changed) != 1 {
		t.Fatalf("expected 1 changed entry, got %d", len(result.Changed))
	}
	if result.Changed[0].From.Message != "old message" {
		t.Errorf("unexpected from message: %s", result.Changed[0].From.Message)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	entries := []snapshot.Entry{
		{Level: "ERROR", Message: "connection refused", Source: "db"},
	}
	base := makeSnapshot(entries)
	current := makeSnapshot(entries)

	result := diff.Compare(base, current)

	if len(result.Added)+len(result.Removed)+len(result.Changed) != 0 {
		t.Error("expected no differences between identical snapshots")
	}
}

func TestResult_Summary(t *testing.T) {
	base := makeSnapshot([]snapshot.Entry{
		{Level: "INFO", Message: "old", Source: "svc"},
	})
	current := makeSnapshot([]snapshot.Entry{
		{Level: "INFO", Message: "new", Source: "svc"},
		{Level: "WARN", Message: "extra", Source: "svc2"},
	})

	result := diff.Compare(base, current)
	summary := result.Summary()

	if !strings.Contains(summary, "+1 added") {
		t.Errorf("summary missing added count: %s", summary)
	}
	if !strings.Contains(summary, "~1 changed") {
		t.Errorf("summary missing changed count: %s", summary)
	}
}
