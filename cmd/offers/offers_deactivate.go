package offers

import (
	"context"
	"fmt"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdOffersDeactivate returns new initialized instance of the deactivate sub command
func NewCmdOffersDeactivate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "deactivate <offer_id>",
		Short: "Deactive an offer",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			// offer_id
			offerID := args[0]
			if !cmdvalidate.IsValidUUID(offerID) {
				fmt.Fprintf(os.Stderr, "offer_id %q is not a valid v4 uuid\n",
					offerID)
				os.Exit(1)
			}

			ctx := context.Background()
			err = client.DeleteOffer(ctx, offerID)
			if err == eclient.ErrOfferNotFound {
				fmt.Fprintf(os.Stderr,
					"offer not found. Use ecom offers list to check.\n")
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
		},
	}
	return cmd
}
