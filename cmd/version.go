package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version contains the version string for this tool
var Version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "prints ecom version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "%s\n", Version)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
