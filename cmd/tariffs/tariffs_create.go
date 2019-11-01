package tariffs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

var countryCodes []string

func init() {
	countryCodes = []string{
		"UK - United Kingdom",
		"AT - Austria",
		"BG - Bulgaria",
		"CZ - Czechia",
		"DK - Denmark",
		"EE - Estonia",
		"FI - Finland",
		"FR - France",
		"DE - Germany",
		"HU - Hungary",
		"IT - Italy",
		"IE - Ireland",
		"LT - Lithuania",
		"LU - Luxembourg",
		"NL - Netherlands",
		"NO - Norway",
		"PL - Poland",
		"PT - Portugal",
		"RO - Romania",
		"SK - Slovakia",
		"SI - Slovenia",
		"ES - Spain (including Balearic islands)",
		"SE - Sweden",
		"CH - Switzerland",
		"US - United States",
	}
}

// NewCmdShippingTarrifsCreate returns new initialized instance of create sub command
func NewCmdShippingTarrifsCreate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new shipping tariff",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}

			// get the request params
			req, err := promptCreateShippingTariff()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			ctx := context.Background()
			tariff, err := client.CreateShippingTariff(ctx, req)
			if err != nil {
				fmt.Printf("%+v\n", err)
				os.Exit(1)
			}

			format := "%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "Shipping Tariff ID:", tariff.ID)
			fmt.Fprintf(tw, format, "Country Code:", tariff.CountryCode)
			fmt.Fprintf(tw, format, "Shipping Code:", tariff.ShippingCode)
			fmt.Fprintf(tw, format, "Name:", tariff.Name)
			fmt.Fprintf(tw, format, "Price:", tariff.Price)
			fmt.Fprintf(tw, format, "Tax Code:", tariff.TaxCode)
			fmt.Fprintf(tw, format, "Created:", tariff.Created.In(location).Format(timeDisplayFormat))
			fmt.Fprintf(tw, format, "Modified:", tariff.Modified.In(location).Format(timeDisplayFormat))
			tw.Flush()
		},
	}
	return cmd
}

func promptCreateShippingTariff() (*eclient.CreateShippingTariffRequest, error) {
	var req eclient.CreateShippingTariffRequest

	// country_code
	var countryCode string
	c := &survey.Select{
		Message: "Shipping Code:",
		Options: countryCodes,
	}
	survey.AskOne(c, &countryCode, survey.Required)
	req.CountryCode = countryCode[0:2]

	// shipping_code
	s := &survey.Input{
		Message: "Shipping Code:",
	}
	survey.AskOne(s, &req.Shippingcode, survey.Required)

	// name
	n := &survey.Input{
		Message: "Name:",
	}
	survey.AskOne(n, &req.Name, survey.Required)

	// price
	p := &survey.Input{
		Message: "Price:",
	}
	survey.AskOne(p, &req.Price, survey.ComposeValidators(
		survey.Required,
		func(val interface{}) error {
			str, ok := val.(string)
			if !ok {
				return errors.New("invalid response")
			}

			v, err := strconv.Atoi(str)
			if err != nil {
				return err
			}

			if v < 0 {
				return errors.New("price must be a positive integer")
			}
			return nil
		},
	))

	// tax_code
	t := &survey.Select{
		Message: "Tax Code:",
		Options: []string{
			"T20",
			"T0",
		},
	}
	survey.AskOne(t, &req.TaxCode, survey.Required)

	return &req, nil
}
