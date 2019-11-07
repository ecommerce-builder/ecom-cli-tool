package inventory

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdInventoryList returns new initialized instance of list sub command
func NewCmdInventoryList() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List inventory",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			ctx := context.Background()
			inv, err := client.GetAllInventory(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			format := "%s\t%s\t%v\t%v\t%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t\n",
				"Inventory ID",
				"Product SKU",
				"Onhand",
				"Overselling",
				"Created",
				"Modified")
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t\n",
				"------------",
				"-----------",
				"------",
				"-----------",
				"-------",
				"--------")
			for _, v := range inv {
				fmt.Fprintf(tw, format,
					v.ID,
					v.ProductSKU,
					v.Onhand,
					v.Overselling,
					v.Created.In(location).Format(timeDisplayFormat),
					v.Modified.In(location).Format(timeDisplayFormat))
			}
			tw.Flush()
		},
	}
	return cmd
}
