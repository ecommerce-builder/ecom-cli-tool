package coupons

import (
	"context"
	"fmt"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdCouponsDelete returns new initialized instance of the delete sub command
func NewCmdCouponsDelete() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "delete <coupon_code>",
		Short: "Delete a coupon",
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

			err = client.DeleteCoupon(ctx, coupon.ID)
			if err == eclient.ErrCouponNotFound {
				fmt.Fprintf(os.Stderr, "coupon not found. Use ecom coupons list to check.\n")
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
		},
	}
	return cmd
}
