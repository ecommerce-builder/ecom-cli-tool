package customers

import (
	"github.com/spf13/cobra"
)

// NewCmdCustomers returns new initialized instance of the customers sub command
func NewCmdCustomers() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "customers",
		Short: "Customer management",
	}
	cmd.AddCommand(NewCmdCustomersList())
	return cmd
}
