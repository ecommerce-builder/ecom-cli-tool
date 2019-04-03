package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var catalogGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the catalog hierarchy",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get catalog")
	},
}

func init() {
	catalogCmd.AddCommand(catalogGetCmd)
}
