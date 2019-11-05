package carts

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdCartUpdateProduct returns new initialized instance of the update-product sub command
func NewCmdCartUpdateProduct() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cartProductID string
	var qty int
	var cmd = &cobra.Command{
		Use:   "update-product <cart_product_id> <qty>",
		Short: "Update a product qty already in a cart",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("requires <cart_product_id> and <qty> arguments")
			}

			cartProductID = args[0]
			if !cmdvalidate.IsValidUUID(cartProductID) {
				return fmt.Errorf("cart_product_id value %q is not a valid v4 uuid", cartProductID)
			}

			qty, err = strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("<qty> must be an integer value: %v", err)
			}

			return nil
		},
		// cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			ctx := context.Background()
			cartProduct, err := client.UpdateCartProduct(ctx, cartProductID, qty)
			if err == eclient.ErrCartProductNotFound {
				fmt.Fprintf(os.Stderr, "cart product %q not found\n", cartProductID)
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			format := "%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "Product cart ID:", cartProduct.ID)
			fmt.Fprintf(tw, format, "SKU:", cartProduct.SKU)
			fmt.Fprintf(tw, format, "Product ID", cartProduct.ProductID)
			fmt.Fprintf(tw, format, "Name:", cartProduct.Name)
			fmt.Fprintf(tw, format, "Qty:", cartProduct.Qty)
			fmt.Fprintf(tw, format, "Unit price:", cartProduct.UnitPrice)
			fmt.Fprintf(tw, format, "Created:",
				cartProduct.Created.In(location).Format(timeDisplayFormat))
			fmt.Fprintf(tw, format, "Modified:",
				cartProduct.Modified.In(location).Format(timeDisplayFormat))
			tw.Flush()
		},
	}
	return cmd
}
