package prices

import (
	"github.com/spf13/cobra"
)

// NewCmdPrices returns new initialized instance of prices sub command
func NewCmdPrices() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "prices",
		Short: "Prices Management",
	}
	cmd.AddCommand(NewCmdPricesList())
	return cmd
}
