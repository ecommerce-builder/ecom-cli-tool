package pcrelations

import (
	"fmt"
	"log"
	"os"
	"sort"

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
				log.Fatal(err)
			}

			pcrelations, err := client.GetCatalogAssocs()
			if err != nil {
				log.Fatal(err)
			}

			// To store the paths in slice in sorted order
			var paths []string
			for p := range pcrelations {
				paths = append(paths, p)
			}
			sort.Strings(paths)

			// Display the associations in order
			fmt.Println("associations:")
			for _, path := range paths {
				fmt.Printf("  %s:\n", path)
				fmt.Println("    products:")
				for _, p := range pcrelations[path].Products {
					fmt.Printf("      - %s\n", p.SKU)
				}
			}
			os.Exit(0)
		},
	}
	return cmd
}
