package cmd

import (
	"github.com/spf13/cobra"
)

// productsCmd represents the products command
var productsCmd = &cobra.Command{
	Use:   "products",
	Short: "Products management",
}

func init() {
	rootCmd.AddCommand(productsCmd)
}
