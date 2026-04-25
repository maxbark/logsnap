package filter_test

import (
	"testing"

	"github.com/logsnap/internal/filter"
)

// TestFilterCombined verifies that multiple criteria are AND-ed together.
func TestFilterCombined(t *testing.T) {
	s := makeSnap()
	out, err := filter.Apply(s, filter.Options{
		Level:      "info",
		ServiceID:  "svc-b",
		FieldKey:   "env",
		FieldValue: "staging",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Entries) != 1 {
		t.Fatalf("expected 1 combined-filter entry, got %d", len(out.Entries))
	}
	if out.Entries[0].ID != "4" {
		t.Fatalf("expected entry ID 4, got %s", out.Entries[0].ID)
	}
}

// TestFilterPreservesSnapshotMeta verifies filtered snapshot carries derived ID.
func TestFilterPreservesSnapshotMeta(t *testing.T) {
	s := makeSnap()
	out, err := filter.Apply(s, filter.Options{Level: "info"})
	if err != nil {
		t.Fatal(err)
	}
	if out.Service != s.Service {
		t.Fatalf("service mismatch: want %s got %s", s.Service, out.Service)
	}
	if out.ID != s.ID+"-filtered" {
		t.Fatalf("unexpected filtered ID: %s", out.ID)
	}
}

// TestFilterNoMatchReturnsEmpty verifies that filters with no matching entries
// return an empty (non-nil) entries slice rather than an error.
func TestFilterNoMatchReturnsEmpty(t *testing.T) {
	s := makeSnap()
	out, err := filter.Apply(s, filter.Options{
		Level:     "debug",
		ServiceID: "svc-nonexistent",
	})
	if err != nil {
		t.Fatalf("expected no error for zero-match filter, got: %v", err)
	}
	if out.Entries == nil {
		t.Fatal("expected non-nil entries slice for zero-match filter")
	}
	if len(out.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(out.Entries))
	}
}
