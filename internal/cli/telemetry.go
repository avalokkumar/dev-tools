package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/devforge/devforge/internal/telemetry"
)

func newTelemetryCmd(_ Stdio) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "telemetry",
		Short: "Inspect local telemetry buffer",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Show whether telemetry is enabled",
		RunE: func(cmd *cobra.Command, _ []string) error {
			r := telemetry.New()
			out := cmd.OutOrStdout()
			if asJSON, _ := cmd.Flags().GetBool("json"); asJSON {
				return json.NewEncoder(out).Encode(map[string]any{"enabled": r.Enabled()})
			}
			if r.Enabled() {
				fmt.Fprintln(out, "telemetry: enabled")
			} else {
				fmt.Fprintln(out, "telemetry: disabled (set DEVFORGE_TELEMETRY=1 to opt in)")
			}
			return nil
		},
	})
	return cmd
}
