// Package ingest provides utilities for parsing structured log lines
// from various formats (JSON, logfmt) and populating a snapshot.Snapshot
// with the resulting entries.
//
// Supported formats:
//
//	"json"   — each line is a JSON object with optional fields:
//	           level, msg/message, time (RFC3339)
//
//	"logfmt" — each line is a sequence of key=value pairs, e.g.:
//	           level=info msg="request handled" time=2024-01-01T00:00:00Z
//
// Example usage:
//
//	snap := snapshot.New("my-service", "v1.2.3")
//	err := ingest.FromReader(os.Stdin, snap, ingest.Options{
//	    ServiceID: "my-service",
//	    Format:    "json",
//	})
package ingest
