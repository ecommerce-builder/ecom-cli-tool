package address

import (
	"github.com/spf13/cobra"
)

// NewCmdAddress returns new initialized instance of address sub command
func NewCmdAddress() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "address",
		Short: "Address Management",
	}
	cmd.AddCommand(NewCmdAddressCreate())
	cmd.AddCommand(NewCmdAddressGet())
	cmd.AddCommand(NewCmdAddressList())
	cmd.AddCommand(NewCmdAddressUpdate())
	cmd.AddCommand(NewCmdAddressDelete())
	return cmd
}
