package address

import (
	"context"
	"fmt"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/ecommerce-builder/ecom-cli-tool/service"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// NewCmdAddressUpdate returns new initialized instance of update sub command
func NewCmdAddressUpdate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "update <address_id>",
		Short: "Update an address",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			// address id
			addrID := args[0]
			if !cmdvalidate.IsValidUUID(addrID) {
				fmt.Fprintf(os.Stderr, "address_id value %q is not a valid v4 uuid\n", addrID)
				os.Exit(1)
			}

			ctx := context.Background()
			addr, err := client.GetAddress(ctx, addrID)
			if err == eclient.ErrAddressNotFound {
				fmt.Fprintf(os.Stderr, "address id %s not found\n", addrID)
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			req, err := promptUpdateAddress(addr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			if req.Type == nil && req.ContactName == nil &&
				req.Addr1 == nil && req.Addr2 == nil &&
				req.City == nil && req.County == nil &&
				req.Postcode == nil &&
				req.CountryCode == nil {
				os.Exit(0)
			}

			updated, err := client.UpdateAddress(ctx, addrID, req)
			if err != nil {
				fmt.Printf("%+v\n", err)
				fmt.Fprintf(os.Stderr, "error updating address: %+v", err)
			}
			showAddress(updated)
		},
	}
	return cmd
}

// only sets the req field if the value has changed, otherwise remained nil.
func promptUpdateAddress(c *eclient.Address) (*eclient.UpdateAddressRequest, error) {
	var req eclient.UpdateAddressRequest

	// type
	var typ string
	t := &survey.Select{
		Message: "Type:",
		Default: c.Typ,
		Options: []string{"shipping", "billing"},
	}
	survey.AskOne(t, &typ, survey.Required)
	if c.Typ != typ {
		req.Type = &typ
	}

	// contact_name
	var contactName string
	n := &survey.Input{
		Message: "Contact Name:",
		Default: c.ContactName,
	}
	survey.AskOne(n, &contactName, survey.Required)
	if c.ContactName != contactName {
		req.ContactName = &contactName
	}

	// addr1
	var addr1 string
	survey.AskOne(
		&survey.Input{
			Message: "Address 1:",
			Default: c.Addr1,
		},
		&addr1,
		survey.Required)
	if c.Addr1 != addr1 {
		req.Addr1 = &addr1
	}

	// addr2
	var addr2 string
	var currentAddr2 string
	if c.Addr2 != nil {
		currentAddr2 = *c.Addr2
	}
	survey.AskOne(
		&survey.Input{
			Message: "Address 2:",
			Default: currentAddr2,
		},
		&addr2,
		nil)

	if addr2 != "" && c.Addr2 != nil {
		if *c.Addr2 != addr2 {
			req.Addr2 = &addr2
		}
	}
	var city string
	survey.AskOne(
		&survey.Input{
			Message: "City:",
			Default: c.City,
		},
		&city,
		survey.Required)
	if c.City != city {
		req.City = &city
	}

	// county
	var county string
	var currentCounty string
	if c.County != nil {
		currentCounty = *c.County
	}
	survey.AskOne(
		&survey.Input{
			Message: "County:",
			Default: currentCounty,
		},
		&county,
		nil)

	if c.County != nil && *c.County != county {
		req.County = &county
	}

	// postcode
	var postcode string
	survey.AskOne(
		&survey.Input{
			Message: "Postcode:",
			Default: c.Postcode,
		}, &postcode, survey.Required)
	if c.Postcode != postcode {
		req.Postcode = &postcode
	}

	// country_code
	var countryCode string
	cy := &survey.Select{
		Message: "Country Code:",
		Default: service.TitleFromCountryCode(c.CountryCode),
		Options: service.CountryCodes(),
	}
	survey.AskOne(cy, &countryCode, survey.Required)
	countryCode = countryCode[0:2]
	if c.CountryCode != countryCode {
		req.CountryCode = &countryCode
	}

	return &req, nil
}
