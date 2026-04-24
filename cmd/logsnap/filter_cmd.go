package main

import (
	"fmt"
	"os"

	"github.com/logsnap/internal/filter"
	"github.com/logsnap/internal/snapshot"
	"github.com/spf13/cobra"
)

func newFilterCmd() *cobra.Command {
	var (
		inputPath  string
		outputPath string
		level      string
		serviceID  string
		messageRe  string
		fieldKey   string
		fieldValue string
	)

	cmd := &cobra.Command{
		Use:   "filter",
		Short: "Filter a snapshot by level, service, message pattern, or field",
		RunE: func(cmd *cobra.Command, args []string) error {
			snap, err := snapshot.Load(inputPath)
			if err != nil {
				return fmt.Errorf("loading snapshot: %w", err)
			}

			opts := filter.Options{
				Level:      level,
				ServiceID:  serviceID,
				MessageRe:  messageRe,
				FieldKey:   fieldKey,
				FieldValue: fieldValue,
			}

			filtered, err := filter.Apply(snap, opts)
			if err != nil {
				return fmt.Errorf("applying filter: %w", err)
			}

			if outputPath == "" {
				for _, e := range filtered.Entries {
					fmt.Fprintf(os.Stdout, "[%s] %s %s\n", e.Level, e.ServiceID, e.Message)
				}
				return nil
			}

			if err := filtered.Save(outputPath); err != nil {
				return fmt.Errorf("saving filtered snapshot: %w", err)
			}
			fmt.Fprintf(os.Stdout, "Filtered snapshot saved to %s (%d entries)\n", outputPath, len(filtered.Entries))
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputPath, "input", "i", "", "Path to input snapshot file (required)")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to save filtered snapshot (optional, prints to stdout if omitted)")
	cmd.Flags().StringVar(&level, "level", "", "Filter by log level (e.g. info, error)")
	cmd.Flags().StringVar(&serviceID, "service", "", "Filter by service ID")
	cmd.Flags().StringVar(&messageRe, "message", "", "Filter by message regex")
	cmd.Flags().StringVar(&fieldKey, "field-key", "", "Filter by field key presence")
	cmd.Flags().StringVar(&fieldValue, "field-value", "", "Filter by field value (requires --field-key)")
	_ = cmd.MarkFlagRequired("input")
	return cmd
}
