package promorules

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// type promoRuleRequest struct {
// 	promoRuleCode    string
// 	name             string
// 	typ              string
// 	target           string
// 	productID        string
// 	categoryID       string
// 	shippingTariffID string
// 	totalTheshold    int
// }

// NewCmdPromoRulesCreate returns new initialized instance of create sub command
func NewCmdPromoRulesCreate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new promo rule",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}

			req, err := promptCreatePromoRule(client)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			fmt.Println(req.PromoRuleCode)
			fmt.Println(req.Name)
			fmt.Println(req.Type)
			fmt.Println(req.Target)
			fmt.Println(req.ProductID)
			fmt.Println(req.CategoryID)
			fmt.Println(req.ShippingTariffID)
			fmt.Println(req.TotalThreshold)

			ctx := context.Background()
			tariff, err := client.CreatePromoRule(ctx, req)
			if err != nil {
				fmt.Printf("%+v\n", err)
				fmt.Fprintf(os.Stderr, "error creating promo rule: %+v", errors.Unwrap(err))
			}

			fmt.Println(tariff)
		},
	}
	return cmd
}

func promptCreatePromoRule(client *eclient.EcomClient) (*eclient.PromoRuleRequest, error) {
	var req eclient.PromoRuleRequest

	// promo_rule_code
	p := &survey.Input{
		Message: "Promo rule code:",
	}
	survey.AskOne(p, &req.PromoRuleCode, nil)

	// name
	n := &survey.Input{
		Message: "Name:",
	}
	survey.AskOne(n, &req.Name, nil)

	// type
	t := &survey.Select{
		Message: "Type:",
		Options: []string{"percentage", "fixed"},
	}
	survey.AskOne(t, &req.Type, nil)

	// target
	g := &survey.Select{
		Message: "Target:",
		Options: []string{"product", "productset", "category", "shipping_tariff", "total"},
	}
	survey.AskOne(g, &req.Target, nil)

	switch req.Target {
	case "product":
		// build a list of products
		products, err := client.GetProducts(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("%w: client.GetProducts(ctx) failed", err)
		}

		productMap := make(map[string]string, 0)
		var productOpts []string
		for _, product := range products {
			productOpts = append(productOpts, product.SKU)
			productMap[product.SKU] = product.ID
		}

		// product_id
		var sku string
		prompt := &survey.Select{
			Message: "Category ID:",
			Options: productOpts,
		}
		survey.AskOne(prompt, &sku, nil)
		req.ProductID = productMap[sku]
	case "category":
		// build a list of categories
		categories, err := client.GetCategories()
		if err != nil {
			return nil, fmt.Errorf("%w: eclient.GetCategories() failed", err)
		}

		categoryMap := make(map[string]string, 0)
		var categoryOpts []string
		for _, category := range categories {
			categoryOpts = append(categoryOpts, category.Path)
			categoryMap[category.Path] = category.ID
		}

		// category_id
		var path string
		prompt := &survey.Select{
			Message: "Category ID:",
			Options: categoryOpts,
		}
		survey.AskOne(prompt, &path, nil)
		req.CategoryID = categoryMap[path]
	case "shipping_tariff":
		// build a list and map of shipping tariffs
		tariffs, err := client.GetShippingTariffs(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("%w: client.GetShippingTariffs(ctx) failed", err)
		}

		tariffMap := make(map[string]string, 0)
		var tariffOpts []string
		for _, tariff := range tariffs {
			tariffOpts = append(tariffOpts, tariff.ShippingCode)
			tariffMap[tariff.ShippingCode] = tariff.ID
		}

		// shipping_tariff_id
		var shippingCode string
		prompt := &survey.Select{
			Message: "Shipping Tariff ID:",
			Options: tariffOpts,
		}
		survey.AskOne(prompt, &shippingCode, nil)
		req.ShippingTariffID = tariffMap[shippingCode]
	case "total":
		// total_threshold
		t := &survey.Input{
			Message: "Order total theshold:",
		}
		survey.AskOne(t, &req.TotalThreshold, nil)
	default:
		return nil, fmt.Errorf("target %q is current not support by this command line tool", req.Target)
	}

	return &req, nil
}
