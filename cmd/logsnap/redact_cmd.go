package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourorg/logsnap/internal/redact"
	"github.com/yourorg/logsnap/internal/snapshot"
)

func newRedactCmd() *cobra.Command {
	var fields []string
	var patterns []string
	var mask string
	var output string

	cmd := &cobra.Command{
		Use:   "redact <snapshot>",
		Short: "Mask sensitive fields or patterns in a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snap, err := snapshot.Load(args[0])
			if err != nil {
				return fmt.Errorf("failed to load snapshot: %w", err)
			}

			if len(fields) == 0 && len(patterns) == 0 {
				return fmt.Errorf("at least one --field or --pattern must be specified")
			}

			opts := redact.Options{
				Fields:   fields,
				Patterns: patterns,
				Mask:     mask,
			}

			result, err := redact.Apply(snap, opts)
			if err != nil {
				return fmt.Errorf("redact failed: %w", err)
			}

			if output != "" {
				if err := result.Save(output); err != nil {
					return fmt.Errorf("failed to save output: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Redacted snapshot saved to %s\n", output)
			} else {
				for _, e := range result.Entries {
					parts := []string{e.Timestamp.Format("2006-01-02T15:04:05Z"), e.Level, e.ServiceID, e.Message}
					fmt.Fprintln(cmd.OutOrStdout(), strings.Join(parts, " | "))
				}
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&fields, "field", nil, "Field names to redact (repeatable)")
	cmd.Flags().StringSliceVar(&patterns, "pattern", nil, "Regex patterns to redact (repeatable)")
	cmd.Flags().StringVar(&mask, "mask", "", "Replacement string (default: [REDACTED])")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Save redacted snapshot to file")
	return cmd
}
