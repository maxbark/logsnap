package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/logsnap/internal/snapshot"
	logsort "github.com/yourorg/logsnap/internal/sort"
)

func newSortCmd() *cobra.Command {
	var (
		outputFile  string
		sortBy      string
		descending  bool
	)

	cmd := &cobra.Command{
		Use:   "sort <snapshot>",
		Short: "Sort snapshot entries by a specified field",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snap, err := snapshot.Load(args[0])
			if err != nil {
				return fmt.Errorf("loading snapshot: %w", err)
			}

			field := logsort.Field(sortBy)
			switch field {
			case logsort.FieldTimestamp, logsort.FieldLevel, logsort.FieldService, logsort.FieldMessage:
			default:
				return fmt.Errorf("unknown sort field %q; valid fields: timestamp, level, service, message", sortBy)
			}

			opts := logsort.Options{
				By:         field,
				Descending: descending,
			}

			out, err := logsort.Apply(snap, opts)
			if err != nil {
				return fmt.Errorf("sorting: %w", err)
			}

			if outputFile != "" {
				if err := out.Save(outputFile); err != nil {
					return fmt.Errorf("saving output: %w", err)
				}
				fmt.Fprintf(os.Stdout, "sorted snapshot saved to %s\n", outputFile)
				return nil
			}

			for _, e := range out.Entries {
				fmt.Fprintf(cmd.OutOrStdout(), "%s [%s] (%s) %s\n",
					e.Timestamp.Format("2006-01-02T15:04:05Z"),
					e.Level, e.ServiceID, e.Message)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "save sorted snapshot to file")
	cmd.Flags().StringVarP(&sortBy, "by", "b", "timestamp", "field to sort by (timestamp, level, service, message)")
	cmd.Flags().BoolVarP(&descending, "desc", "d", false, "sort in descending order")

	return cmd
}
