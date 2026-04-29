package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/user/logsnap/internal/pivot"
	"github.com/user/logsnap/internal/snapshot"
)

func newPivotCmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "pivot <snapshot> <field>",
		Short: "Group snapshot entries by a field and display counts",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			snapPath := args[0]
			field := args[1]

			snap, err := snapshot.Load(snapPath)
			if err != nil {
				return fmt.Errorf("failed to load snapshot: %w", err)
			}

			res, err := pivot.Apply(snap, field)
			if err != nil {
				return fmt.Errorf("pivot failed: %w", err)
			}

			switch outputFormat {
			case "json":
				type jsonGroup struct {
					Value string `json:"value"`
					Count int    `json:"count"`
				}
				groups := make([]jsonGroup, 0, len(res.Keys))
				for _, k := range res.Keys {
					groups = append(groups, jsonGroup{Value: k, Count: len(res.Groups[k])})
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(map[string]interface{}{
					"field":  res.Field,
					"groups": groups,
				})
			default:
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintf(w, "FIELD: %s\n\n", res.Field)
				fmt.Fprintln(w, "VALUE\tCOUNT")
				for _, k := range res.Keys {
					fmt.Fprintf(w, "%s\t%d\n", k, len(res.Groups[k]))
				}
				return w.Flush()
			}
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "format", "f", "text", "Output format: text or json")
	return cmd
}
