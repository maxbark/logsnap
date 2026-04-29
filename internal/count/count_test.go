package count_test

import (
	"testing"
	"time"

	"logsnap/internal/count"
	"logsnap/internal/snapshot"
)

func makeSnap(entries []snapshot.Entry) *snapshot.Snapshot {
	s := snapshot.New("test-snap", "test-source")
	s.Entries = entries
	return s
}

func entry(level, service, msg string) snapshot.Entry {
	return snapshot.Entry{
		Timestamp: time.Now(),
		Level:     level,
		ServiceID: service,
		Message:   msg,
	}
}

func TestApply_ByLevel(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{
		entry("info", "svc-a", "started"),
		entry("error", "svc-b", "failed"),
		entry("info", "svc-a", "stopped"),
		entry("warn", "svc-c", "slow"),
	})

	res, err := count.Apply(snap, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 4 {
		t.Errorf("expected total 4, got %d", res.Total)
	}
	if res.Counts["info"] != 2 {
		t.Errorf("expected info=2, got %d", res.Counts["info"])
	}
	if res.Counts["error"] != 1 {
		t.Errorf("expected error=1, got %d", res.Counts["error"])
	}
}

func TestApply_ByServiceID(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{
		entry("info", "svc-a", "m1"),
		entry("info", "svc-a", "m2"),
		entry("info", "svc-b", "m3"),
	})

	res, err := count.Apply(snap, "service_id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Counts["svc-a"] != 2 {
		t.Errorf("expected svc-a=2, got %d", res.Counts["svc-a"])
	}
	if res.Counts["svc-b"] != 1 {
		t.Errorf("expected svc-b=1, got %d", res.Counts["svc-b"])
	}
}

func TestApply_EmptyServiceFallsToNone(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{
		entry("info", "", "no-service"),
	})
	res, err := count.Apply(snap, "service_id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Counts["(none)"] != 1 {
		t.Errorf("expected (none)=1, got %d", res.Counts["(none)"])
	}
}

func TestApply_InvalidField(t *testing.T) {
	snap := makeSnap(nil)
	_, err := count.Apply(snap, "unknown")
	if err == nil {
		t.Error("expected error for unsupported field")
	}
}

func TestApply_NilSnapshot(t *testing.T) {
	_, err := count.Apply(nil, "level")
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}

func TestSortedKeys_OrderedByCountDesc(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{
		entry("info", "a", "m"),
		entry("info", "a", "m"),
		entry("error", "b", "m"),
	})
	res, _ := count.Apply(snap, "level")
	keys := res.SortedKeys()
	if len(keys) == 0 || keys[0] != "info" {
		t.Errorf("expected first key to be 'info', got %v", keys)
	}
}
