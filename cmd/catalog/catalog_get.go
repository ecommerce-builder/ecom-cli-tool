package catalog

import (
	"fmt"
	"log"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdCatalogGet returns new initialized instance of list get command
func NewCmdCatalogGet() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Get the catalog hierarchy",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}
			root, err := client.GetCatalog()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			treeView(root, 0, false)
		},
	}
	return cmd
}

func treeView(node *eclient.Category, depth int, lastSibling bool) {
	// fmt.Printf("%+v\n", node)
	// fmt.Printf("node.Name=%s last sibling=%t\n", node.Name, lastSibling)
	var arm string
	if lastSibling {
		arm = "└── "
	} else {
		arm = "├── "
	}

	if depth == 0 {
	} else if depth == 1 {
		fmt.Print(arm)
	} else {
		fmt.Print("│   ")
		for i := 0; i < depth-2; i++ {
			fmt.Print("    ")
		}
		fmt.Print(arm)
	}
	fmt.Printf("%s (%s)\n", node.Segment, node.Name)
	lastIdx := len(node.Categories) - 1
	for i, n := range node.Categories {
		treeView(n, depth+1, lastIdx == i)
	}
}
