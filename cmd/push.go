package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pushRecursive bool
var pushDryRun bool

var pushCmd = &cobra.Command{
	Use:   "push <input-dir> <parent-page-id>",
	Short: "Import markdown files to Notion",
	Long: `Push uploads markdown files from your local filesystem and creates
corresponding pages in Notion under the specified parent page.

Example:
  notion-copy push ./docs/ abc123def456
  notion-copy push ./docs/ abc123def456 --recursive --dry-run`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputDir := args[0]
		parentPageID := args[1]

		fmt.Printf("Pushing %s to page %s\n", inputDir, parentPageID)
		fmt.Printf("  Recursive: %v\n", pushRecursive)
		fmt.Printf("  Dry run: %v\n", pushDryRun)

		// TODO: Implement push logic
		// 1. Initialize Notion client
		// 2. Scan input directory for markdown files
		// 3. Parse markdown to AST
		// 4. Convert AST to Notion blocks
		// 5. Create pages in Notion
		// 6. If recursive, handle subdirectories

		return fmt.Errorf("not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.Flags().BoolVarP(&pushRecursive, "recursive", "r", false, "Include subdirectories as nested pages")
	pushCmd.Flags().BoolVar(&pushDryRun, "dry-run", false, "Preview what would be created")
}
