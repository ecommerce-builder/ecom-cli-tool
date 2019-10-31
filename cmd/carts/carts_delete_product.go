package carts

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdCartDeleteProduct returns new initialized instance of the delete-product sub command
func NewCmdCartDeleteProduct() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "delete-product <cart_product_id>",
		Short: "Remove a product from a cart",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}

			cartProductID := args[0]
			if !cmdvalidate.IsValidUUID(cartProductID) {
				fmt.Fprintf(os.Stderr, "cart_product_id value (%q) is not a valid v4 uuid\n", cartProductID)
				os.Exit(1)
			}

			ctx := context.Background()
			err = client.CartsRemoveProduct(ctx, cartProductID)
			if err == eclient.ErrCartProductNotFound {
				fmt.Fprintf(os.Stderr, "cart product %q not found\n", cartProductID)
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
