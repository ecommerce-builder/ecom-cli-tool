package cmd

import (
	"github.com/spf13/cobra"
)

var assocsCmd = &cobra.Command{
	Use:   "assocs",
	Short: "Associations management",
}

func init() {
	rootCmd.AddCommand(assocsCmd)
}
