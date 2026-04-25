package validate_test

import (
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
	"github.com/yourorg/logsnap/internal/validate"
)

func makeSnap(label string) *snapshot.Snapshot {
	snap := snapshot.New(label)
	return snap
}

func addEntry(snap *snapshot.Snapshot, ts time.Time, level, service, msg string) {
	snap.AddEntry(snapshot.Entry{
		Timestamp: ts,
		Level:     level,
		ServiceID: service,
		Message:   msg,
	})
}

func TestCheck_ValidSnapshot(t *testing.T) {
	snap := makeSnap("prod-v1")
	addEntry(snap, time.Now(), "info", "api", "server started")

	res, err := validate.Check(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Valid {
		t.Errorf("expected valid, got errors: %v", res.Errors)
	}
}

func TestCheck_MissingMessage(t *testing.T) {
	snap := makeSnap("prod-v1")
	addEntry(snap, time.Now(), "info", "api", "")

	res, _ := validate.Check(snap)
	if res.Valid {
		t.Error("expected invalid due to missing message")
	}
	if len(res.Errors) == 0 {
		t.Error("expected at least one error")
	}
}

func TestCheck_MissingTimestamp(t *testing.T) {
	snap := makeSnap("prod-v1")
	addEntry(snap, time.Time{}, "info", "api", "some message")

	res, _ := validate.Check(snap)
	if res.Valid {
		t.Error("expected invalid due to zero timestamp")
	}
}

func TestCheck_UnrecognizedLevel(t *testing.T) {
	snap := makeSnap("prod-v1")
	addEntry(snap, time.Now(), "verbose", "svc", "a message")

	res, _ := validate.Check(snap)
	if !res.Valid {
		t.Error("unrecognized level should be a warning, not an error")
	}
	if len(res.Warnings) == 0 {
		t.Error("expected a warning for unrecognized level")
	}
}

func TestCheck_NoLabel(t *testing.T) {
	snap := makeSnap("")
	addEntry(snap, time.Now(), "info", "svc", "hello")

	res, _ := validate.Check(snap)
	if !res.Valid {
		t.Error("missing label should only be a warning")
	}
	if len(res.Warnings) == 0 {
		t.Error("expected warning for missing label")
	}
}

func TestCheck_NilSnapshot(t *testing.T) {
	_, err := validate.Check(nil)
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}
