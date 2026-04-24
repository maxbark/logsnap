package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourorg/logsnap/internal/diff"
	"github.com/yourorg/logsnap/internal/snapshot"
)

func newDiffCmd() *cobra.Command {
	var outputFile string

	cmd := &cobra.Command{
		Use:   "diff <snapshot-a> <snapshot-b>",
		Short: "Compare two snapshots and display differences",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			snapA, err := snapshot.Load(args[0])
			if err != nil {
				return fmt.Errorf("loading snapshot A: %w", err)
			}

			snapB, err := snapshot.Load(args[1])
			if err != nil {
				return fmt.Errorf("loading snapshot B: %w", err)
			}

			result := diff.Compare(snapA, snapB)

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "STATUS\tSERVICE\tLEVEL\tMESSAGE")
			fmt.Fprintln(w, "------\t-------\t-----\t-------")

			for _, e := range result.Added {
				fmt.Fprintf(w, "ADDED\t%s\t%s\t%s\n", e.ServiceID, e.Level, e.Message)
			}
			for _, e := range result.Removed {
				fmt.Fprintf(w, "REMOVED\t%s\t%s\t%s\n", e.ServiceID, e.Level, e.Message)
			}
			for _, c := range result.Changed {
				fmt.Fprintf(w, "CHANGED\t%s\t%s\t%s -> %s\n", c.Before.ServiceID, c.Before.Level, c.Before.Message, c.After.Message)
			}
			w.Flush()

			if outputFile != "" {
				f, err := os.Create(outputFile)
				if err != nil {
					return fmt.Errorf("creating output file: %w", err)
				}
				defer f.Close()
				tw := tabwriter.NewWriter(f, 0, 0, 2, ' ', 0)
				fmt.Fprintln(tw, "STATUS\tSERVICE\tLEVEL\tMESSAGE")
				for _, e := range result.Added {
					fmt.Fprintf(tw, "ADDED\t%s\t%s\t%s\n", e.ServiceID, e.Level, e.Message)
				}
				for _, e := range result.Removed {
					fmt.Fprintf(tw, "REMOVED\t%s\t%s\t%s\n", e.ServiceID, e.Level, e.Message)
				}
				for _, c := range result.Changed {
					fmt.Fprintf(tw, "CHANGED\t%s\t%s\t%s -> %s\n", c.Before.ServiceID, c.Before.Level, c.Before.Message, c.After.Message)
				}
				tw.Flush()
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write diff results to file")
	return cmd
}
