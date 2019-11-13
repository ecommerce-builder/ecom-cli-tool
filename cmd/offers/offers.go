package offers

import (
	"github.com/spf13/cobra"
)

// NewCmdOffers returns new initialized instance of offers sub command
func NewCmdOffers() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "offers",
		Short: "Offers management",
	}
	cmd.AddCommand(NewCmdOffersActivate())
	cmd.AddCommand(NewCmdOffersList())
	cmd.AddCommand(NewCmdOffersGet())
	cmd.AddCommand(NewCmdOffersDeactivate())
	return cmd
}
