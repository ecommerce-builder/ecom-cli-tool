package cmd

import (
	"github.com/spf13/cobra"
)

// profilesCmd represents the profiles command
var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Profile management",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(profilesCmd)
}
