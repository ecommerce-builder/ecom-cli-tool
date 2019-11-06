package pcrelations

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdPCRelationsApply returns new initialized instance of apply sub command
func NewCmdPCRelationsApply() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "apply <product-category-relations.yaml>",
		Short: "Replace all product to category relations",
		Long:  ``,

		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			err := client.SetToken(&current)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			data, err := ioutil.ReadFile(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			var relationships eclient.ProductCategoryRelationsYAML
			err = yaml.Unmarshal([]byte(data), &relationships)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			// retrieve a list of all products and build a map
			// of sku -> product ids
			products, err := client.GetProducts(context.TODO())
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get products: %+v", err)
				os.Exit(1)
			}
			productSKUToID := make(map[string]string)
			for _, p := range products {
				productSKUToID[p.SKU] = p.ID
			}

			// fmt.Printf("%#v\n", productSKUToID)

			// retrieve a list of all categories and build a map
			// of path -> category ids
			categories, err := client.GetCategories()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get categories: %+v", err)
				os.Exit(1)
			}
			categoryPathToID := make(map[string]string)
			for _, c := range categories {
				categoryPathToID[c.Path] = c.ID
			}

			for path, productset := range relationships.Rels {
				if _, ok := categoryPathToID[path]; !ok {
					fmt.Fprintf(os.Stderr, "Category path %s not found.\n", path)
				}

				for _, sku := range productset.Products {
					if _, ok := productSKUToID[sku]; !ok {
						fmt.Fprintf(os.Stderr, "Product SKU=%s in Category path=%s section not found.\n", sku, path)
					}
				}
			}

			rels := make([]*eclient.CreateProductsCategories, 0)
			for path, productset := range relationships.Rels {
				for _, sku := range productset.Products {
					c := eclient.CreateProductsCategories{
						CategoryID: categoryPathToID[path],
						ProductID:  productSKUToID[sku],
					}
					rels = append(rels, &c)
				}
			}

			err = client.UpdateProductCategoryRelations(rels)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
		},
	}
	return cmd
}
