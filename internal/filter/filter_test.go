package filter_test

import (
	"testing"

	"github.com/logsnap/internal/filter"
	"github.com/logsnap/internal/snapshot"
)

func makeSnap() *snapshot.Snapshot {
	s := snapshot.New("snap-1", "svc-a")
	s.AddEntry(snapshot.Entry{ID: "1", Level: "info", ServiceID: "svc-a", Message: "started", Fields: map[string]string{"env": "prod"}})
	s.AddEntry(snapshot.Entry{ID: "2", Level: "error", ServiceID: "svc-a", Message: "connection failed", Fields: map[string]string{"env": "prod"}})
	s.AddEntry(snapshot.Entry{ID: "3", Level: "warn", ServiceID: "svc-b", Message: "retrying", Fields: map[string]string{"env": "staging"}})
	s.AddEntry(snapshot.Entry{ID: "4", Level: "info", ServiceID: "svc-b", Message: "done", Fields: map[string]string{"env": "staging"}})
	return s
}

func TestFilterByLevel(t *testing.T) {
	out, err := filter.Apply(makeSnap(), filter.Options{Level: "error"})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Entries) != 1 || out.Entries[0].ID != "2" {
		t.Fatalf("expected 1 error entry, got %d", len(out.Entries))
	}
}

func TestFilterByServiceID(t *testing.T) {
	out, err := filter.Apply(makeSnap(), filter.Options{ServiceID: "svc-b"})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestFilterByMessageRegex(t *testing.T) {
	out, err := filter.Apply(makeSnap(), filter.Options{MessageRe: "^conn"})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out.Entries))
	}
}

func TestFilterInvalidRegex(t *testing.T) {
	_, err := filter.Apply(makeSnap(), filter.Options{MessageRe: "["})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestFilterByFieldKey(t *testing.T) {
	out, err := filter.Apply(makeSnap(), filter.Options{FieldKey: "env", FieldValue: "staging"})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 staging entries, got %d", len(out.Entries))
	}
}

func TestFilterNoMatch(t *testing.T) {
	out, err := filter.Apply(makeSnap(), filter.Options{Level: "debug"})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(out.Entries))
	}
}
