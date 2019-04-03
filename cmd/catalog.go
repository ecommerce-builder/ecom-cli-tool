package cmd

import (
	"github.com/spf13/cobra"
)

// catalogCmd represents the catalog command
var catalogCmd = &cobra.Command{
	Use:   "catalog",
	Short: "Catalog management",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(catalogCmd)
}
