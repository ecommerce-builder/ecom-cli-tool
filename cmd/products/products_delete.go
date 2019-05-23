package products

import (
	"fmt"
	"log"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
	"gopkg.in/go-playground/validator.v9"
)

// NewCmdProductsDelete returns new initialized instance of the delete sub command
func NewCmdProductsDelete() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "delete <sku>",
		Short: "Delete product",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}

			sku := args[0]
			validate := validator.New()
			err = validate.Var(sku, "gt=1,lt=10")
			if err != nil {
				log.Fatal(err)
			}

			err = client.DeleteProduct(sku)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	return cmd
}
