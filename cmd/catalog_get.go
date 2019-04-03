package cmd

import (
	"fmt"
	"log"
	"os"

	"bitbucket.org/andyfusniakteam/ecom-api-go/utils/nestedset"
	"bitbucket.org/andyfusniakteam/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

var catalogGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the catalog hierarchy",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		current := rc.Configurations[currentConfigName]
		ecomClient := eclient.NewEcomClient(current.FirebaseAPIKey, current.Endpoint, timeout)
		err := ecomClient.SetToken(&current)
		if err != nil {
			log.Fatal(err)
		}

		nodes, err := ecomClient.GetCatalog()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		root := nestedset.BuildTree(nodes)

		root.PreorderTraversalPrint(os.Stdout)

	},
}

func init() {
	catalogCmd.AddCommand(catalogGetCmd)
}
