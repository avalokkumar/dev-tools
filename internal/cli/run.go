package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/devforge/devforge/internal/mcpserver"
	"github.com/devforge/devforge/internal/tools"
)

// newRunCmd builds `devforge run <tool>_<op>` — a generic dispatcher that
// invokes any Operation registered with the Registry. Useful for tools that
// don't have a dedicated cobra subcommand yet (yaml, csv, diff, regex, cron,
// tz, faker), and for scripting AI agents that prefer a uniform shape.
func newRunCmd(_ Stdio) *cobra.Command {
	var (
		argsFile string
		argsRaw  string
	)
	cmd := &cobra.Command{
		Use:   "run <operation>",
		Short: "Invoke any registered Operation by name with JSON args",
		Long: "Reads JSON arguments from --args, --args-file, or stdin and invokes " +
			"the named Operation. Operation names are <tool>_<op>, e.g. " +
			"yaml_format, jwt_decode, faker_generate.\n\n" +
			"Use `devforge run --list` to enumerate registered Operations.",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reg := mcpserver.NewRegistry()
			if err := tools.Register(reg); err != nil {
				return err
			}
			plugins := loadPluginsInto(cmd.Context(), reg, "")
			defer func() {
				for _, p := range plugins {
					_ = p.Close()
				}
			}()

			list, _ := cmd.Flags().GetBool("list")
			if list || len(args) == 0 {
				return printOpsList(cmd.OutOrStdout(), reg)
			}

			name := args[0]
			op, ok := reg.Get(name)
			if !ok {
				return fmt.Errorf("operation %q not found (use --list to see options)", name)
			}

			argBytes, err := readArgs(cmd, argsRaw, argsFile)
			if err != nil {
				return err
			}
			if len(argBytes) == 0 {
				argBytes = []byte("{}")
			}
			out, err := op.Handler(cmd.Context(), json.RawMessage(argBytes))
			if err != nil {
				return err
			}
			_, _ = cmd.OutOrStdout().Write(out)
			_, _ = fmt.Fprintln(cmd.OutOrStdout())
			return nil
		},
	}
	cmd.Flags().StringVar(&argsRaw, "args", "", "JSON argument literal")
	cmd.Flags().StringVar(&argsFile, "args-file", "", "read JSON arguments from FILE")
	cmd.Flags().Bool("list", false, "list registered Operations")
	return cmd
}

func readArgs(cmd *cobra.Command, raw, file string) ([]byte, error) {
	switch {
	case raw != "":
		return []byte(raw), nil
	case file != "":
		return os.ReadFile(file)
	default:
		// If stdin is a terminal, return empty.
		st, _ := os.Stdin.Stat()
		if st != nil && st.Mode()&os.ModeCharDevice != 0 {
			return nil, nil
		}
		return io.ReadAll(cmd.InOrStdin())
	}
}

func printOpsList(w io.Writer, reg *mcpserver.Registry) error {
	for _, op := range reg.List() {
		fmt.Fprintf(w, "%-24s  %s\n", op.Name(), op.Description)
	}
	return nil
}
