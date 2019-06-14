package token

import (
	"github.com/spf13/cobra"
)

// NewCmdToken returns new initialized instance of token sub command
func NewCmdToken() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "token",
		Short: "Token management",
	}
	cmd.AddCommand(NewCmdTokenShow())
	return cmd
}
