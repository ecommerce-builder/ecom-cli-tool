package prices

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

const timeDisplayFormat = "2006-01-02 15:04"

var location *time.Location

func init() {
	var err error
	location, err = time.LoadLocation("Europe/London")
	if err != nil {
		fmt.Fprintf(os.Stderr, "time.LoadLocation(%q) failed: %+v", "Europe/London", err.Error())
		return
	}
}

// NewCmdPricesList returns new initialized instance of the list sub command
func NewCmdPricesList() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "list",
		Short: "list all prices for all products",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			ctx := context.Background()
			prices, err := client.GetPrices(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			format := "%s\t%s\t%s\t%v\t£%.4f\t%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n",
				"Price ID", "Product SKU",
				"Price List Code", "Break", "Unit Price",
				"Created", "Modified")
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n",
				"--------", "-----------",
				"---------------", "-----", "----------",
				"-------", "--------")

			for _, p := range prices {
				fmt.Fprintf(tw, format,
					p.ID, p.ProductSKU,
					p.PriceListCode, p.Break,
					float64(p.UnitPrice)/10000.0,
					p.Created.In(location).Format(timeDisplayFormat),
					p.Modified.In(location).Format(timeDisplayFormat))
			}
			tw.Flush()
		},
	}
	return cmd
}
