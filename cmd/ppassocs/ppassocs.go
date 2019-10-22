package ppassocs

import (
	"github.com/spf13/cobra"
)

// NewCmdPPAssocs returns new initialized instance of the ppassocs sub command
func NewCmdPPAssocs() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "ppassocs",
		Short: "Product to product associations management",
	}
	cmd.AddCommand(NewCmdPPAssocsList())
	cmd.AddCommand(NewCmdPPAssocsDelete())
	return cmd
}
