package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/logsnap/internal/snapshot"
	"github.com/yourorg/logsnap/internal/summary"
)

func newSummaryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "summary <snapshot>",
		Short: "Print aggregated statistics for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snapPath := args[0]

			snap, err := snapshot.Load(snapPath)
			if err != nil {
				return fmt.Errorf("loading snapshot: %w", err)
			}

			stats, err := summary.Compute(snap)
			if err != nil {
				return fmt.Errorf("computing summary: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Snapshot : %s\n", snap.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "Version  : %s\n", snap.Version)
			fmt.Fprintln(cmd.OutOrStdout())
			summary.Print(cmd.OutOrStdout(), stats)
			return nil
		},
	}
	return cmd
}

func init() {
	_ = os.Stderr // ensure os import used
}
