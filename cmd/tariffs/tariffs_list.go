package tariffs

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

// NewCmdShippingTariffsList returns new initialized instance of the list sub command
func NewCmdShippingTariffsList() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "list",
		Short: "list shipping tariffs",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			tariffs, err := client.GetShippingTariffs(context.TODO())
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			format := "%s\t%s\t%s\t%s\t%v\t%s\t%v\t%v\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "Shipping Tariff ID", "Shipping Code", "Country Code",
				"Name", "Price", "Tax code", "Created", "Modified")
			fmt.Fprintf(tw, format, "------------------", "-------------", "------------",
				"----", "-----", "--------", "-------", "--------")

			for _, t := range tariffs {
				fmt.Fprintf(tw, format,
					t.ID, t.ShippingCode, t.CountryCode,
					t.Name, t.Price, t.TaxCode,
					t.Created.In(location).Format(timeDisplayFormat),
					t.Modified.In(location).Format(timeDisplayFormat))
			}
			tw.Flush()
		},
	}
	return cmd
}
