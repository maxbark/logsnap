package rename_test

import (
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/rename"
	"github.com/yourorg/logsnap/internal/snapshot"
)

func makeSnap(source, label string) *snapshot.Snapshot {
	s := snapshot.New(source, label)
	s.AddEntry(snapshot.Entry{
		Timestamp: time.Now().UTC(),
		Level:     "info",
		Message:   "hello world",
		ServiceID: "svc-a",
	})
	return s
}

func TestApply_RenamesLabel(t *testing.T) {
	snap := makeSnap("prod", "old-label")
	out, err := rename.Apply(snap, rename.Options{Label: "new-label"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Meta.Label != "new-label" {
		t.Errorf("expected label 'new-label', got %q", out.Meta.Label)
	}
	if out.Meta.Source != "prod" {
		t.Errorf("source should be unchanged, got %q", out.Meta.Source)
	}
}

func TestApply_RenamesSource(t *testing.T) {
	snap := makeSnap("prod", "snap1")
	out, err := rename.Apply(snap, rename.Options{Source: "staging"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Meta.Source != "staging" {
		t.Errorf("expected source 'staging', got %q", out.Meta.Source)
	}
}

func TestApply_MergesTags(t *testing.T) {
	snap := makeSnap("prod", "snap1")
	snap.Meta.Tags = map[string]string{"env": "prod", "team": "platform"}
	out, err := rename.Apply(snap, rename.Options{
		Tags:      map[string]string{"env": "staging", "version": "v2"},
		MergeTags: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Meta.Tags["env"] != "staging" {
		t.Errorf("expected env=staging, got %q", out.Meta.Tags["env"])
	}
	if out.Meta.Tags["team"] != "platform" {
		t.Errorf("expected team=platform to be preserved")
	}
	if out.Meta.Tags["version"] != "v2" {
		t.Errorf("expected version=v2 to be added")
	}
}

func TestApply_ReplacesTags(t *testing.T) {
	snap := makeSnap("prod", "snap1")
	snap.Meta.Tags = map[string]string{"env": "prod"}
	out, err := rename.Apply(snap, rename.Options{
		Tags:      map[string]string{"region": "us-east"},
		MergeTags: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out.Meta.Tags["env"]; ok {
		t.Errorf("expected old tag 'env' to be removed")
	}
	if out.Meta.Tags["region"] != "us-east" {
		t.Errorf("expected region=us-east")
	}
}

func TestApply_NilSnapshotReturnsError(t *testing.T) {
	_, err := rename.Apply(nil, rename.Options{Label: "x"})
	if err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestApply_NoOptionsReturnsError(t *testing.T) {
	snap := makeSnap("prod", "snap1")
	_, err := rename.Apply(snap, rename.Options{})
	if err == nil {
		t.Fatal("expected error when no rename options provided")
	}
}

func TestApply_PreservesEntries(t *testing.T) {
	snap := makeSnap("prod", "snap1")
	out, err := rename.Apply(snap, rename.Options{Label: "renamed"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != len(snap.Entries) {
		t.Errorf("expected %d entries, got %d", len(snap.Entries), len(out.Entries))
	}
}
