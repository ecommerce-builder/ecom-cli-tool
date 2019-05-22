package cmd

import (
	"github.com/spf13/cobra"
)

var productsCmd = &cobra.Command{
	Use:   "products",
	Short: "Products management",
}

func init() {
	rootCmd.AddCommand(productsCmd)
}
