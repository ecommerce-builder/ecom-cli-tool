package ppassocs

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

// NewCmdPPAssocsList returns new initialized instance of the list sub command
func NewCmdPPAssocsList() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "list <ppa_group_code>",
		Short: "list product to product associations for a given group",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			ppaGroupCode := args[0]
			ctx := context.Background()

			groups, err := client.GetPPAGroups(ctx)
			if err != nil {
				fmt.Printf("%+v\n", err)
				os.Exit(1)
			}
			var ppaGroupID string
			for _, g := range groups {
				if g.Code == ppaGroupCode {
					ppaGroupID = g.ID
					break
				}
			}
			if ppaGroupID == "" {
				fmt.Fprintf(os.Stderr, "ppa group code %q not found\n",
					ppaGroupCode)
				os.Exit(1)
			}

			assocs, err := client.GetPPAssocs(ctx, ppaGroupID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			// create a lookup of product id to product skus
			products, err := client.GetProducts(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "")
			}

			productMap := make(map[string]string)
			for _, v := range products {
				productMap[v.ID] = v.SKU
			}

			format := "%s\t%s\t%s\t%s\t%v\t%v\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "PP Assocs ID", "Group Code", "Product From",
				"Product To", "Created", "Modified")
			fmt.Fprintf(tw, format, "------------", "----------", "------------",
				"----------", "-------", "--------")
			for _, v := range assocs {
				fmt.Fprintf(tw, format,
					v.ID,
					ppaGroupCode,
					productMap[v.ProductFromID],
					productMap[v.ProductToID],
					v.Created.In(location).Format(timeDisplayFormat),
					v.Modified.In(location).Format(timeDisplayFormat))
			}
			tw.Flush()
		},
	}
	return cmd
}
