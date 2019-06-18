package assocs

import (
	"fmt"
	"log"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdAssocsList returns new initialized instance of list sub command
func NewCmdAssocsList() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List all catalog associations",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}

			assocs, err := client.GetCatalogAssocs()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("associations:")
			for path, assoc := range assocs {
				fmt.Printf("  %s:\n", path)
				fmt.Println("    products:")
				for _, p := range assoc.Products {
					fmt.Printf("      - %s\n", p.SKU)
				}
			}
			os.Exit(0)
		},
	}
	return cmd
}
