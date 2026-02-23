package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Set via ldflags
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("notion-copy %s\n", Version)
		fmt.Printf("  commit:  %s\n", Commit)
		fmt.Printf("  built:   %s\n", Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
