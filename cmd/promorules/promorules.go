package promorules

import (
	"github.com/spf13/cobra"
)

// NewCmdPromoRules returns new initialized instance of promorules sub command
func NewCmdPromoRules() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "promorules",
		Short: "Promotion Rules Management",
	}
	cmd.AddCommand(NewCmdPromoRulesCreate())
	cmd.AddCommand(NewCmdPromoRulesGet())
	cmd.AddCommand(NewCmdPromoRulesList())
	cmd.AddCommand(NewCmdPromoRulesDelete())
	return cmd
}
