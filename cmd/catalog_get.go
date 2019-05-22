package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

var catalogGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the catalog hierarchy",
	Run: func(cmd *cobra.Command, args []string) {
		current := rc.Configurations[currentConfigName]
		client := eclient.New(current.Endpoint, timeout)
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

func init() {
	catalogCmd.AddCommand(catalogGetCmd)
}
