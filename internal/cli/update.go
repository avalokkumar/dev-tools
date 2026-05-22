package cli

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/devforge/devforge/internal/update"
)

func newUpdateCmd(_ Stdio) *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Check for or apply binary updates",
		Long: "Resolves the latest available release via the configured source.\n" +
			"Without a release channel wired in, this is a no-op that explains how to " +
			"configure one. Use --dry-run to simulate.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			plan, err := update.Apply(cmd.Context(), update.ApplyOptions{
				DryRun: dryRun,
				Out:    cmd.OutOrStdout(),
			})
			if err != nil {
				return err
			}
			if asJSON, _ := cmd.Flags().GetBool("json"); asJSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(plan)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "describe the update without applying it")
	return cmd
}
