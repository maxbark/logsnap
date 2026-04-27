package sample_test

import (
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/sample"
	"github.com/yourorg/logsnap/internal/snapshot"
)

func makeSnap(n int) *snapshot.Snapshot {
	snap := snapshot.New("test", "test-source")
	for i := 0; i < n; i++ {
		snap.AddEntry(snapshot.Entry{
			Timestamp: time.Now(),
			Level:     "info",
			Message:   fmt.Sprintf("message %d", i),
			ServiceID: "svc",
		})
	}
	return snap
}

import "fmt"

func TestApply_ReturnsNEntries(t *testing.T) {
	snap := makeSnap(20)
	opts := sample.Options{N: 5, Deterministic: true, Seed: 42}
	out, err := sample.Apply(snap, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 5 {
		t.Errorf("expected 5 entries, got %d", len(out.Entries))
	}
}

func TestApply_NGreaterThanTotal(t *testing.T) {
	snap := makeSnap(3)
	opts := sample.Options{N: 10, Deterministic: true, Seed: 1}
	out, err := sample.Apply(snap, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out.Entries))
	}
}

func TestApply_PreservesOrder(t *testing.T) {
	snap := makeSnap(10)
	opts := sample.Options{N: 4, Deterministic: true, Seed: 7}
	out, err := sample.Apply(snap, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(out.Entries); i++ {
		if out.Entries[i].Timestamp.Before(out.Entries[i-1].Timestamp) {
			t.Errorf("entries are not in order at index %d", i)
		}
	}
}

func TestApply_NilSnapshot(t *testing.T) {
	_, err := sample.Apply(nil, sample.DefaultOptions())
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}

func TestApply_InvalidN(t *testing.T) {
	snap := makeSnap(5)
	_, err := sample.Apply(snap, sample.Options{N: 0})
	if err == nil {
		t.Error("expected error for N=0")
	}
}

func TestApply_Deterministic(t *testing.T) {
	snap := makeSnap(50)
	opts := sample.Options{N: 10, Deterministic: true, Seed: 99}
	a, _ := sample.Apply(snap, opts)
	b, _ := sample.Apply(snap, opts)
	if len(a.Entries) != len(b.Entries) {
		t.Fatal("deterministic runs returned different lengths")
	}
	for i := range a.Entries {
		if a.Entries[i].Message != b.Entries[i].Message {
			t.Errorf("entry %d differs: %q vs %q", i, a.Entries[i].Message, b.Entries[i].Message)
		}
	}
}
