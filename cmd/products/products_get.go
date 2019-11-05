package products

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdProductsGet returns new initialized instance of the get sub command
func NewCmdProductsGet() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "get <sku>",
		Short: "Get product",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
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

			product, err := client.GetProduct(ctx, productID)
			if err == eclient.ErrProductNotFound {
				fmt.Printf("Product %s not found.\n", sku)
				os.Exit(0)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			format := "%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "Product ID:", product.ID)
			fmt.Fprintf(tw, format, "Path:", product.Path)
			fmt.Fprintf(tw, format, "SKU:", product.SKU)
			fmt.Fprintf(tw, format, "Name:", product.Name)
			fmt.Fprintf(tw, format, "Created:",
				product.Created.In(location).Format(timeDisplayFormat))
			fmt.Fprintf(tw, format, "Modified:",
				product.Modified.In(location).Format(timeDisplayFormat))
			tw.Flush()
		},
	}
	return cmd
}
