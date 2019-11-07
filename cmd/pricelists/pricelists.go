package pricelists

import (
	"github.com/spf13/cobra"
)

// NewCmdPriceLists returns new initialized instance of pricelist sub command
func NewCmdPriceLists() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "pricelists",
		Short: "Price list management",
	}
	cmd.AddCommand(NewCmdPriceListsCreate())
	cmd.AddCommand(NewCmdPriceListsGet())
	cmd.AddCommand(NewCmdPriceListsList())
	cmd.AddCommand(NewCmdPriceListUpdate())
	cmd.AddCommand(NewCmdPriceListsDelete())
	return cmd
}
