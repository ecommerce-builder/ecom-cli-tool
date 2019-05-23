package assocs

import (
	"github.com/spf13/cobra"
)

// NewCmdAssocs returns new initialized instance of assocs sub command
func NewCmdAssocs() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "assocs",
		Short: "Associations management",
	}
	cmd.AddCommand(NewCmdAssocsApply())
	cmd.AddCommand(NewCmdAssocsList())
	cmd.AddCommand(NewCmdAssocsPurge())
	return cmd
}
