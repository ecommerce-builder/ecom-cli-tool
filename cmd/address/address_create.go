package address

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/ecommerce-builder/ecom-cli-tool/service"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// NewCmdAddressCreate returns new initialized instance of create sub command
func NewCmdAddressCreate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "create <email>",
		Short: "Create a new address for a given user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			email := args[0]
			ctx := context.Background()
			users, err := client.GetUsers(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			userMap := make(map[string]string, 0)
			var userOpts []string
			for _, user := range users {
				userOpts = append(userOpts, user.Email)
				userMap[user.Email] = user.ID
			}

			if _, ok := userMap[email]; !ok {
				fmt.Fprintf(os.Stderr, "Email %s did not match any users\n", email)
				os.Exit(1)
			}

			req, err := promptCreateAddress()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			req.UserID = userMap[email]

			addr, err := client.CreateAddress(ctx, req)
			if err != nil {
				fmt.Printf("%+v\n", err)
				fmt.Fprintf(os.Stderr, "error creating address: %+v", errors.Unwrap(err))
			}
			showAddress(addr)
		},
	}
	return cmd
}

func promptCreateAddress() (*eclient.CreateAddressRequest, error) {
	var req eclient.CreateAddressRequest

	// type
	t := &survey.Select{
		Message: "Type:",
		Options: []string{"shipping", "billing"},
	}
	survey.AskOne(t, &req.Type, survey.Required)

	// contact_name
	n := &survey.Input{
		Message: "Contact Name:",
	}
	survey.AskOne(n, &req.ContactName, survey.Required)

	// addr1
	survey.AskOne(
		&survey.Input{
			Message: "Address 1:",
		},
		&req.Addr1,
		survey.Required)

	// addr2
	var addr2 string
	survey.AskOne(
		&survey.Input{
			Message: "Address 2:",
		},
		&addr2,
		nil)
	if addr2 != "" {
		req.Addr2 = &addr2
	}

	// city
	survey.AskOne(
		&survey.Input{
			Message: "City:",
		},
		&req.City,
		nil)

	// county
	var county string
	survey.AskOne(
		&survey.Input{
			Message: "County:",
		},
		&county,
		nil)
	if county != "" {
		req.County = &county
	}

	// postcode
	survey.AskOne(
		&survey.Input{
			Message: "Postcode:",
		}, &req.Postcode, survey.Required)

	// country_code
	var countryCode string
	c := &survey.Select{
		Message: "Country Code:",
		Options: service.CountryCodes(),
	}
	survey.AskOne(c, &countryCode, survey.Required)
	req.CountryCode = countryCode[0:2]

	return &req, nil
}
