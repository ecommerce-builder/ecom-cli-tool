package cmd

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v3"

	"bitbucket.org/andyfusniakteam/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
	"gopkg.in/go-playground/validator.v9"
)

var productsGetCmd = &cobra.Command{
	Use:   "get <sku>|<file.yaml>",
	Short: "Get product",
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

		product, err := ecomClient.GetProduct(sku)
		if err != nil {
			log.Fatal(err)
		}

		b, err := yaml.Marshal(product)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", string(b))
	},
}

func init() {
	productsCmd.AddCommand(productsGetCmd)
}
