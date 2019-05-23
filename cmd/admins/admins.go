package admins

import (
	"github.com/spf13/cobra"
)

// NewCmdAdmins returns new initialized instance of admins sub command
func NewCmdAdmins() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "admins",
		Short: "Administrator management",
	}
	cmd.AddCommand(NewCmdAdminsCreate())
	cmd.AddCommand(NewCmdAdminsList())
	cmd.AddCommand(NewCmdAdminsRemove())
	return cmd
}
