package cmd

import (
	"log"

	"bitbucket.org/andyfusniakteam/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
	"gopkg.in/go-playground/validator.v9"
)

var productsDeleteCmd = &cobra.Command{
	Use:   "delete <sku>",
	Short: "Delete product",
	Long:  ``,

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		current := rc.Configurations[currentConfigName]
		ecomClient := eclient.New(current.Endpoint, timeout)
		err := ecomClient.SetToken(&current)
		if err != nil {
			log.Fatal(err)
		}

		sku := args[0]
		validate := validator.New()
		err = validate.Var(sku, "gt=1,lt=10")
		if err != nil {
			log.Fatal(err)
		}

		err = ecomClient.DeleteProduct(sku)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	productsCmd.AddCommand(productsDeleteCmd)
}
