package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourorg/logsnap/internal/replay"
	"github.com/yourorg/logsnap/internal/snapshot"
)

func newReplayCmd() *cobra.Command {
	var (
		format    string
		delayMs   int
	)

	cmd := &cobra.Command{
		Use:   "replay <snapshot>",
		Short: "Replay a snapshot's log entries to stdout",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snap, err := snapshot.Load(args[0])
			if err != nil {
				return fmt.Errorf("failed to load snapshot: %w", err)
			}

			opts := replay.Options{
				Format: format,
				Delay:  time.Duration(delayMs) * time.Millisecond,
			}

			if err := replay.Run(os.Stdout, snap, opts); err != nil {
				return fmt.Errorf("replay failed: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, json, logfmt")
	cmd.Flags().IntVarP(&delayMs, "delay", "d", 0, "Delay in milliseconds between entries")

	return cmd
}
