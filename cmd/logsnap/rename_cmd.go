package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourorg/logsnap/internal/rename"
	"github.com/yourorg/logsnap/internal/snapshot"
)

func newRenameCmd() *cobra.Command {
	var (
		label     string
		source    string
		tagsRaw   []string
		mergeTags bool
		output    string
	)

	cmd := &cobra.Command{
		Use:   "rename <snapshot>",
		Short: "Rename or relabel a snapshot's metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snap, err := snapshot.Load(args[0])
			if err != nil {
				return fmt.Errorf("rename: failed to load snapshot: %w", err)
			}

			tags := make(map[string]string)
			for _, raw := range tagsRaw {
				parts := strings.SplitN(raw, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("rename: invalid tag format %q, expected key=value", raw)
				}
				tags[parts[0]] = parts[1]
			}

			opts := rename.Options{
				Label:     label,
				Source:    source,
				Tags:      tags,
				MergeTags: mergeTags,
			}

			out, err := rename.Apply(snap, opts)
			if err != nil {
				return fmt.Errorf("rename: %w", err)
			}

			dest := output
			if dest == "" {
				dest = args[0]
			}

			if err := out.Save(dest); err != nil {
				return fmt.Errorf("rename: failed to save snapshot: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "snapshot saved to %s (label=%q source=%q)\n",
				dest, out.Meta.Label, out.Meta.Source)
			return nil
		},
	}

	cmd.Flags().StringVar(&label, "label", "", "new label for the snapshot")
	cmd.Flags().StringVar(&source, "source", "", "new source identifier for the snapshot")
	cmd.Flags().StringArrayVar(&tagsRaw, "tag", nil, "metadata tag in key=value format (repeatable)")
	cmd.Flags().BoolVar(&mergeTags, "merge-tags", true, "merge new tags with existing tags (default true)")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output path (defaults to overwriting input)")

	return cmd
}
