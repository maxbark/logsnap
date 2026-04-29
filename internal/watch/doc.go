// Package watch implements live-tail polling for logsnap snapshot files.
//
// It periodically re-reads a snapshot file on disk and emits any entries
// appended since the last poll, supporting text, JSON, and logfmt output
// formats. Watching stops cleanly when the provided context is cancelled,
// making it suitable for use with OS signal handling in CLI commands.
//
// Example usage:
//
//	opts := watch.DefaultOptions()
//	opts.Format = "json"
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	if err := watch.Run(ctx, "./my.snap", os.Stdout, opts); err != nil {
//		log.Fatal(err)
//	}
package watch
