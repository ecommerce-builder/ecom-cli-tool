package cmd

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v3"

	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

var productsGetCmd = &cobra.Command{
	Use:   "get <sku>|<file.yaml>",
	Short: "Get product",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		current := rc.Configurations[currentConfigName]
		ecomClient := eclient.New(current.Endpoint, timeout)
		if err := ecomClient.SetToken(&current); err != nil {
			log.Fatal(err)
		}
		sku := args[0]
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
