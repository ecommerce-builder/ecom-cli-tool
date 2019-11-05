package coupons

import (
	"github.com/spf13/cobra"
)

// NewCmdCoupons returns new initialized instance of coupons sub command
func NewCmdCoupons() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "coupons",
		Short: "Coupons management",
	}
	cmd.AddCommand(NewCmdCouponsCreate())
	cmd.AddCommand(NewCmdCouponsGet())
	cmd.AddCommand(NewCmdCouponsList())
	cmd.AddCommand(NewCmdCouponsDelete())
	cmd.AddCommand(NewCmdCouponsVoid())
	return cmd
}
