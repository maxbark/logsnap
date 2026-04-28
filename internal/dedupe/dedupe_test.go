package dedupe_test

import (
	"testing"
	"time"

	"logsnap/internal/dedupe"
	"logsnap/internal/snapshot"
)

func makeSnap(label string, entries []snapshot.Entry) *snapshot.Snapshot {
	s := snapshot.New(label, "test-source")
	s.Entries = entries
	return s
}

func entry(msg, level, svc string) snapshot.Entry {
	return snapshot.Entry{
		Message:   msg,
		Level:     level,
		ServiceID: svc,
		Timestamp: time.Now(),
	}
}

func TestApply_DedupeByMessage(t *testing.T) {
	snap := makeSnap("test", []snapshot.Entry{
		entry("hello", "info", "svc-a"),
		entry("hello", "warn", "svc-b"),
		entry("world", "info", "svc-a"),
	})
	out, err := dedupe.Apply(snap, dedupe.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestApply_KeepLast(t *testing.T) {
	snap := makeSnap("test", []snapshot.Entry{
		entry("dup", "info", "svc-a"),
		entry("dup", "warn", "svc-b"),
	})
	opts := dedupe.Options{Fields: []string{"message"}, KeepFirst: false}
	out, err := dedupe.Apply(snap, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out.Entries))
	}
	if out.Entries[0].Level != "warn" {
		t.Errorf("expected last entry (warn), got %q", out.Entries[0].Level)
	}
}

func TestApply_MultiFieldKey(t *testing.T) {
	snap := makeSnap("test", []snapshot.Entry{
		entry("msg", "info", "svc-a"),
		entry("msg", "info", "svc-b"), // different service — not a dup
		entry("msg", "info", "svc-a"), // true dup
	})
	opts := dedupe.Options{Fields: []string{"message", "service_id"}, KeepFirst: true}
	out, err := dedupe.Apply(snap, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestApply_NilSnapshot(t *testing.T) {
	_, err := dedupe.Apply(nil, dedupe.DefaultOptions())
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}

func TestApply_UnsupportedField(t *testing.T) {
	snap := makeSnap("test", nil)
	_, err := dedupe.Apply(snap, dedupe.Options{Fields: []string{"unknown_field"}})
	if err == nil {
		t.Error("expected error for unsupported field")
	}
}

func TestApply_PreservesMetadata(t *testing.T) {
	snap := makeSnap("my-label", nil)
	snap.Tags = map[string]string{"env": "prod"}
	out, err := dedupe.Apply(snap, dedupe.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Label != "my-label" {
		t.Errorf("expected label %q, got %q", "my-label", out.Label)
	}
	if out.Tags["env"] != "prod" {
		t.Errorf("expected tag env=prod, got %q", out.Tags["env"])
	}
}
