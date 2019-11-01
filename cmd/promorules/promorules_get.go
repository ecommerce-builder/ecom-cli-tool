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
			if err == eclient.ErrPromoRuleNotFound {
				fmt.Fprintf(os.Stderr, "promo rule %q (%q) not found\n", promoRuleCode, promoRuleID)
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			showPromoRule(promoRule)
		},
	}
	return cmd
}

func showPromoRule(promoRule *eclient.PromoRule) {
	format := "%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Product Rule ID:", promoRule.ID)
	fmt.Fprintf(tw, format, "Promo Rule Code:", promoRule.PromoRuleCode)
	fmt.Fprintf(tw, format, "Name:", promoRule.Name)
	fmt.Fprintf(tw, format, "Type:", promoRule.Type)
	if promoRule.Type == "percentage" {
		fmt.Fprintf(tw, "%v\t%.2f%%\t\n", "Amount:", float64(promoRule.Amount)/100.0)
	} else {
		fmt.Fprintf(tw, "%v\t%.4f\t\n", "Amount", float64(promoRule.Amount)/100.0)
	}

	var startAt, endAt string
	if promoRule.StartAt == nil {
		startAt = "-"
	} else {
		startAt = promoRule.StartAt.In(location).Format(timeDisplayFormat)
	}
	if promoRule.EndAt == nil {
		endAt = "-"
	} else {
		endAt = promoRule.EndAt.In(location).Format(timeDisplayFormat)
	}
	fmt.Fprintf(tw, format, "Start At:", startAt)
	fmt.Fprintf(tw, format, "End At:", endAt)
	fmt.Fprintf(tw, format, "Target:", promoRule.Target)
	switch promoRule.Target {
	case "product":
		fmt.Fprintf(tw, format, "Product ID:", *promoRule.ProductID)
		fmt.Fprintf(tw, format, "Product Path:", *promoRule.ProductPath)
		fmt.Fprintf(tw, format, "Product SKU:", *promoRule.ProductSKU)
	case "productset":
		fmt.Fprintf(tw, format, "Product Set ID", *promoRule.ProductSetID)
	case "category":
		fmt.Fprintf(tw, format, "Category ID:", *promoRule.CategoryID)
		fmt.Fprintf(tw, format, "Category Path:", *promoRule.CategoryPath)
	case "shipping_tariff":
		fmt.Fprintf(tw, format, "Shipping Tariff ID:", *promoRule.ShippingTariffID)
		fmt.Fprintf(tw, format, "Shipping Tariff Code:", *promoRule.ShippingTariffCode)
	case "total":
		fmt.Fprintf(tw, "%v\t£%.2f\t\n", "Total Threshold",
			float64(*promoRule.TotalThreshold)/10000.0)
	default:
		fmt.Fprintf(os.Stderr, "unknown promo rule target %q\n", promoRule.Target)
		os.Exit(1)
	}
	fmt.Fprintf(tw, format, "Created:", promoRule.Created.In(location).Format(timeDisplayFormat))
	fmt.Fprintf(tw, format, "Modified:", promoRule.Modified.In(location).Format(timeDisplayFormat))
	tw.Flush()
}
