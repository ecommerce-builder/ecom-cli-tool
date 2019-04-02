package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// customersCmd represents the customers command
var customersCmd = &cobra.Command{
	Use:   "customers",
	Short: "Customer management",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("customers called")
	},
}

func init() {
	rootCmd.AddCommand(customersCmd)
}
