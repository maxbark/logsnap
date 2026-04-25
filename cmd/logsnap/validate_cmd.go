package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/logsnap/internal/snapshot"
	"github.com/yourorg/logsnap/internal/validate"
)

func newValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate <snapshot>",
		Short: "Check a snapshot for structural and semantic errors",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snap, err := snapshot.Load(args[0])
			if err != nil {
				return fmt.Errorf("loading snapshot: %w", err)
			}

			result, err := validate.Check(snap)
			if err != nil {
				return fmt.Errorf("validation error: %w", err)
			}

			w := cmd.OutOrStdout()

			if len(result.Warnings) > 0 {
				fmt.Fprintln(w, "Warnings:")
				for _, warn := range result.Warnings {
					fmt.Fprintf(w, "  ⚠  %s\n", warn)
				}
			}

			if len(result.Errors) > 0 {
				fmt.Fprintln(w, "Errors:")
				for _, e := range result.Errors {
					fmt.Fprintf(w, "  ✗  %s\n", e)
				}
			}

			if result.Valid {
				fmt.Fprintln(w, "✓ Snapshot is valid.")
				return nil
			}

			fmt.Fprintln(w, "✗ Snapshot is invalid.")
			os.Exit(1)
			return nil
		},
	}
	return cmd
}
