package cmd

import (
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/admins"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/assocs"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/catalog"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/customers"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/products"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/profiles"
	"github.com/ecommerce-builder/ecom-cli-tool/cmd/token"

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
	cmd.AddCommand(admins.NewCmdAdmins())
	cmd.AddCommand(assocs.NewCmdAssocs())
	cmd.AddCommand(catalog.NewCmdCatalog())
	cmd.AddCommand(customers.NewCmdCustomers())
	cmd.AddCommand(products.NewCmdProducts())
	cmd.AddCommand(profiles.NewCmdProfiles())
	cmd.AddCommand(NewCmdCompletion())
	cmd.AddCommand(NewCmdSysInfo())
	cmd.AddCommand(token.NewCmdToken())
	cmd.AddCommand(NewCmdVersion())
	return cmd
}
