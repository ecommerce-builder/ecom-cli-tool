package offers

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/ecommerce-builder/ecom-cli-tool/service"
	"github.com/spf13/cobra"
)

// NewCmdOffersList returns new initialized instance of list sub command
func NewCmdOffersList() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List offers",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			ctx := context.Background()
			offers, err := client.GetOffers(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			format := "%v\t%v\t%v\t%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t\n",
				"Offer ID",
				"Promo Rule Code",
				"Promo Rule ID",
				"Created",
				"Modified")
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t\n",
				"--------",
				"---------------",
				"-------------",
				"-------",
				"--------")
			for _, v := range offers {
				fmt.Fprintf(tw, format,
					v.ID,
					v.PromoRuleCode,
					v.PromoRuleID,
					v.Created.In(service.Location).Format(service.TimeDisplayFormat),
					v.Modified.In(service.Location).Format(service.TimeDisplayFormat))
			}
			tw.Flush()
		},
	}
	return cmd
}
