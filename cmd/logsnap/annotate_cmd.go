package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/logsnap/internal/annotate"
	"github.com/yourorg/logsnap/internal/snapshot"
)

func newAnnotateCmd() *cobra.Command {
	var (
		note      string
		author    string
		index     int
		overwrite bool
		output    string
	)

	cmd := &cobra.Command{
		Use:   "annotate <snapshot>",
		Short: "Attach a note to one or all entries in a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snap, err := snapshot.Load(args[0])
			if err != nil {
				return fmt.Errorf("loading snapshot: %w", err)
			}

			out, err := annotate.Apply(snap, annotate.Options{
				EntryIndex: index,
				Note:       note,
				Author:     author,
				Overwrite:  overwrite,
			})
			if err != nil {
				return err
			}

			dest := output
			if dest == "" {
				dest = args[0]
			}
			if err := out.Save(dest); err != nil {
				return fmt.Errorf("saving snapshot: %w", err)
			}
			fmt.Fprintf(os.Stdout, "annotation applied → %s\n", dest)
			return nil
		},
	}

	cmd.Flags().StringVarP(&note, "note", "n", "", "annotation text (required)")
	cmd.Flags().StringVarP(&author, "author", "a", "", "optional author label")
	cmd.Flags().IntVarP(&index, "index", "i", -1, "entry index to annotate (-1 = all)")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing annotation")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output path (defaults to input file)")
	_ = cmd.MarkFlagRequired("note")

	return cmd
}
