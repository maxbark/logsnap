package watch_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
	"github.com/yourorg/logsnap/internal/watch"
)

func makeWatchSnap(t *testing.T, entries int) (string, *snapshot.Snapshot) {
	t.Helper()
	snap := snapshot.New("watch-test", "svc")
	for i := 0; i < entries; i++ {
		snap.AddEntry(snapshot.Entry{
			Timestamp: time.Now(),
			Level:     "info",
			ServiceID: "svc",
			Message:   "entry",
		})
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
	return path, snap
}

func TestRun_EmitsNewEntries(t *testing.T) {
	path, _ := makeWatchSnap(t, 1)

	out, err := os.CreateTemp(t.TempDir(), "out-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	opts := watch.Options{PollInterval: 100 * time.Millisecond, Format: "text"}

	// Append a new entry after a short delay
	go func() {
		time.Sleep(300 * time.Millisecond)
		snap, _ := snapshot.Load(path)
		snap.AddEntry(snapshot.Entry{
			Timestamp: time.Now(),
			Level:     "warn",
			ServiceID: "svc",
			Message:   "new-entry",
		})
		_ = snap.Save(path)
		time.Sleep(400 * time.Millisecond)
		cancel()
	}()

	_ = watch.Run(ctx, path, out, opts)

	_ = out.Sync()
	b, _ := os.ReadFile(out.Name())
	if !strings.Contains(string(b), "new-entry") {
		t.Errorf("expected new-entry in output, got: %s", string(b))
	}
}

func TestRun_JSONFormat(t *testing.T) {
	path, _ := makeWatchSnap(t, 0)

	out, _ := os.CreateTemp(t.TempDir(), "out-*.txt")
	defer out.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	opts := watch.Options{PollInterval: 100 * time.Millisecond, Format: "json"}

	go func() {
		time.Sleep(200 * time.Millisecond)
		snap, _ := snapshot.Load(path)
		snap.AddEntry(snapshot.Entry{
			Timestamp: time.Now(), Level: "error", ServiceID: "svc", Message: "boom",
		})
		_ = snap.Save(path)
		time.Sleep(400 * time.Millisecond)
		cancel()
	}()

	_ = watch.Run(ctx, path, out, opts)
	_ = out.Sync()
	b, _ := os.ReadFile(out.Name())
	if !strings.Contains(string(b), `"message":"boom"`) {
		t.Errorf("expected JSON with message boom, got: %s", string(b))
	}
}

func TestRun_MissingFileNoError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	out, _ := os.CreateTemp(t.TempDir(), "out-*.txt")
	defer out.Close()

	opts := watch.DefaultOptions()
	err := watch.Run(ctx, "/nonexistent/path/snap.json", out, opts)
	if err != nil {
		t.Errorf("expected nil error for missing file, got %v", err)
	}
}
