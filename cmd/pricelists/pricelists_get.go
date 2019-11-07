package pricelists

import (
	"context"
	"fmt"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdPriceListsGet returns new initialized instance of the get sub command
func NewCmdPriceListsGet() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "get <price_list_code>",
		Short: "Get price list",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			priceListCode := args[0]
			ctx := context.Background()
			priceLists, err := client.GetPriceLists(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			var priceListID string
			for _, v := range priceLists {
				if v.PriceListCode == priceListCode {
					priceListID = v.ID
					break
				}
			}

			priceList, err := client.GetPriceList(ctx, priceListID)
			if err == eclient.ErrProductNotFound {
				fmt.Printf("price list %s not found.\n", priceListCode)
				os.Exit(0)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			showPriceList(priceList)
		},
	}
	return cmd
}
