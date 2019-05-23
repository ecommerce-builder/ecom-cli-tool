package products

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdProductsApply returns new initialized instance of the apply sub command
func NewCmdProductsApply() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "apply <product.yaml>",
		Short: "Create or update an exising product",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}

			data, err := ioutil.ReadFile(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			var p struct {
				Product eclient.Product `yaml:"product"`
			}
			err = yaml.Unmarshal([]byte(data), &p)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			exists, err := client.ProductExists(p.Product.SKU)
			if err != nil {
				log.Fatalf("error: %+v", err)
			}
			pc := eclient.ProductCreate{
				SKU:  p.Product.SKU,
				EAN:  p.Product.EAN,
				URL:  p.Product.URL,
				Name: p.Product.Name,
				Data: p.Product.Data,
			}
			if !exists {
				_, err = client.CreateProduct(&pc)
				if err != nil {
					log.Fatalf("create product failed sku=%q: %+v", p.Product.SKU, err)
				}
				os.Exit(0)
			}
			pu := eclient.ProductUpdate{
				EAN:  p.Product.EAN,
				URL:  p.Product.URL,
				Name: p.Product.Name,
				Data: p.Product.Data,
			}
			_, err = client.UpdateProduct(p.Product.SKU, &pu)
			if err != nil {
				log.Fatalf("update product failed sku=%q: %+v", p.Product.SKU, err)
			}
			os.Exit(0)
		},
	}
	return cmd
}
