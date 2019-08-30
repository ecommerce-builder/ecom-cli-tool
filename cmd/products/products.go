package products

import (
	"github.com/spf13/cobra"
)

// NewCmdProducts returns new initialized instance of products sub command
func NewCmdProducts() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "products",
		Short: "Products management",
	}
	cmd.AddCommand(NewCmdProductsApply())
	cmd.AddCommand(NewCmdProductsDelete())
	cmd.AddCommand(NewCmdProductsGet())
	cmd.AddCommand(NewCmdProductsList())
	return cmd
}
