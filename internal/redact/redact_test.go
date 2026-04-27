package redact_test

import (
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/redact"
	"github.com/yourorg/logsnap/internal/snapshot"
)

func makeSnap() *snapshot.Snapshot {
	snap := snapshot.New("test", "unit")
	snap.AddEntry(snapshot.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "user logged in",
		ServiceID: "auth",
		Fields:    map[string]string{"email": "alice@example.com", "ip": "1.2.3.4", "action": "login"},
	})
	snap.AddEntry(snapshot.Entry{
		Timestamp: time.Now(),
		Level:     "error",
		Message:   "payment failed token=abc123secret",
		ServiceID: "billing",
		Fields:    map[string]string{"amount": "99.99", "token": "abc123secret"},
	})
	return snap
}

func TestApply_RedactsByField(t *testing.T) {
	snap := makeSnap()
	out, err := redact.Apply(snap, redact.Options{Fields: []string{"email", "token"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Fields["email"] != "[REDACTED]" {
		t.Errorf("expected email to be redacted, got %q", out.Entries[0].Fields["email"])
	}
	if out.Entries[0].Fields["ip"] != "1.2.3.4" {
		t.Errorf("non-sensitive field should be preserved")
	}
	if out.Entries[1].Fields["token"] != "[REDACTED]" {
		t.Errorf("expected token to be redacted")
	}
}

func TestApply_RedactsByPattern(t *testing.T) {
	snap := makeSnap()
	out, err := redact.Apply(snap, redact.Options{Patterns: []string{`\d+\.\d+\.\d+\.\d+`}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Fields["ip"] != "[REDACTED]" {
		t.Errorf("expected ip to be redacted by pattern")
	}
	if out.Entries[0].Fields["action"] != "login" {
		t.Errorf("non-matching field should be preserved")
	}
}

func TestApply_CustomMask(t *testing.T) {
	snap := makeSnap()
	out, err := redact.Apply(snap, redact.Options{Fields: []string{"email"}, Mask: "***"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Fields["email"] != "***" {
		t.Errorf("expected custom mask, got %q", out.Entries[0].Fields["email"])
	}
}

func TestApply_InvalidPattern(t *testing.T) {
	snap := makeSnap()
	_, err := redact.Apply(snap, redact.Options{Patterns: []string{`[invalid`}})
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestApply_NilSnapshot(t *testing.T) {
	_, err := redact.Apply(nil, redact.Options{})
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}

func TestApply_PreservesMetadata(t *testing.T) {
	snap := makeSnap()
	out, err := redact.Apply(snap, redact.Options{Fields: []string{"email"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Label != snap.Label || out.Source != snap.Source {
		t.Errorf("snapshot metadata should be preserved")
	}
}
