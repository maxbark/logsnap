package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourorg/logsnap/internal/export"
	"github.com/yourorg/logsnap/internal/snapshot"
)

func newExportCmd() *cobra.Command {
	var formatFlag string
	var outputFlag string

	cmd := &cobra.Command{
		Use:   "export <snapshot>",
		Short: "Export a snapshot to JSON, CSV, or text format",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snap, err := snapshot.Load(args[0])
			if err != nil {
				return fmt.Errorf("loading snapshot: %w", err)
			}

			fmt := export.Format(strings.ToLower(formatFlag))

			if outputFlag != "" {
				f, err := os.Create(outputFlag)
				if err != nil {
					return fmt.Errorf("creating output file: %w", err)
				}
				defer f.Close()
				if err := export.Write(snap, fmt, f); err != nil {
					return fmt.Errorf("exporting snapshot: %w", err)
				}
				cmd.Printf("exported to %s\n", outputFlag)
				return nil
			}

			return export.Write(snap, fmt, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&formatFlag, "format", "f", "text", "output format: json, csv, text")
	cmd.Flags().StringVarP(&outputFlag, "output", "o", "", "write output to file instead of stdout")

	return cmd
}
