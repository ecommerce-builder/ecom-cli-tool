package pricelists

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/ecommerce-builder/ecom-cli-tool/service"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// NewCmdPriceListsCreate returns new initialized instance of create sub command
func NewCmdPriceListsCreate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create a price list",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			req, err := promptCreatePriceList()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			// attempt to create the price list
			ctx := context.Background()
			priceList, err := client.CreatePriceList(ctx, req)
			if err != nil {
				fmt.Printf("%+v\n", err)
				os.Exit(1)
			}

			showPriceList(priceList)
		},
	}
	return cmd
}

func promptCreatePriceList() (*eclient.CreatePriceListRequest, error) {
	var req eclient.CreatePriceListRequest

	// price_list_code
	p := &survey.Input{
		Message: "Price List Code:",
	}
	survey.AskOne(p, &req.PriceListCode, survey.Required)

	// currency_code
	c := &survey.Select{
		Message: "Currency Code:",
		Options: []string{
			"GBP",
			"EUR",
			"USD",
		},
	}
	survey.AskOne(c, &req.CurrencyCode, survey.Required)

	// strategy
	s := &survey.Select{
		Message: "Strategy:",
		Options: []string{
			"simple",
			"volume",
			"tiered",
		},
	}
	survey.AskOne(s, &req.Strategy, survey.Required)

	// inc_tax
	i := &survey.Confirm{
		Message: "Inc Tax?",
		Default: false,
	}
	survey.AskOne(i, &req.IncTax, survey.Required)

	// name
	n := &survey.Input{
		Message: "Name:",
	}
	survey.AskOne(n, &req.Name, survey.Required)

	// description
	d := &survey.Input{
		Message: "Description",
	}
	survey.AskOne(d, &req.Description, survey.Required)

	return &req, nil
}

func showPriceList(v *eclient.PriceList) {
	format := "%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Price List ID:", v.ID)
	fmt.Fprintf(tw, format, "Price List Code:", v.PriceListCode)
	fmt.Fprintf(tw, format, "Currency Code:", v.CurrencyCode)
	fmt.Fprintf(tw, format, "Strategy:", v.Strategy)
	fmt.Fprintf(tw, format, "Inc Tax:", v.IncTax)
	fmt.Fprintf(tw, format, "Name:", v.Name)
	fmt.Fprintf(tw, format, "Description:", v.Description)
	fmt.Fprintf(tw, format, "Created:", v.Created.In(service.Location).Format(service.TimeDisplayFormat))
	fmt.Fprintf(tw, format, "Modified:", v.Modified.In(service.Location).Format(service.TimeDisplayFormat))
	tw.Flush()
}
