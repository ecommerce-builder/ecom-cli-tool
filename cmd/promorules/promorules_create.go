package promorules

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

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

			ctx := context.Background()
			promoRule, err := client.CreatePromoRule(ctx, req)
			if err != nil {
				fmt.Printf("%+v\n", err)
				fmt.Fprintf(os.Stderr, "error creating promo rule: %+v\n", errors.Unwrap(err))
				os.Exit(1)
			}
			showPromoRule(promoRule)
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
	survey.AskOne(p, &req.PromoRuleCode, survey.Required)

	// name
	n := &survey.Input{
		Message: "Name:",
	}
	survey.AskOne(n, &req.Name, survey.Required)

	// use date a date range?
	useRange := false
	d := &survey.Confirm{
		Message: "Use a date range?",
	}
	survey.AskOne(d, &useRange, survey.Required)

	if useRange {
		st := &survey.Input{
			Message: "Start At (DD/MM/YYYY):",
		}
		var startAt string
		var startAtT time.Time
		survey.AskOne(st, &startAt, survey.ComposeValidators(
			survey.Required,
			func(val interface{}) error {
				ts, ok := val.(string)
				if !ok {
					return errors.New("invalid response")
				}
				ts = ts + " 23:59:59"

				var err error
				startAtT, err = time.Parse("02/01/2006 15:04:05", ts)
				if err != nil {
					return err
				}
				return nil
			},
		))
		req.StartAt = &startAtT

		en := &survey.Input{
			Message: "End At (DD/MM/YYYY):",
		}
		var endAt string
		var endAtT time.Time
		survey.AskOne(en, &endAt, survey.ComposeValidators(
			survey.Required,
			func(val interface{}) error {
				ts, ok := val.(string)
				if !ok {
					return errors.New("invalid response")
				}
				ts = ts + " 00:00:00"

				var err error
				endAtT, err = time.Parse("02/01/2006 15:04:05", ts)
				if err != nil {
					return err
				}
				return nil
			},
		))
		req.EndAt = &endAtT
	}

	// type
	t := &survey.Select{
		Message: "Type:",
		Options: []string{"percentage", "fixed"},
	}
	survey.AskOne(t, &req.Type, survey.Required)

	// amount
	a := &survey.Input{
		Message: "Amount:",
	}
	survey.AskOne(a, &req.Amount, survey.ComposeValidators(
		survey.Required,
		func(val interface{}) error {
			str, ok := val.(string)
			if !ok {
				return errors.New("invalid response")
			}

			v, err := strconv.Atoi(str)
			if err != nil {
				return err
			}

			if req.Type == "percentage" {
				if v < 0 || v > 10000 {
					return errors.New("amount must be between 0 and 10000 (0.00% to 100.00%)")
				}
			}
			if req.Type == "fixed" {
				if v < 0 {
					return errors.New("amount must be a positive integer")
				}
			}
			return nil
		},
	))

	// target
	g := &survey.Select{
		Message: "Target:",
		Options: []string{"product", "productset", "category", "shipping_tariff", "total"},
	}
	survey.AskOne(g, &req.Target, survey.Required)

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
			Message: "Product ID:",
			Options: productOpts,
		}
		survey.AskOne(prompt, &sku, survey.Required)
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
		survey.AskOne(prompt, &path, survey.Required)
		req.CategoryID = categoryMap[path]
	case "shipping_tariff":
		// build a list and map of shipping tariffs
		ctx := context.Background()
		tariffs, err := client.GetShippingTariffs(ctx)
		if err != nil {
			return nil, fmt.Errorf("%w: client.GetShippingTariffs(ctx) failed", err)
		}
		if len(tariffs) < 1 {
			fmt.Fprintf(os.Stderr, "No shipping tariffs have been created.\n")
			os.Exit(1)
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
		survey.AskOne(prompt, &shippingCode, survey.Required)
		req.ShippingTariffID = tariffMap[shippingCode]
	case "total":
		// total_threshold
		t := &survey.Input{
			Message: "Order total theshold:",
		}
		survey.AskOne(t, &req.TotalThreshold, survey.ComposeValidators(
			survey.Required,
			func(val interface{}) error {
				str, ok := val.(string)
				if !ok {
					return errors.New("invalid response")
				}

				v, err := strconv.Atoi(str)
				if err != nil {
					return err
				}

				if v <= 0 {
					return errors.New("amount must be a positive integer")
				}
				return nil
			},
		))
	default:
		return nil, fmt.Errorf("target %q is current not support by this command line tool", req.Target)
	}

	return &req, nil
}
