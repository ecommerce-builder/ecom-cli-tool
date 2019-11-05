package carts

import (
	"context"
	"fmt"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdCartEmptyProducts returns new initialized instance of the empty sub command.
func NewCmdCartEmptyProducts() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "empty-products",
		Short: "empty all products from the cart",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			cartID := os.Getenv("ECOM_CLI_CART_ID")
			if !cmdvalidate.IsValidUUID(cartID) {
				fmt.Fprintf(os.Stderr, "ECOM_CLI_CART_ID value (%q) is not a valid v4 uuid\n", cartID)
				os.Exit(1)
			}

			ctx := context.Background()
			err := client.EmptyCartProducts(ctx, cartID)
			if err == eclient.ErrCartNotFound {
				fmt.Fprintf(os.Stderr, "cart %q not found. Set the environment variable ECOM_CLI_CART_ID to a valid v4 uuid.\n", cartID)
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
