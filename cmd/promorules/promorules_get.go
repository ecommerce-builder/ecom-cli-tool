package promorules

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

type productResponseContainer struct {
	Product *eclient.ProductResponse `yaml:"product"`
}

// NewCmdPromoRulesGet returns new initialized instance of the get sub command
func NewCmdPromoRulesGet() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "get <promo_rule_code>",
		Short: "Get promo rule",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}

			// promo_rule_code to id
			promoRuleCode := args[0]
			ctx := context.Background()
			promoRules, err := client.GetPromoRules(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			var promoRuleID string
			for _, pr := range promoRules {
				if pr.PromoRuleCode == promoRuleCode {
					promoRuleID = pr.ID
					break
				}
			}
			if promoRuleID == "" {
				fmt.Fprintf(os.Stderr, "promo_rule_code %q not found\n", promoRuleCode)
				os.Exit(1)
			}

			promoRule, err := client.GetPromoRule(ctx, promoRuleID)
			if err != nil {
				if err == eclient.ErrProductNotFound {
					fmt.Fprintf(os.Stderr, "Promo rule %q not found.\n", promoRuleID)
					os.Exit(0)
				}
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			// fmt.Printf("%+v\n", promoRule)

			format := "%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)

			fmt.Fprintf(tw, format, "Promo rule code", promoRule.PromoRuleCode)
			fmt.Fprintf(tw, format, "Name", promoRule.Name)
			fmt.Fprintf(tw, format, "Type", promoRule.Type)
			fmt.Fprintf(tw, format, "Target", promoRule.Target)

			switch promoRule.Target {
			case "product":
				fmt.Fprintf(tw, format, "Product ID", promoRule.ProductID)
			default:
				fmt.Fprintf(os.Stderr, "unrecognised promo rule target %q\n", promoRule.Target)
				os.Exit(1)
			}
			tw.Flush()
		},
	}
	return cmd
}
