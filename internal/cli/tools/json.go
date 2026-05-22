package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/devforge/devforge/pkg/jsonfmt"
)

// NewJSONCmd returns the `devforge json` Tool command tree.
func NewJSONCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "json",
		Short: "Format and validate JSON",
	}
	cmd.AddCommand(newJSONFmtCmd())
	cmd.AddCommand(newJSONValidateCmd())
	return cmd
}

func newJSONFmtCmd() *cobra.Command {
	var (
		indent          int
		sortKeys        bool
		trailingNewline bool
		inFile          string
	)
	cmd := &cobra.Command{
		Use:   "fmt",
		Short: "Pretty-print or compact JSON (reads stdin or --in)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			input, err := readInputBytes(cmd, inFile)
			if err != nil {
				return err
			}
			res, err := jsonfmt.Format(input, jsonfmt.FormatOptions{
				Indent: indent, SortKeys: sortKeys, TrailingNewline: trailingNewline,
			})
			if err != nil {
				return err
			}
			asJSON, _ := cmd.Flags().GetBool("json")
			out := cmd.OutOrStdout()
			if asJSON {
				return jsonEncode(out, map[string]any{
					"output":      string(res.Output),
					"diagnostics": res.Diagnostics,
				})
			}
			_, _ = out.Write(res.Output)
			for _, d := range res.Diagnostics {
				fmt.Fprintf(cmd.ErrOrStderr(), "diagnostic: [%s] %s\n", d.Code, d.Message)
			}
			if hasErrorDiag(res.Diagnostics) {
				return errExitOne
			}
			return nil
		},
	}
	cmd.Flags().IntVarP(&indent, "indent", "i", 2, "spaces per indent (0 = compact)")
	cmd.Flags().BoolVar(&sortKeys, "sort-keys", false, "sort object keys")
	cmd.Flags().BoolVar(&trailingNewline, "trailing-newline", false, "append a newline")
	cmd.Flags().StringVarP(&inFile, "in", "f", "", "read input from FILE instead of stdin")
	return cmd
}

func newJSONValidateCmd() *cobra.Command {
	var (
		inFile     string
		schemaFile string
	)
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate JSON syntax (and optional schema)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			input, err := readInputBytes(cmd, inFile)
			if err != nil {
				return err
			}
			var schema []byte
			if schemaFile != "" {
				schema, err = os.ReadFile(schemaFile)
				if err != nil {
					return err
				}
			}
			res, err := jsonfmt.Validate(input, schema)
			if err != nil {
				return err
			}
			asJSON, _ := cmd.Flags().GetBool("json")
			out := cmd.OutOrStdout()
			if asJSON {
				return jsonEncode(out, res)
			}
			if res.Valid {
				fmt.Fprintln(out, "valid")
			} else {
				fmt.Fprintln(out, "invalid")
			}
			for _, d := range res.Diagnostics {
				fmt.Fprintf(cmd.ErrOrStderr(), "[%s] %s\n", d.Code, d.Message)
			}
			if !res.Valid {
				return errExitOne
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&inFile, "in", "f", "", "input file")
	cmd.Flags().StringVarP(&schemaFile, "schema", "s", "", "JSON Schema file")
	return cmd
}

// readInputBytes reads from --in file if set, else from cobra's InOrStdin.
func readInputBytes(cmd *cobra.Command, file string) ([]byte, error) {
	if file != "" {
		return os.ReadFile(file)
	}
	return io.ReadAll(cmd.InOrStdin())
}

func jsonEncode(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
