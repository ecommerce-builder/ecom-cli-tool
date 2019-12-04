package orders

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/ecommerce-builder/ecom-cli-tool/service"
	"github.com/spf13/cobra"
)

// NewCmdOrdersList returns new initialized instance of list sub command
func NewCmdOrdersList() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List orders",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			ctx := context.Background()
			orders, err := client.GetOrders(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			format := "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n",
				"ID",
				"Order ID",
				"Status",
				"Payment",
				"Contact name",
				"Email",
				"Currency",
				"Total ex VAT",
				"VAT Total",
				"Total inc VAT",
				"Created")
			fmt.Fprintf(tw, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t\n",
				"--------",
				"-----------",
				"------",
				"-------",
				"------------",
				"-----",
				"--------",
				"------------",
				"---------",
				"-------------",
				"-------")
			for _, v := range orders {
				fmt.Fprintf(tw, format,
					v.ID,
					v.OrderID,
					v.Status,
					v.Payment,
					v.User.ContactName,
					v.User.Email,
					v.Currency,
					service.IntPriceToString(v.TotalExVAT),
					service.IntPriceToString(v.VATTotal),
					service.IntPriceToString(v.TotalIncVAT),
					v.Created.In(service.Location).Format(service.TimeDisplayFormat))
			}
			tw.Flush()
		},
	}
	return cmd
}
