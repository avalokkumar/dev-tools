// Package tools wraps Engines as cobra subcommands.
//
// Each file here is a thin Adapter: parse flags → build Options → call Engine →
// write Result to stdout. No logic.
package tools

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/devforge/devforge/pkg/uuidx"
)

// NewUUIDCmd returns the `devforge uuid` Tool command tree.
func NewUUIDCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uuid",
		Short: "Generate UUIDs and cryptographic digests",
	}
	cmd.AddCommand(newUUIDGenCmd())
	cmd.AddCommand(newUUIDHashCmd())
	return cmd
}

func newUUIDGenCmd() *cobra.Command {
	var (
		version int
		count   int
		format  string
	)
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate one or more UUIDs",
		RunE: func(cmd *cobra.Command, _ []string) error {
			res, err := uuidx.Generate(uuidx.GenerateOptions{
				Version: version,
				Count:   count,
				Format:  format,
			})
			if err != nil {
				return err
			}
			asJSON, _ := cmd.Flags().GetBool("json")
			out := cmd.OutOrStdout()
			if asJSON {
				enc := json.NewEncoder(out)
				enc.SetIndent("", "  ")
				return enc.Encode(res)
			}
			for _, d := range res.Diagnostics {
				fmt.Fprintf(cmd.ErrOrStderr(), "diagnostic: [%s] %s\n", d.Code, d.Message)
			}
			for _, v := range res.Values {
				fmt.Fprintln(out, v)
			}
			if hasErrorDiag(res.Diagnostics) {
				return errExitOne
			}
			return nil
		},
	}
	cmd.Flags().IntVarP(&version, "version", "v", 4, "UUID version (4 or 7)")
	cmd.Flags().IntVarP(&count, "count", "n", 1, "number of UUIDs to generate")
	cmd.Flags().StringVar(&format, "format", "std", "output format: std|compact|urn")
	return cmd
}

func newUUIDHashCmd() *cobra.Command {
	var (
		algos    []string
		encoding string
		input    string
	)
	cmd := &cobra.Command{
		Use:   "hash",
		Short: "Compute cryptographic digests of an input string",
		RunE: func(cmd *cobra.Command, _ []string) error {
			res, err := uuidx.Hash([]byte(input), uuidx.HashOptions{
				Algos:    algos,
				Encoding: encoding,
			})
			if err != nil {
				return err
			}
			asJSON, _ := cmd.Flags().GetBool("json")
			out := cmd.OutOrStdout()
			if asJSON {
				enc := json.NewEncoder(out)
				enc.SetIndent("", "  ")
				return enc.Encode(res)
			}
			for k, v := range res.Digests {
				fmt.Fprintf(out, "%s\t%s\n", k, v)
			}
			for _, d := range res.Diagnostics {
				fmt.Fprintf(cmd.ErrOrStderr(), "diagnostic: [%s] %s\n", d.Code, d.Message)
			}
			if hasErrorDiag(res.Diagnostics) {
				return errExitOne
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&input, "input", "i", "", "string to hash")
	cmd.Flags().StringSliceVarP(&algos, "algo", "a", []string{"sha256"}, "algorithms (md5,sha1,sha256,sha512)")
	cmd.Flags().StringVar(&encoding, "encoding", "hex", "encoding: hex|base64")
	return cmd
}
