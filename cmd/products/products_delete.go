package products

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdProductsDelete returns new initialized instance of the delete sub command
func NewCmdProductsDelete() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "delete <sku>",
		Short: "Delete product",
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

			err = client.DeleteProduct(ctx, productID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
		},
	}
	return cmd
}
