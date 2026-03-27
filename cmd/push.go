package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kris-hansen/notion-copy/internal/convert"
	"github.com/kris-hansen/notion-copy/internal/notion"
	"github.com/spf13/cobra"
)

var pushRecursive bool
var pushDryRun bool

var pushCmd = &cobra.Command{
	Use:   "push <input-path> [parent-page-id]",
	Short: "Import markdown files to Notion",
	Long: `Push uploads markdown files from your local filesystem and creates
corresponding pages in Notion under the specified parent page.

If no parent-page-id is provided, uses NOTION_DEFAULT_PAGE_ID from environment.

Supports both single files and directories.

Example:
  notion-copy push ./doc.md                      # Uses default page
  notion-copy push ./doc.md abc123def456         # Explicit parent
  notion-copy push ./docs/ -r                    # Recursive, default page
  notion-copy push ./docs/ abc123def456 --dry-run`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]

		// Get parent page ID from args or environment
		var parentPageID string
		if len(args) >= 2 {
			parentPageID = args[1]
		} else {
			parentPageID = os.Getenv("NOTION_DEFAULT_PAGE_ID")
			if parentPageID == "" {
				return fmt.Errorf("no parent page specified and NOTION_DEFAULT_PAGE_ID not set")
			}
		}

		// Check if input exists
		info, err := os.Stat(inputPath)
		if err != nil {
			return fmt.Errorf("cannot access %s: %w", inputPath, err)
		}

		// Initialize Notion client
		var client *notion.Client
		if !pushDryRun {
			client, err = notion.New()
			if err != nil {
				return fmt.Errorf("failed to create Notion client: %w", err)
			}
		}

		ctx := context.Background()

		if info.IsDir() {
			return pushDirectory(ctx, client, inputPath, parentPageID)
		}
		return pushFile(ctx, client, inputPath, parentPageID)
	},
}

func pushDirectory(ctx context.Context, client *notion.Client, dir string, parentPageID string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			if pushRecursive {
				// Create a page for the subdirectory, then push its contents
				fmt.Printf("📁 [dir] %s\n", entry.Name())
				if pushDryRun {
					// In dry run, just recurse with same parent
					if err := pushDirectory(ctx, client, path, parentPageID); err != nil {
						return err
					}
				} else {
					// Create a page for the directory
					page, err := client.CreatePage(ctx, parentPageID, entry.Name(), nil)
					if err != nil {
						return fmt.Errorf("failed to create page for %s: %w", entry.Name(), err)
					}
					// Recurse with new page as parent
					if err := pushDirectory(ctx, client, path, string(page.ID)); err != nil {
						return err
					}
				}
			}
			continue
		}

		// Only process markdown files
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			continue
		}

		if err := pushFile(ctx, client, path, parentPageID); err != nil {
			return err
		}
	}

	return nil
}

func pushFile(ctx context.Context, client *notion.Client, filePath string, parentPageID string) error {
	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	markdown := string(content)

	// Extract title from markdown or use filename
	title := convert.ExtractTitle(markdown)
	if title == "" {
		title = strings.TrimSuffix(filepath.Base(filePath), ".md")
	}

	// Convert markdown to Notion blocks
	blocks, err := convert.ToNotionBlocks(markdown)
	if err != nil {
		return fmt.Errorf("failed to convert %s: %w", filePath, err)
	}

	if pushDryRun {
		fmt.Printf("📄 [dry-run] %s → \"%s\" (%d blocks)\n", filePath, title, len(blocks))
		return nil
	}

	// Create the page
	page, err := client.CreatePage(ctx, parentPageID, title, blocks)
	if err != nil {
		return fmt.Errorf("failed to create page for %s: %w", filePath, err)
	}

	fmt.Printf("📄 Created: %s → \"%s\" (id: %s)\n", filePath, title, page.ID)
	return nil
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.Flags().BoolVarP(&pushRecursive, "recursive", "r", false, "Include subdirectories as nested pages")
	pushCmd.Flags().BoolVar(&pushDryRun, "dry-run", false, "Preview what would be created")
}
