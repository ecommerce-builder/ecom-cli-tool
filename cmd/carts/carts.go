package carts

import (
	"github.com/spf13/cobra"
)

// NewCmdCarts returns new initialized instance of carts sub command
func NewCmdCarts() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "carts",
		Short: "Carts management",
	}
	cmd.AddCommand(NewCmdCartsCreate())
	cmd.AddCommand(NewCmdCartListProducts())
	cmd.AddCommand(NewCmdCartsAddProduct())
	cmd.AddCommand(NewCmdCartUpdateProduct())
	cmd.AddCommand(NewCmdCartDeleteProduct())
	cmd.AddCommand(NewCmdCartEmptyProducts())
	return cmd
}
