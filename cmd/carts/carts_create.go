package carts

import (
	"context"
	"fmt"
	"log"
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

// NewCmdCartsCreate returns new initialized instance of the create sub command
func NewCmdCartsCreate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new shopping cart",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}

			// load all price lists
			ctx := context.Background()
			cart, err := client.CreateCart(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			format := "%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "Cart ID:", cart.ID)
			fmt.Fprintf(tw, format, "Locked:", cart.Locked)
			fmt.Fprintf(tw, format, "Created:",
				cart.Created.In(location).Format(timeDisplayFormat))
			fmt.Fprintf(tw, format, "Modified:",
				cart.Modified.In(location).Format(timeDisplayFormat))
			tw.Flush()

			fmt.Fprintf(os.Stdout, "export ECOM_CLI_CART_ID=%s\n", cart.ID)
		},
	}
	return cmd
}
