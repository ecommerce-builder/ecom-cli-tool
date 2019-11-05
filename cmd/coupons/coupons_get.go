package coupons

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdCouponsGet returns new initialized instance of the get sub command
func NewCmdCouponsGet() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "get <coupon_code>",
		Short: "Get coupon code",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			// coupon_code to id
			couponCode := args[0]
			ctx := context.Background()
			coupons, err := client.GetCoupons(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			var coupon *eclient.Coupon
			for _, c := range coupons {
				if c.CouponCode == couponCode {
					coupon = c
					break
				}
			}
			if coupon == nil {
				fmt.Fprintf(os.Stderr, "coupon_code %q not found\n", couponCode)
				os.Exit(1)
			}

			showCoupon(coupon)
		},
	}
	return cmd
}

func showCoupon(v *eclient.Coupon) {
	format := "%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Coupon ID:", v.ID)
	fmt.Fprintf(tw, format, "Promo Rule ID:", v.PromoRuleID)
	fmt.Fprintf(tw, format, "Promo Rule Code:", v.PromoRuleCode)
	fmt.Fprintf(tw, format, "Resuable:", v.Resuable)
	fmt.Fprintf(tw, format, "Void:", v.Void)
	fmt.Fprintf(tw, format, "Spent Count:", v.SpendCount)
	fmt.Fprintf(tw, format, "Created:",
		v.Created.In(location).Format(timeDisplayFormat))
	fmt.Fprintf(tw, format, "Modified:",
		v.Modified.In(location).Format(timeDisplayFormat))
	tw.Flush()
}
