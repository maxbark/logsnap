package annotate_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/annotate"
	"github.com/yourorg/logsnap/internal/snapshot"
)

func makeSnap(label string, messages ...string) *snapshot.Snapshot {
	s := snapshot.New(label)
	for _, m := range messages {
		s.AddEntry(snapshot.Entry{
			Timestamp: time.Now().UTC(),
			Level:     "info",
			Message:   m,
			ServiceID: "svc",
			Fields:    map[string]string{},
		})
	}
	return s
}

func TestApply_SingleEntry(t *testing.T) {
	s := makeSnap("test", "alpha", "beta")
	out, err := annotate.Apply(s, annotate.Options{EntryIndex: 0, Note: "check this"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.Entries[0].Fields["_annotation"], "check this") {
		t.Errorf("expected annotation on entry 0, got %q", out.Entries[0].Fields["_annotation"])
	}
	if _, ok := out.Entries[1].Fields["_annotation"]; ok {
		t.Error("entry 1 should not be annotated")
	}
}

func TestApply_AllEntries(t *testing.T) {
	s := makeSnap("test", "alpha", "beta", "gamma")
	out, err := annotate.Apply(s, annotate.Options{EntryIndex: -1, Note: "global note"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, e := range out.Entries {
		if !strings.Contains(e.Fields["_annotation"], "global note") {
			t.Errorf("entry %d missing annotation", i)
		}
	}
}

func TestApply_WithAuthor(t *testing.T) {
	s := makeSnap("test", "msg")
	out, err := annotate.Apply(s, annotate.Options{EntryIndex: 0, Note: "looks odd", Author: "alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	anno := out.Entries[0].Fields["_annotation"]
	if !strings.Contains(anno, "alice: looks odd") {
		t.Errorf("expected author prefix, got %q", anno)
	}
}

func TestApply_AppendsByDefault(t *testing.T) {
	s := makeSnap("test", "msg")
	s.Entries[0].Fields["_annotation"] = "[2024-01-01T00:00:00Z] first"
	out, err := annotate.Apply(s, annotate.Options{EntryIndex: 0, Note: "second"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	anno := out.Entries[0].Fields["_annotation"]
	if !strings.Contains(anno, "first") || !strings.Contains(anno, "second") {
		t.Errorf("expected both notes, got %q", anno)
	}
}

func TestApply_Overwrite(t *testing.T) {
	s := makeSnap("test", "msg")
	s.Entries[0].Fields["_annotation"] = "old note"
	out, err := annotate.Apply(s, annotate.Options{EntryIndex: 0, Note: "new note", Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	anno := out.Entries[0].Fields["_annotation"]
	if strings.Contains(anno, "old note") {
		t.Errorf("old note should have been overwritten, got %q", anno)
	}
}

func TestApply_NilSnapshot(t *testing.T) {
	_, err := annotate.Apply(nil, annotate.Options{Note: "x"})
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}

func TestApply_EmptyNote(t *testing.T) {
	s := makeSnap("test", "msg")
	_, err := annotate.Apply(s, annotate.Options{EntryIndex: 0, Note: ""})
	if err == nil {
		t.Error("expected error for empty note")
	}
}

func TestApply_OutOfRangeIndex(t *testing.T) {
	s := makeSnap("test", "msg")
	_, err := annotate.Apply(s, annotate.Options{EntryIndex: 99, Note: "x"})
	if err == nil {
		t.Error("expected error for out-of-range index")
	}
}
