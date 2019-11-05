package address

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

// NewCmdAddressList returns new initialized instance of the list sub command
func NewCmdAddressList() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "list <email>",
		Short: "list addressess for a user",
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
				fmt.Fprintf(os.Stderr, "email %s did not match any users\n", email)
				os.Exit(1)
			}

			addresses, err := client.GetAddressesByUser(ctx, userMap[email])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			format := "%s\t%s\t%s\t%v\t%s\t%v\t%s\t%s\t%v\t%v\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "Address ID", "Contact Name",
				"Address 1", "Address 2", "City", "County",
				"Postcode", "Country Code", "Created", "Modified")
			fmt.Fprintf(tw, format, "----------", "------------",
				"---------", "---------", "----", "------",
				"--------", "------------", "-------", "--------")
			for _, a := range addresses {
				addr2 := "-"
				if a.Addr2 != nil {
					addr2 = *a.Addr2
				}
				county := "-"
				if a.County != nil {
					county = *a.County
				}
				fmt.Fprintf(tw, format, a.ID,
					a.ContactName, a.Addr1, addr2,
					a.City, county, a.Postcode,
					a.CountryCode,
					a.Created.In(location).Format(timeDisplayFormat),
					a.Modified.In(location).Format(timeDisplayFormat))
			}
			fmt.Fprintf(tw, "(%d addresses found)\n", len(addresses))
			tw.Flush()
		},
	}
	return cmd
}
