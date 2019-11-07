package pricelists

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// NewCmdPriceListUpdate returns new initialized instance of the update sub command
func NewCmdPriceListUpdate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "update <price_list_code>",
		Short: "Update a price list",
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
			if priceListID == "" {
				fmt.Fprintf(os.Stderr,
					"price list with code %q not found\n",
					priceListCode)
				os.Exit(1)
			}

			existing, err := client.GetPriceList(ctx, priceListID)
			if err == eclient.ErrPriceListNotFound {
				fmt.Fprintf(os.Stderr,
					"price list %s not found\n",
					priceListCode)
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			req, err := promptUpdatePriceList(existing)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			priceList, err := client.UpdatePriceList(ctx, priceListID, req)
			if err == eclient.ErrPriceListNotFound {
				fmt.Fprintf(os.Stderr,
					"price list %s not found\n",
					priceListCode)
				os.Exit(1)
			}
			if err == eclient.ErrPriceListCodeExists {
				fmt.Fprintf(os.Stderr,
					"price list %s is already exists\n",
					priceListCode)
				os.Exit(1)
			}
			if err != nil {
				fmt.Printf("%+v\n", err)
				fmt.Fprintf(os.Stderr, "error creating user: %+v\n", errors.Unwrap(err))
				os.Exit(1)
			}

			showPriceList(priceList)
		},
	}
	return cmd
}

func promptUpdatePriceList(existing *eclient.PriceList) (*eclient.UpdatePriceListRequest, error) {
	var req eclient.UpdatePriceListRequest

	// price_list_code
	u := &survey.Input{
		Message: "Price List Code:",
		Default: existing.PriceListCode,
	}
	survey.AskOne(u, &req.PriceListCode, nil)

	// currency_code
	c := &survey.Select{
		Message: "Currency Code:",
		Options: []string{
			"GBP",
			"EUR",
			"USD",
		},
		Default: existing.CurrencyCode,
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
		Default: existing.Strategy,
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
		Default: existing.Name,
	}
	survey.AskOne(n, &req.Name, survey.Required)

	// description
	d := &survey.Input{
		Message: "Description",
		Default: existing.Description,
	}
	survey.AskOne(d, &req.Description, survey.Required)

	return &req, nil
}
