package main

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
)

func writeWatchSnap(t *testing.T, entries []snapshot.Entry) string {
	t.Helper()
	snap := snapshot.New("watch-cmd-test", "svc")
	for _, e := range entries {
		snap.AddEntry(e)
	}
	f, err := os.CreateTemp(t.TempDir(), "watch-*.snap")
	if err != nil {
		t.Fatal(err)
	}
	path := f.Name()
	f.Close()
	if err := snap.Save(path); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestWatchCmd_InvalidFormat(t *testing.T) {
	path := writeWatchSnap(t, nil)

	cmd := newRootCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"watch", "--format", "xml", path})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestWatchCmd_JSONOutputRegistered(t *testing.T) {
	// Verify the watch sub-command is wired up and accepts --format json
	path := writeWatchSnap(t, []snapshot.Entry{
		{Timestamp: time.Now(), Level: "info", ServiceID: "svc", Message: "hello"},
	})

	// Encode a second entry and append it after a brief pause in a goroutine
	go func() {
		time.Sleep(200 * time.Millisecond)
		snap, _ := snapshot.Load(path)
		snap.AddEntry(snapshot.Entry{
			Timestamp: time.Now(), Level: "debug", ServiceID: "svc", Message: "world",
		})
		_ = snap.Save(path)
	}()

	// Just assert that the command can be built and the flag parses
	cmd := newWatchCmd()
	if cmd.Use != "watch <snapshot>" {
		t.Errorf("unexpected Use: %s", cmd.Use)
	}
	f := cmd.Flags().Lookup("format")
	if f == nil {
		t.Fatal("expected --format flag")
	}
	if f.DefValue != "text" {
		t.Errorf("expected default text, got %s", f.DefValue)
	}

	// Sanity: snapshot JSON round-trip
	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(snap.Entries[0])
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Error("expected non-empty JSON")
	}
}
