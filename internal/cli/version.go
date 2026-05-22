package cli

import (
	"encoding/json"
	"fmt"

	"github.com/devforge/devforge/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCmd(_ Stdio) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, _ []string) error {
			info := version.Get()
			asJSON, _ := cmd.Flags().GetBool("json")
			out := cmd.OutOrStdout()
			if asJSON {
				enc := json.NewEncoder(out)
				enc.SetIndent("", "  ")
				return enc.Encode(info)
			}
			_, err := fmt.Fprintf(out, "devforge %s (commit %s, built %s)\n",
				info.Version, info.Commit, info.Date)
			return err
		},
	}
}
