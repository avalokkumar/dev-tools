package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func newDocsCmd(io Stdio) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Generate documentation artifacts",
	}
	cmd.AddCommand(newDocsGenCmd(io))
	return cmd
}

func newDocsGenCmd(io Stdio) *cobra.Command {
	var outDir string
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate Markdown docs for every CLI subcommand",
		Long: "Walks the cobra command tree and writes one Markdown file per " +
			"subcommand to --out. Used to publish CLI reference docs.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if outDir == "" {
				return fmt.Errorf("--out is required")
			}
			if err := os.MkdirAll(outDir, 0o755); err != nil {
				return err
			}
			// Walk a fresh root tree so the generated docs do not include the
			// `docs` subcommand itself recursing.
			root := newRootCmd(io)
			root.DisableAutoGenTag = true
			return doc.GenMarkdownTree(root, outDir)
		},
	}
	cmd.Flags().StringVarP(&outDir, "out", "o", "", "output directory")
	return cmd
}
