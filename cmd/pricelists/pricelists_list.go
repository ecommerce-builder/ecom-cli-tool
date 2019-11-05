package pricelists

import (
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

			priceLists, err := client.GetPriceLists()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			format := "%v\t%v\t%v\t%v\t%v\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "Price List Code", "Currency Code", "Strategy", "Inc Tax", "Name")
			fmt.Fprintf(tw, format, "---------------", "-------------", "--------", "-------", "----")
			for _, l := range priceLists {
				fmt.Fprintf(tw, format, l.PriceListCode, l.CurrencyCode, l.Strategy, l.IncTax, l.Name)
			}
			tw.Flush()
		},
	}
	return cmd
}
