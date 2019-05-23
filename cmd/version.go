package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version contains the version string for this tool
var Version string

// NewCmdVersion returns new initialized instance of the version sub command
func NewCmdVersion() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Displays the ecom version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(os.Stdout, "%s\n", Version)
			os.Exit(0)
		},
	}
	return cmd
}
