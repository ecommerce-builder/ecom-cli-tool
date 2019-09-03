package pricelists

import (
	"github.com/spf13/cobra"
)

// NewCmdPriceLists returns new initialized instance of pricelist sub command
func NewCmdPriceLists() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "price-lists",
		Short: "Price list management",
	}
	cmd.AddCommand(NewCmdPriceListsList())
	return cmd
}
