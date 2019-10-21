package users

import (
	"github.com/spf13/cobra"
)

// NewCmdUsers returns new initialized instance of the customers sub command
func NewCmdUsers() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "users",
		Short: "User management",
	}
	cmd.AddCommand(NewCmdUsersCreate())
	cmd.AddCommand(NewCmdUsersList())
	return cmd
}
