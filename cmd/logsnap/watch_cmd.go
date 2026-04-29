package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourorg/logsnap/internal/watch"
)

func newWatchCmd() *cobra.Command {
	var (
		format   string
		interval int
	)

	cmd := &cobra.Command{
		Use:   "watch <snapshot>",
		Short: "Live-tail a snapshot file, printing new entries as they appear",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			if format != "text" && format != "json" && format != "logfmt" {
				return fmt.Errorf("unsupported format %q: use text, json, or logfmt", format)
			}

			opts := watch.Options{
				PollInterval: time.Duration(interval) * time.Millisecond,
				Format:       format,
			}

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			fmt.Fprintf(cmd.OutOrStdout(), "Watching %s (press Ctrl+C to stop)\n", path)
			return watch.Run(ctx, path, os.Stdout, opts)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, json, logfmt")
	cmd.Flags().IntVarP(&interval, "interval", "i", 500, "Poll interval in milliseconds")
	return cmd
}
