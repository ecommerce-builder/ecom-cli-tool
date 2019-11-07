package inventory

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
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

// NewCmdInventoryGet returns new initialized instance of the get sub command
func NewCmdInventoryGet() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "get <inventory_id>",
		Short: "Get an inventory by id",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			invID := args[0]
			if !cmdvalidate.IsValidUUID(invID) {
				fmt.Fprintf(os.Stderr, "inventory_id %q is not a valid v4 uuid\n", invID)
				os.Exit(1)
			}

			ctx := context.Background()
			inventory, err := client.GetInventory(ctx, invID)
			if err == eclient.ErrInventoryNotFound {
				fmt.Fprintf(os.Stderr, "inventory %q not found\n", invID)
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			showInventory(inventory)
		},
	}
	return cmd
}

func showInventory(v *eclient.Inventory) {
	format := "%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Inventory ID:", v.ID)
	fmt.Fprintf(tw, format, "Product ID:", v.ProductID)
	fmt.Fprintf(tw, format, "Product Path:", v.ProductPath)
	fmt.Fprintf(tw, format, "Product SKU:", v.ProductSKU)
	fmt.Fprintf(tw, format, "Onhand:", v.Onhand)
	fmt.Fprintf(tw, format, "Overselling:", v.Overselling)
	fmt.Fprintf(tw, format, "Created:",
		v.Created.In(location).Format(timeDisplayFormat))
	fmt.Fprintf(tw, format, "Modified:",
		v.Modified.In(location).Format(timeDisplayFormat))
	tw.Flush()
}
