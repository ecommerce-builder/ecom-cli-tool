package pricelists

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdPriceListsList returns new initialized instance of the get sub command
func NewCmdPriceListsList() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "list price lists",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			ctx := context.Background()
			priceLists, err := client.GetPriceLists(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			format := "%v\t%v\t%v\t%v\t%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format,
				"Price List Code",
				"Currency Code",
				"Strategy",
				"Inc Tax",
				"Name",
				"Description")
			fmt.Fprintf(tw, format,
				"---------------",
				"-------------",
				"--------",
				"-------",
				"----",
				"-----------")
			for _, v := range priceLists {
				fmt.Fprintf(tw, format,
					v.PriceListCode,
					v.CurrencyCode,
					v.Strategy,
					v.IncTax,
					v.Name,
					v.Description)
			}
			tw.Flush()
		},
	}
	return cmd
}
