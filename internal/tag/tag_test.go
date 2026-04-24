package tag_test

import (
	"testing"
	"time"

	"github.com/user/logsnap/internal/snapshot"
	"github.com/user/logsnap/internal/tag"
)

func makeSnap() *snapshot.Snapshot {
	snap := snapshot.New("test-snap", "v1.0.0")
	snap.Meta.CapturedAt = time.Now().UTC()
	return snap
}

func TestApply_SingleTag(t *testing.T) {
	snap := makeSnap()
	if err := tag.Apply(snap, []string{"env=production"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := snap.Meta.Tags["env"]; got != "production" {
		t.Errorf("expected env=production, got %q", got)
	}
}

func TestApply_MultipleTags(t *testing.T) {
	snap := makeSnap()
	err := tag.Apply(snap, []string{"env=staging", "region=us-east-1", "team=platform"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := map[string]string{
		"env":    "staging",
		"region": "us-east-1",
		"team":   "platform",
	}
	for k, v := range expected {
		if snap.Meta.Tags[k] != v {
			t.Errorf("tag[%s]: expected %q, got %q", k, v, snap.Meta.Tags[k])
		}
	}
}

func TestApply_OverwritesExistingTag(t *testing.T) {
	snap := makeSnap()
	_ = tag.Apply(snap, []string{"env=dev"})
	_ = tag.Apply(snap, []string{"env=prod"})
	if snap.Meta.Tags["env"] != "prod" {
		t.Errorf("expected overwritten value prod, got %q", snap.Meta.Tags["env"])
	}
}

func TestApply_InvalidFormat(t *testing.T) {
	snap := makeSnap()
	if err := tag.Apply(snap, []string{"noequals"}); err == nil {
		t.Error("expected error for missing '=' but got nil")
	}
}

func TestApply_InvalidKey(t *testing.T) {
	snap := makeSnap()
	if err := tag.Apply(snap, []string{"123bad=value"}); err == nil {
		t.Error("expected error for invalid key but got nil")
	}
}

func TestApply_NilSnapshot(t *testing.T) {
	if err := tag.Apply(nil, []string{"env=prod"}); err == nil {
		t.Error("expected error for nil snapshot but got nil")
	}
}

func TestRemove_ExistingKey(t *testing.T) {
	snap := makeSnap()
	_ = tag.Apply(snap, []string{"env=prod", "region=eu"})
	if err := tag.Remove(snap, []string{"env"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := snap.Meta.Tags["env"]; exists {
		t.Error("expected env tag to be removed")
	}
	if snap.Meta.Tags["region"] != "eu" {
		t.Error("expected region tag to remain")
	}
}

func TestRemove_MissingKeyIsNoOp(t *testing.T) {
	snap := makeSnap()
	if err := tag.Remove(snap, []string{"nonexistent"}); err != nil {
		t.Errorf("unexpected error removing missing key: %v", err)
	}
}

func TestRemove_NilSnapshot(t *testing.T) {
	if err := tag.Remove(nil, []string{"env"}); err == nil {
		t.Error("expected error for nil snapshot but got nil")
	}
}
