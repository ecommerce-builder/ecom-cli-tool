package cmd

import (
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/carts"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/categoriestree"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/devkeys"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/pcrelations"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/ppagroups"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/ppassocs"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/pricelists"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/products"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/profiles"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/promorules"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/tariffs"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/token"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/users"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/webhooks"

	"github.com/spf13/cobra"
)

// var rc *configmgr.EcomConfigurations
// var currentConfigName string

// NewEcomCmd creates the `ecom` command.
func NewEcomCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "ecom",
		Short: "ecom is a CLI tool for administering ecommerce systems",
		Long:  `See the user guide for more details.`,
	}
	cmd.AddCommand(carts.NewCmdCarts())
	cmd.AddCommand(pcrelations.NewCmdPCRelations())
	cmd.AddCommand(categoriestree.NewCmdCategoriesTree())
	cmd.AddCommand(users.NewCmdUsers())
	cmd.AddCommand(devkeys.NewCmdDevKeys())
	cmd.AddCommand(products.NewCmdProducts())
	cmd.AddCommand(ppassocs.NewCmdPPAssocs())
	cmd.AddCommand(ppagroups.NewCmdPPAGroups())
	cmd.AddCommand(pricelists.NewCmdPriceLists())
	cmd.AddCommand(profiles.NewCmdProfiles())
	cmd.AddCommand(promorules.NewCmdPromoRules())
	cmd.AddCommand(tariffs.NewCmdShippingTariffsRules())
	cmd.AddCommand(webhooks.NewCmdWebhooks())
	cmd.AddCommand(NewCmdCompletion())
	cmd.AddCommand(NewCmdSysInfo())
	cmd.AddCommand(token.NewCmdToken())
	cmd.AddCommand(NewCmdVersion())
	return cmd
}
