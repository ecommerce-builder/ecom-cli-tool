package coupons

import (
	"context"
	"fmt"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// NewCmdCouponsCreate returns new initialized instance of create sub command
func NewCmdCouponsCreate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Mints a new coupon",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			// get the request params
			ctx := context.Background()
			req, err := promptCreateCoupon(ctx, client)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			coupon, err := client.CreateCoupon(ctx, req)
			if err != nil {
				fmt.Printf("%+v\n", err)
				os.Exit(1)
			}
			showCoupon(coupon)
		},
	}
	return cmd

}

func promptCreateCoupon(ctx context.Context, client *eclient.EcomClient) (*eclient.CreateCouponRequest, error) {
	var req eclient.CreateCouponRequest

	// build a list of promo rules
	promoRules, err := client.GetPromoRules(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "get promo rules")
	}

	promoRulesMap := make(map[string]string, 0)
	var promoRuleOpts []string
	for _, p := range promoRules {
		promoRuleOpts = append(promoRuleOpts, p.PromoRuleCode)
		promoRulesMap[p.PromoRuleCode] = p.ID
	}

	// promo_rule_id
	var promoRule string
	prompt := &survey.Select{
		Message: "Promo Rule Code:",
		Options: promoRuleOpts,
	}
	survey.AskOne(prompt, &promoRule, survey.Required)
	req.PromoRuleID = promoRulesMap[promoRule]

	// coupon_code
	c := &survey.Input{
		Message: "Coupon Code:",
	}
	survey.AskOne(c, &req.CouponCode, survey.Required)

	// reusable
	d := &survey.Confirm{
		Message: "Make coupon reusable?",
	}
	survey.AskOne(d, &req.Resuable, survey.Required)

	return &req, nil
}
