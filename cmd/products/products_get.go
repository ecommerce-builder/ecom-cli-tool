package products

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

type productResponseContainer struct {
	Product *eclient.ProductResponse `yaml:"product"`
}

// NewCmdProductsGet returns new initialized instance of the get sub command
func NewCmdProductsGet() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "get <sku>",
		Short: "Get product",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}
			sku := args[0]
			product, err := client.GetProduct(sku)
			if err != nil {
				if err == eclient.ErrProductNotFound {
					fmt.Printf("Product %s not found.\n", sku)
					os.Exit(0)
				}
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			container := productResponseContainer{
				Product: product,
			}
			enc := yaml.NewEncoder(os.Stdout)
			//enc.SetIndent(2)
			if err := enc.Encode(container); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			if err := enc.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			//b, err := yaml.Marshal(product)
			//if err != nil {
			//	log.Fatal(err)
			//}
			//fmt.Printf("%s\n", string(b))
		},
	}
	return cmd
}
