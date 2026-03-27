package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jomei/notionapi"
	"github.com/kris-hansen/notion-copy/internal/convert"
	"github.com/kris-hansen/notion-copy/internal/notion"
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

		// Initialize Notion client
		client, err := notion.New()
		if err != nil {
			return fmt.Errorf("failed to create Notion client: %w", err)
		}

		ctx := context.Background()

		// Create output directory
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		return pullPage(ctx, client, pageID, outputDir)
	},
}

func pullPage(ctx context.Context, client *notion.Client, pageID string, outputDir string) error {
	// Get page info
	page, err := client.GetPage(ctx, pageID)
	if err != nil {
		return fmt.Errorf("failed to get page %s: %w", pageID, err)
	}

	// Extract title from page properties
	title := extractPageTitle(page)
	if title == "" {
		title = "Untitled"
	}

	// Get all blocks for this page
	blocks, err := client.GetBlocks(ctx, pageID)
	if err != nil {
		return fmt.Errorf("failed to get blocks for %s: %w", pageID, err)
	}

	// Separate child pages from content blocks
	var contentBlocks []notionapi.Block
	var childPages []childPageInfo

	for _, block := range blocks {
		if cpb, ok := block.(*notionapi.ChildPageBlock); ok {
			childPages = append(childPages, childPageInfo{
				ID:    string(block.GetID()),
				Title: cpb.ChildPage.Title,
			})
		} else {
			contentBlocks = append(contentBlocks, block)
		}
	}

	// Convert blocks to markdown
	markdown, err := convert.ToMarkdown(contentBlocks)
	if err != nil {
		return fmt.Errorf("failed to convert blocks: %w", err)
	}

	// Add title as H1 if not already present
	if !strings.HasPrefix(strings.TrimSpace(markdown), "# ") {
		markdown = "# " + title + "\n\n" + markdown
	}

	// Write to file
	filename := sanitizeFilename(title) + ".md"
	filePath := filepath.Join(outputDir, filename)

	if err := os.WriteFile(filePath, []byte(markdown), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filePath, err)
	}

	fmt.Printf("📄 Pulled: %s → %s\n", title, filePath)

	// Handle child pages if recursive
	if pullRecursive && len(childPages) > 0 {
		var childDir string
		if pullFlatten {
			childDir = outputDir
		} else {
			// Create subdirectory for children
			childDir = filepath.Join(outputDir, sanitizeFilename(title))
			if err := os.MkdirAll(childDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", childDir, err)
			}
		}

		for _, child := range childPages {
			if err := pullPage(ctx, client, child.ID, childDir); err != nil {
				// Log error but continue with other pages
				fmt.Fprintf(os.Stderr, "⚠️  Failed to pull %s: %v\n", child.Title, err)
			}
		}
	}

	return nil
}

type childPageInfo struct {
	ID    string
	Title string
}

// extractPageTitle gets the title from a Notion page's properties
func extractPageTitle(page *notionapi.Page) string {
	// Try "title" property (standard for pages)
	if titleProp, ok := page.Properties["title"]; ok {
		if tp, ok := titleProp.(*notionapi.TitleProperty); ok {
			var parts []string
			for _, rt := range tp.Title {
				parts = append(parts, rt.PlainText)
			}
			return strings.Join(parts, "")
		}
	}

	// Try "Name" property (common in databases)
	if nameProp, ok := page.Properties["Name"]; ok {
		if tp, ok := nameProp.(*notionapi.TitleProperty); ok {
			var parts []string
			for _, rt := range tp.Title {
				parts = append(parts, rt.PlainText)
			}
			return strings.Join(parts, "")
		}
	}

	// Iterate all properties looking for title type
	for _, prop := range page.Properties {
		if tp, ok := prop.(*notionapi.TitleProperty); ok {
			var parts []string
			for _, rt := range tp.Title {
				parts = append(parts, rt.PlainText)
			}
			if len(parts) > 0 {
				return strings.Join(parts, "")
			}
		}
	}

	return ""
}

// sanitizeFilename removes characters that aren't safe for filenames
func sanitizeFilename(name string) string {
	// Replace problematic characters
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "",
	)
	name = replacer.Replace(name)

	// Trim spaces and dots from ends
	name = strings.TrimSpace(name)
	name = strings.Trim(name, ".")

	// Limit length
	if len(name) > 100 {
		name = name[:100]
	}

	if name == "" {
		name = "untitled"
	}

	return name
}

func init() {
	rootCmd.AddCommand(pullCmd)

	pullCmd.Flags().BoolVarP(&pullRecursive, "recursive", "r", false, "Include nested subpages")
	pullCmd.Flags().BoolVar(&pullIncludeDatabases, "include-databases", false, "Export linked databases as CSV")
	pullCmd.Flags().BoolVar(&pullFlatten, "flatten", false, "No subdirectories, flat file structure")
}
