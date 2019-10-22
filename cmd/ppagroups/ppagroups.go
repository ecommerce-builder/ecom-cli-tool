package ppagroups

import (
	"github.com/spf13/cobra"
)

// NewCmdPPAGroups returns new initialized instance of the ppagroups sub command
func NewCmdPPAGroups() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "ppagroups",
		Short: "Product to product associations groups Management",
	}
	cmd.AddCommand(NewCmdPPAGroupCreate())
	cmd.AddCommand(NewCmdPPAGroupsList())
	cmd.AddCommand(NewCmdPPAGroupsGet())
	cmd.AddCommand(NewCmdPPAGroupsDelete())
	return cmd
}
