package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/logsnap/internal/ingest"
	"github.com/yourorg/logsnap/internal/snapshot"
)

func newIngestCmd() *cobra.Command {
	var (
		output    string
		serviceID string
		version   string
		format    string
		inputFile string
	)

	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Ingest structured log lines into a snapshot file",
		RunE: func(cmd *cobra.Command, args []string) error {
			var src *os.File
			if inputFile != "" {
				f, err := os.Open(inputFile)
				if err != nil {
					return fmt.Errorf("opening input file: %w", err)
				}
				defer f.Close()
				src = f
			} else {
				src = os.Stdin
			}

			snap := snapshot.New(serviceID, version)
			opts := ingest.Options{
				ServiceID: serviceID,
				Format:    format,
			}
			if err := ingest.FromReader(src, snap, opts); err != nil {
				return fmt.Errorf("ingesting logs: %w", err)
			}

			if err := snap.Save(output); err != nil {
				return fmt.Errorf("saving snapshot: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Ingested %d entries into %s\n", len(snap.Entries), output)
			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "snapshot.snap", "Output snapshot file path")
	cmd.Flags().StringVar(&serviceID, "service", "", "Service identifier")
	cmd.Flags().StringVar(&version, "version", "unknown", "Deployment version label")
	cmd.Flags().StringVar(&format, "format", "json", "Log format: json or logfmt")
	cmd.Flags().StringVarP(&inputFile, "file", "f", "", "Input log file (defaults to stdin)")

	return cmd
}
