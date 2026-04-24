package truncate_test

import (
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
	"github.com/yourorg/logsnap/internal/truncate"
)

func makeSnap() *snapshot.Snapshot {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	snap := snapshot.New("svc-test")
	for i := 0; i < 5; i++ {
		snap.AddEntry(snapshot.Entry{
			Timestamp: base.Add(time.Duration(i) * time.Minute),
			Level:     "info",
			Message:   "entry",
			ServiceID: "svc-test",
		})
	}
	return snap
}

func TestApply_MaxEntries(t *testing.T) {
	snap := makeSnap()
	out, err := truncate.Apply(snap, truncate.Options{MaxEntries: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out.Entries))
	}
	// Should keep the most recent 3 (indices 2,3,4).
	base := time.Date(2024, 1, 1, 12, 2, 0, 0, time.UTC)
	if !out.Entries[0].Timestamp.Equal(base) {
		t.Errorf("expected first kept entry at %v, got %v", base, out.Entries[0].Timestamp)
	}
}

func TestApply_Since(t *testing.T) {
	snap := makeSnap()
	since := time.Date(2024, 1, 1, 12, 2, 0, 0, time.UTC)
	out, err := truncate.Apply(snap, truncate.Options{Since: since})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out.Entries))
	}
}

func TestApply_Until(t *testing.T) {
	snap := makeSnap()
	until := time.Date(2024, 1, 1, 12, 2, 0, 0, time.UTC)
	out, err := truncate.Apply(snap, truncate.Options{Until: until})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out.Entries))
	}
}

func TestApply_PreservesMetadata(t *testing.T) {
	snap := makeSnap()
	out, err := truncate.Apply(snap, truncate.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ServiceID != snap.ServiceID {
		t.Errorf("ServiceID mismatch: got %q, want %q", out.ServiceID, snap.ServiceID)
	}
	if out.ID != snap.ID {
		t.Errorf("ID mismatch: got %q, want %q", out.ID, snap.ID)
	}
}

func TestApply_NilSnapshot(t *testing.T) {
	_, err := truncate.Apply(nil, truncate.Options{})
	if err == nil {
		t.Error("expected error for nil snapshot, got nil")
	}
}

func TestApply_NegativeMaxEntries(t *testing.T) {
	snap := makeSnap()
	_, err := truncate.Apply(snap, truncate.Options{MaxEntries: -1})
	if err == nil {
		t.Error("expected error for negative MaxEntries, got nil")
	}
}
