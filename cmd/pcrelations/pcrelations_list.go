package pcrelations

import (
	"fmt"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdPCRelationsList returns new initialized instance of list sub command
func NewCmdPCRelationsList() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List all product to category relations",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			pcrelations, err := client.GetProductCategoryRelations()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			categoryPathToProductList := make(map[string][]string)

			for _, rel := range pcrelations {
				categoryPathToProductList[rel.CategoryPath] = append(categoryPathToProductList[rel.CategoryPath], rel.ProductSKU)
			}

			// Display the associations in order
			fmt.Println("product_category_relations:")
			for categoryPath := range categoryPathToProductList {
				fmt.Printf("  %s:\n", categoryPath)
				fmt.Println("    products:")
				for _, sku := range categoryPathToProductList[categoryPath] {
					fmt.Printf("      - %s\n", sku)
				}
			}
			os.Exit(0)
		},
	}
	return cmd
}
