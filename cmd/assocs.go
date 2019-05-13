package cmd

import (
	"github.com/spf13/cobra"
)

// catalogCmd represents the catalog command
var assocsCmd = &cobra.Command{
	Use:   "assocs",
	Short: "Associations management",
}

func init() {
	rootCmd.AddCommand(assocsCmd)
}
