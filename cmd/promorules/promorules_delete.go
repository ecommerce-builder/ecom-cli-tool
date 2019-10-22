package promorules

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdPromoRulesDelete returns new initialized instance of the delete sub command
func NewCmdPromoRulesDelete() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "delete <promo_rule_code>",
		Short: "Delete a promo rule",
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

			err = client.DeletePromoRule(ctx, promoRuleID)
			if err == eclient.ErrBadRequest {
				fmt.Fprintf(os.Stderr, "Bad request - this is likely an error with the command line tool - please report this\n")
				os.Exit(1)
			}
			if err == eclient.ErrPromoRuleNotFound {
				fmt.Fprintf(os.Stderr, "Promo rule not found. Use ecom promorules list to check.\n")
				os.Exit(1)
			}
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	return cmd
}
