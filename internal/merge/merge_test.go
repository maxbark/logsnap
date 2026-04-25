package merge_test

import (
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/merge"
	"github.com/yourorg/logsnap/internal/snapshot"
)

func makeSnap(label string, entries []snapshot.Entry) *snapshot.Snapshot {
	s := snapshot.New(label, time.Now())
	for _, e := range entries {
		s.AddEntry(e)
	}
	return s
}

func TestMerge_CombinesEntries(t *testing.T) {
	a := makeSnap("a", []snapshot.Entry{
		{Level: "info", Message: "started", ServiceID: "svc-a"},
	})
	b := makeSnap("b", []snapshot.Entry{
		{Level: "error", Message: "failed", ServiceID: "svc-b"},
	})

	out, err := merge.Merge(merge.DefaultOptions(), a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestMerge_DeduplicateByMessage(t *testing.T) {
	entry := snapshot.Entry{Level: "warn", Message: "disk full", ServiceID: "svc-a"}
	a := makeSnap("a", []snapshot.Entry{entry})
	b := makeSnap("b", []snapshot.Entry{entry, {Level: "info", Message: "ok", ServiceID: "svc-b"}})

	opts := merge.Options{DeduplicateByMessage: true, Label: "deduped"}
	out, err := merge.Merge(opts, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Errorf("expected 2 entries after dedup, got %d", len(out.Entries))
	}
}

func TestMerge_SetsLabel(t *testing.T) {
	a := makeSnap("a", []snapshot.Entry{{Level: "info", Message: "x"}})
	b := makeSnap("b", []snapshot.Entry{{Level: "info", Message: "y"}})

	opts := merge.Options{Label: "release-combined"}
	out, err := merge.Merge(opts, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Label != "release-combined" {
		t.Errorf("expected label 'release-combined', got %q", out.Label)
	}
}

func TestMerge_RequiresAtLeastTwo(t *testing.T) {
	a := makeSnap("a", []snapshot.Entry{{Level: "info", Message: "x"}})
	_, err := merge.Merge(merge.DefaultOptions(), a)
	if err == nil {
		t.Error("expected error for fewer than 2 snapshots")
	}
}

func TestMerge_NilSnapshotReturnsError(t *testing.T) {
	a := makeSnap("a", []snapshot.Entry{{Level: "info", Message: "x"}})
	_, err := merge.Merge(merge.DefaultOptions(), a, nil)
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}
