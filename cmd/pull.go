package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pullRecursive bool
var pullIncludeDatabases bool
var pullFlatten bool

var pullCmd = &cobra.Command{
	Use:   "pull <page-id> <output-dir>",
	Short: "Export Notion pages to markdown files",
	Long: `Pull downloads a Notion page (and optionally its children) and converts
them to markdown files on your local filesystem.

Example:
  notion-copy pull abc123def456 ./output/
  notion-copy pull abc123def456 ./backup/ --recursive`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		pageID := args[0]
		outputDir := args[1]

		fmt.Printf("Pulling page %s to %s\n", pageID, outputDir)
		fmt.Printf("  Recursive: %v\n", pullRecursive)
		fmt.Printf("  Include databases: %v\n", pullIncludeDatabases)
		fmt.Printf("  Flatten: %v\n", pullFlatten)

		// TODO: Implement pull logic
		// 1. Initialize Notion client
		// 2. Fetch page and blocks
		// 3. Convert blocks to markdown
		// 4. Write to output directory
		// 5. If recursive, fetch child pages

		return fmt.Errorf("not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	pullCmd.Flags().BoolVarP(&pullRecursive, "recursive", "r", false, "Include nested subpages")
	pullCmd.Flags().BoolVar(&pullIncludeDatabases, "include-databases", false, "Export linked databases as CSV")
	pullCmd.Flags().BoolVar(&pullFlatten, "flatten", false, "No subdirectories, flat file structure")
}
