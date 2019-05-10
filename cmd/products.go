package cmd

import (
	"github.com/spf13/cobra"
)

// projectsCmd represents the projects command
var productsCmd = &cobra.Command{
	Use:   "products",
	Short: "Products management",
}

func init() {
	rootCmd.AddCommand(productsCmd)
}
