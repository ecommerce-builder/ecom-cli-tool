package orders

import (
	"context"
	"fmt"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdOrdersStripeCheckout returns new initialized instance of the stripecheckout sub command
func NewCmdOrdersStripeCheckout() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "stripecheckout <order_id>",
		Short: "Stripe checkout an order",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			orderID := args[0]
			if !cmdvalidate.IsValidUUID(orderID) {
				fmt.Fprintf(os.Stderr, "order_id %q is not a valid v4 uuid\n",
					orderID)
				os.Exit(1)
			}

			ctx := context.Background()
			sessionID, err := client.StripeCheckout(ctx, orderID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Stripe checkout id: %s\n", sessionID)
		},
	}
	return cmd
}
