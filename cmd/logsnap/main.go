// Package main is the entry point for the logsnap CLI tool.
// logsnap captures and diffs structured log output across deployments.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// version is set at build time via -ldflags.
var version = "dev"

func main() {
	root := newRootCmd()
	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// newRootCmd constructs the root cobra command and registers all subcommands.
func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "logsnap",
		Short:   "Capture and diff structured log output across deployments",
		Version: version,
		Long: `logsnap is a CLI tool for capturing, comparing, filtering,
and replaying structured log snapshots across deployments.

Snapshots are stored as JSON files and can be ingested from
JSON or logfmt formatted log streams.

Examples:
  # Ingest logs from stdin
  cat app.log | logsnap ingest --format json --output snap.json

  # Diff two snapshots
  logsnap diff --base before.json --head after.json

  # Filter a snapshot by level
  logsnap filter --input snap.json --level error

  # Export a snapshot to CSV
  logsnap export --input snap.json --format csv --output out.csv

  # Replay a snapshot to stdout
  logsnap replay --input snap.json --format logfmt

  # Summarise a snapshot
  logsnap summary --input snap.json

  # Annotate entries
  logsnap annotate --input snap.json --note "reviewed" --output snap.json

  # Validate a snapshot
  logsnap validate --input snap.json`,
		// Silence default error printing — we handle it in main.
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// Register all subcommands.
	root.AddCommand(
		newIngestCmd(),
		newDiffCmd(),
		newFilterCmd(),
		newExportCmd(),
		newReplayCmd(),
		newSummaryCmd(),
		newAnnotateCmd(),
		newValidateCmd(),
	)

	return root
}
