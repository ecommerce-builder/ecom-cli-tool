package address

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdAddressGet returns new initialized instance of the get sub command
func NewCmdAddressGet() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "get <address_id>",
		Short: "Get address by id",
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

			showAddress(addr)
		},
	}
	return cmd
}

func showAddress(v *eclient.Address) {
	addr2 := "-"
	if v.Addr2 != nil {
		addr2 = *v.Addr2
	}
	county := "-"
	if v.County != nil {
		county = *v.County
	}
	format := "%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Address ID:", v.ID)
	fmt.Fprintf(tw, format, "User ID:", v.UserID)
	fmt.Fprintf(tw, format, "Type:", v.Typ)
	fmt.Fprintf(tw, format, "Contact Name:", v.ContactName)
	fmt.Fprintf(tw, format, "Address 1:", v.Addr1)
	fmt.Fprintf(tw, format, "Address 2:", addr2)
	fmt.Fprintf(tw, format, "City:", v.City)
	fmt.Fprintf(tw, format, "County:", county)
	fmt.Fprintf(tw, format, "Postcode:", v.Postcode)
	fmt.Fprintf(tw, format, "Country Code:", v.CountryCode)
	fmt.Fprintf(tw, format, "Created",
		v.Created.In(location).Format(timeDisplayFormat))
	fmt.Fprintf(tw, format, "Modified",
		v.Modified.In(location).Format(timeDisplayFormat))
	tw.Flush()
}
