package orders

import (
	"github.com/spf13/cobra"
)

// NewCmdOrders returns new initialized instance of orders sub command
func NewCmdOrders() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "orders",
		Short: "Orders management",
	}
	cmd.AddCommand(NewCmdOrdersCreate())
	cmd.AddCommand(NewCmdOrdersGet())
	cmd.AddCommand(NewCmdOrdersList())
	cmd.AddCommand(NewCmdOrdersStripeCheckout())
	return cmd
}
