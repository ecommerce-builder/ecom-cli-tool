package carts

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdCartListProducts returns new initialized instance of the list sub command.
func NewCmdCartListProducts() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "list-products",
		Short: "list products in cart",
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
			cartProducts, err := client.GetCartProducts(ctx, cartID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			format := "%s\t%s\t%s\t%v\t%v\t%v\t%v\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "Cart Product ID", "SKU", "Name",
				"Qty", "Unit price", "Created", "Modified")
			fmt.Fprintf(tw, format, "--", "----", "---",
				"---", "----------", "-------", "--------")

			for _, v := range cartProducts {
				fmt.Fprintf(tw, format,
					v.ID, v.SKU, v.Name, v.Qty, v.UnitPrice,
					v.Created.In(location).Format(timeDisplayFormat),
					v.Modified.In(location).Format(timeDisplayFormat))
			}
			tw.Flush()
		},
	}

	return cmd
}
