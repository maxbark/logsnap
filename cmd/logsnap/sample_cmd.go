package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/logsnap/internal/sample"
	"github.com/yourorg/logsnap/internal/snapshot"
)

func newSampleCmd() *cobra.Command {
	var (
		outputPath    string
		n             int
		seed          int64
		deterministic bool
	)

	cmd := &cobra.Command{
		Use:   "sample <snapshot>",
		Short: "Randomly sample N entries from a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snap, err := snapshot.Load(args[0])
			if err != nil {
				return fmt.Errorf("loading snapshot: %w", err)
			}

			opts := sample.Options{
				N:             n,
				Seed:          seed,
				Deterministic: deterministic,
			}

			out, err := sample.Apply(snap, opts)
			if err != nil {
				return fmt.Errorf("sampling: %w", err)
			}

			if outputPath != "" {
				if err := out.Save(outputPath); err != nil {
					return fmt.Errorf("saving output: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Sampled %d entries saved to %s\n", len(out.Entries), outputPath)
				return nil
			}

			for _, e := range out.Entries {
				fmt.Fprintf(cmd.OutOrStdout(), "%s [%s] %s\n", e.Timestamp.Format("2006-01-02T15:04:05Z07:00"), e.Level, e.Message)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&n, "count", "n", 10, "Number of entries to sample")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Save sampled snapshot to file")
	cmd.Flags().Int64Var(&seed, "seed", 0, "Random seed for deterministic sampling")
	cmd.Flags().BoolVar(&deterministic, "deterministic", false, "Use seed for reproducible sampling")

	_ = os.Stderr // suppress unused import
	return cmd
}
