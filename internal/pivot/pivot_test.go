package pivot_test

import (
	"testing"
	"time"

	"github.com/user/logsnap/internal/pivot"
	"github.com/user/logsnap/internal/snapshot"
)

func makeSnap(entries []snapshot.Entry) *snapshot.Snapshot {
	snap := snapshot.New("test", "pivot-test")
	snap.Entries = entries
	return snap
}

func TestApply_ByLevel(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{
		{Level: "info", Message: "a", Timestamp: time.Now()},
		{Level: "error", Message: "b", Timestamp: time.Now()},
		{Level: "info", Message: "c", Timestamp: time.Now()},
	})

	res, err := pivot.Apply(snap, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Field != "level" {
		t.Errorf("expected field 'level', got %q", res.Field)
	}
	if len(res.Groups["info"]) != 2 {
		t.Errorf("expected 2 info entries, got %d", len(res.Groups["info"]))
	}
	if len(res.Groups["error"]) != 1 {
		t.Errorf("expected 1 error entry, got %d", len(res.Groups["error"]))
	}
}

func TestApply_ByServiceID(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{
		{ServiceID: "svc-a", Message: "x", Timestamp: time.Now()},
		{ServiceID: "svc-b", Message: "y", Timestamp: time.Now()},
		{ServiceID: "svc-a", Message: "z", Timestamp: time.Now()},
	})

	res, err := pivot.Apply(snap, "service_id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Groups["svc-a"]) != 2 {
		t.Errorf("expected 2 svc-a entries, got %d", len(res.Groups["svc-a"]))
	}
}

func TestApply_MissingFieldFallsToNone(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{
		{Message: "no-level", Timestamp: time.Now()},
	})

	res, err := pivot.Apply(snap, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Groups["(none)"]) != 1 {
		t.Errorf("expected 1 entry under '(none)', got %d", len(res.Groups["(none)"]))
	}
}

func TestApply_ByTagField(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{
		{Message: "a", Timestamp: time.Now(), Tags: map[string]string{"env": "prod"}},
		{Message: "b", Timestamp: time.Now(), Tags: map[string]string{"env": "staging"}},
		{Message: "c", Timestamp: time.Now(), Tags: map[string]string{"env": "prod"}},
	})

	res, err := pivot.Apply(snap, "env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Groups["prod"]) != 2 {
		t.Errorf("expected 2 prod entries, got %d", len(res.Groups["prod"]))
	}
	if len(res.Keys) != 2 {
		t.Errorf("expected 2 sorted keys, got %d", len(res.Keys))
	}
}

func TestApply_NilSnapshot(t *testing.T) {
	_, err := pivot.Apply(nil, "level")
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}

func TestApply_EmptyField(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{})
	_, err := pivot.Apply(snap, "")
	if err == nil {
		t.Error("expected error for empty field")
	}
}
