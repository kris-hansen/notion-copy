package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "notion-copy",
	Short: "Copy Notion documents to the filesystem and back",
	Long: `notion-copy is a CLI tool for moving content between Notion and your local filesystem.

Pull: Notion pages → Markdown files
Push: Markdown files → Notion pages

No sync, no magic — just straightforward copy operations.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Try loading .env from current directory
	godotenv.Load()

	// Try loading from ~/.config/notion-copy/.env
	home, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(home, ".config", "notion-copy", ".env")
		godotenv.Load(configPath)
	}
}
