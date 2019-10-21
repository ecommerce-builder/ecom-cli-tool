package tariffs

import (
	"github.com/spf13/cobra"
)

// NewCmdShippingTariffsRules returns new initialized instance of tariffs sub command
func NewCmdShippingTariffsRules() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "tariffs",
		Short: "Shipping Tariff Management",
	}
	cmd.AddCommand(NewCmdShippingTariffsList())
	return cmd
}
