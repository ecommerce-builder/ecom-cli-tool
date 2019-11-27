package orders

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/ecommerce-builder/ecom-cli-tool/service"
	"github.com/spf13/cobra"
)

// NewCmdOrdersGet returns new initialized instance of the get sub command
func NewCmdOrdersGet() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "get <offer_id>",
		Short: "Get an order by id",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			orderID := args[0]
			if !cmdvalidate.IsValidUUID(orderID) {
				fmt.Fprintf(os.Stderr, "order_id %q is not a valid v4 uuid\n",
					orderID)
				os.Exit(1)
			}

			ctx := context.Background()
			order, err := client.GetOrder(ctx, orderID)
			if err == eclient.ErrOrderNotFound {
				fmt.Fprintf(os.Stderr, "order %q not found\n", orderID)
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			showOrder(order)
		},
	}
	return cmd
}

func showOrder(v *eclient.Order) {
	format := "%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "ID:", v.ID)
	fmt.Fprintf(tw, format, "Order ID:", v.OrderID)
	fmt.Fprintf(tw, format, "Status:", v.Status)
	fmt.Fprintf(tw, format, "Payment:", v.Payment)
	fmt.Fprintf(tw, format, "Contact name:", v.ContactName)
	fmt.Fprintf(tw, format, "Email:", v.Email)
	fmt.Fprintf(tw, format, "Currency:", v.Currency)
	fmt.Fprintf(tw, format, "Total ex VAT:", service.IntPriceToString(v.TotalExVAT))
	fmt.Fprintf(tw, format, "VAT Total:", service.IntPriceToString(v.VATTotal))
	fmt.Fprintf(tw, format, "Total inc VAT:", service.IntPriceToString(v.TotalIncVAT))
	fmt.Fprintf(tw, format, "Created:", v.Created.In(service.Location).Format(service.TimeDisplayFormat))
	fmt.Fprintf(tw, format, "Modified:", v.Modified.In(service.Location).Format(service.TimeDisplayFormat))

	fmt.Fprintln(tw)
	fmt.Fprintf(tw, format, "Billing address", "")
	fmt.Fprintf(tw, format, "---------------", "")
	showAddr(tw, &v.Billing)

	fmt.Fprintln(tw)
	fmt.Fprintf(tw, format, "Shipping address", "")
	fmt.Fprintf(tw, format, "----------------", "")
	showAddr(tw, &v.Shipping)

	fmt.Fprintln(tw)
	fmt.Fprintf(tw, format, "Order Items", "")
	fmt.Fprintf(tw, format, "-----------", "")

	fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n",
		"Qty",
		"SKU",
		"Currency",
		"Unit price",
		"Tax Code",
		"VAT",
		"Total inc VAT")
	fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n",
		"---",
		"---",
		"--------",
		"----------",
		"--------",
		"---",
		"------------")
	for _, v := range v.Items {
		fmt.Fprintf(tw, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t\n",
			v.Qty,
			v.SKU,
			v.Currency,
			service.IntPriceToString(v.UnitPrice),
			v.TaxCode,
			service.IntPriceToString(v.VAT),
			service.IntPriceToString(v.VAT+v.UnitPrice))
	}

	fmt.Fprintf(tw, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t\n",
		"",
		"",
		"",
		service.IntPriceToString(v.TotalExVAT),
		"",
		service.IntPriceToString(v.VATTotal),
		service.IntPriceToString(v.TotalIncVAT))
	tw.Flush()
}

func showAddr(tw *tabwriter.Writer, v *eclient.OrderAddr) {
	format := "%v\t%v\t\n"
	fmt.Fprintf(tw, format, "Contact name:", v.ContactName)
	fmt.Fprintf(tw, format, "Address 1:", v.Addr1)
	var addr2 string
	if v.Addr2 != nil {
		addr2 = *v.Addr2
	}
	fmt.Fprintf(tw, format, "Address 2:", addr2)
	fmt.Fprintf(tw, format, "City:", v.City)
	var county string
	if v.County != nil {
		county = *v.County
	}
	fmt.Fprintf(tw, format, "County:", county)
	fmt.Fprintf(tw, format, "Postcode:", v.Postcode)
	fmt.Fprintf(tw, format, "Country code:", v.CountryCode)
}
