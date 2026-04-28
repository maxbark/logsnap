package sort_test

import (
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
	logsort "github.com/yourorg/logsnap/internal/sort"
)

func makeSnap(entries []snapshot.Entry) *snapshot.Snapshot {
	s := snapshot.New("test-snap", "test-source")
	s.Entries = entries
	return s
}

var (
	t1 = time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 = time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	t3 = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
)

func TestApply_SortByTimestampAsc(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{
		{Timestamp: t3, Message: "c"},
		{Timestamp: t1, Message: "a"},
		{Timestamp: t2, Message: "b"},
	})
	out, err := logsort.Apply(snap, logsort.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Message != "a" || out.Entries[2].Message != "c" {
		t.Errorf("expected ascending timestamp order, got %v", out.Entries)
	}
}

func TestApply_SortByTimestampDesc(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{
		{Timestamp: t1, Message: "a"},
		{Timestamp: t3, Message: "c"},
		{Timestamp: t2, Message: "b"},
	})
	out, err := logsort.Apply(snap, logsort.Options{By: logsort.FieldTimestamp, Descending: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Message != "c" || out.Entries[2].Message != "a" {
		t.Errorf("expected descending timestamp order, got %v", out.Entries)
	}
}

func TestApply_SortByLevel(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{
		{Level: "warn", Message: "w"},
		{Level: "error", Message: "e"},
		{Level: "info", Message: "i"},
	})
	out, _ := logsort.Apply(snap, logsort.Options{By: logsort.FieldLevel})
	if out.Entries[0].Level != "error" || out.Entries[2].Level != "warn" {
		t.Errorf("expected level sort order error<info<warn, got %v", out.Entries)
	}
}

func TestApply_SortByMessage(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{
		{Message: "zebra"},
		{Message: "apple"},
		{Message: "mango"},
	})
	out, _ := logsort.Apply(snap, logsort.Options{By: logsort.FieldMessage})
	if out.Entries[0].Message != "apple" || out.Entries[2].Message != "zebra" {
		t.Errorf("expected alphabetical message order, got %v", out.Entries)
	}
}

func TestApply_NilSnapshot(t *testing.T) {
	out, err := logsort.Apply(nil, logsort.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Errorf("expected nil output for nil input")
	}
}

func TestApply_PreservesMetadata(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{{Timestamp: t1, Message: "x"}})
	snap.Label = "my-label"
	out, _ := logsort.Apply(snap, logsort.DefaultOptions())
	if out.Label != "my-label" {
		t.Errorf("expected label to be preserved, got %q", out.Label)
	}
}

func TestApply_EmptyEntries(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{})
	out, err := logsort.Apply(snap, logsort.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error for empty entries: %v", err)
	}
	if len(out.Entries) != 0 {
		t.Errorf("expected empty entries, got %d entries", len(out.Entries))
	}
}
