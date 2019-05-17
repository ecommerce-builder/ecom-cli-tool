package cmd

import (
	"github.com/spf13/cobra"
)

// projectsCmd represents the projects command
var projectsCmd = &cobra.Command{
	Use:   "configs",
	Short: "Configuration management",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(projectsCmd)
}
