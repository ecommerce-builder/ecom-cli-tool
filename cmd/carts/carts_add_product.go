package carts

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdCartsAddProduct returns new initialized instance of the add-product sub command
func NewCmdCartsAddProduct() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "add-product <sku>",
		Short: "Add a product to an existing shopping cart",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}

			sku := args[0]
			ctx := context.Background()
			products, err := client.GetProducts(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			var productID string
			for _, v := range products {
				if v.SKU == sku {
					productID = v.ID
					break
				}
			}
			if productID == "" {
				fmt.Fprintf(os.Stderr, "product with sku %q not found\n", sku)
				os.Exit(1)
			}

			cartID := os.Getenv("ECOM_CLI_CART_ID")
			if !cmdvalidate.IsValidUUID(cartID) {
				fmt.Fprintf(os.Stderr, "ECOM_CLI_CART_ID value (%q) is not a valid v4 uuid\n", cartID)
				os.Exit(1)
			}

			cartProductRequest := eclient.CartProductRequest{
				CartID:    cartID,
				ProductID: productID,
				Qty:       1,
			}

			cartProduct, err := client.CartAddProduct(ctx, &cartProductRequest)
			if err == eclient.ErrCartNotFound {
				fmt.Fprintf(os.Stderr, "cart %q not found\n", cartID)
				os.Exit(1)
			}
			if err == eclient.ErrCartProductExists {
				fmt.Fprintf(os.Stderr, "cart product (sku=%q) already in the cart\n", sku)
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to add product id=%q to cart id=%s: %v\n",
					productID, cartID, err)
				os.Exit(1)
			}

			format := "%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "Cart Item ID:", cartProduct.ID)
			fmt.Fprintf(tw, format, "Product ID:", cartProduct.ProductID)
			fmt.Fprintf(tw, format, "SKU ID:", cartProduct.SKU)
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
