package promorules

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

// NewCmdPromoRulesList returns new initialized instance of the list sub command
func NewCmdPromoRulesList() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "list",
		Short: "list price lists",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			ctx := context.Background()
			promoRules, err := client.GetPromoRules(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			format := "%s\t%s\t%s\t%s\t%s\t%v\t%s\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "Promo Rule code", "Name", "Start At", "End At", "Type", "Amount", "Target")
			fmt.Fprintf(tw, format, "---------------", "----", "--------", "------", "----", "------", "------")
			for _, p := range promoRules {
				var startAt, endAt string
				if p.StartAt == nil {
					startAt = "-"
				} else {
					startAt = p.StartAt.In(location).Format(timeDisplayFormat)
				}
				if p.EndAt == nil {
					endAt = "-"
				} else {
					endAt = p.EndAt.In(location).Format(timeDisplayFormat)
				}

				var amount string
				if p.Type == "percentage" {
					amount = fmt.Sprintf("%.2f%%", float64(p.Amount)/100.0)
				} else {
					amount = fmt.Sprintf("£%.4f", float64(p.Amount)/10000.0)
				}
				fmt.Fprintf(tw, format, p.PromoRuleCode, p.Name, startAt,
					endAt, p.Type, amount, p.Target)
			}
			tw.Flush()
		},
	}
	return cmd
}
