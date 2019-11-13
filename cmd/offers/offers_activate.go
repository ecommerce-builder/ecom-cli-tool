package offers

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/ecommerce-builder/ecom-cli-tool/service"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// NewCmdOffersActivate returns new initialized instance of activate sub command
func NewCmdOffersActivate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "activate",
		Short: "Active an offer",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			// get the request params
			ctx := context.Background()
			req, err := promptCreateOffer(ctx, client)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			offer, err := client.CreateOffer(ctx, req)
			if err == eclient.ErrOfferExists {
				fmt.Fprintf(os.Stderr,
					"offer with promo rule id %s already active\n",
					req.PromoRuleID)
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			showOffer(offer)
		},
	}
	return cmd
}

func promptCreateOffer(ctx context.Context, client *eclient.EcomClient) (*eclient.CreateOfferRequest, error) {
	var req eclient.CreateOfferRequest

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

	return &req, nil
}

func showOffer(v *eclient.Offer) {
	format := "%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Offer ID:", v.ID)
	fmt.Fprintf(tw, format, "Promo Rule ID:", v.PromoRuleID)
	fmt.Fprintf(tw, format, "Promo Rule Code:", "todo")
	fmt.Fprintf(tw, format, "Created:", v.Created.In(service.Location).Format(service.TimeDisplayFormat))
	fmt.Fprintf(tw, format, "Modified:", v.Modified.In(service.Location).Format(service.TimeDisplayFormat))
	tw.Flush()
}
