package cmd

import (
	"github.com/spf13/cobra"
)

var adminsCmd = &cobra.Command{
	Use:   "admins",
	Short: "Administrator management",
}

func init() {
	rootCmd.AddCommand(adminsCmd)
}
