package cmd

import (
	"fmt"
	"log"
	"os"

	"bitbucket.org/andyfusniakteam/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

var assocsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all catalog associations",
	Run: func(cmd *cobra.Command, args []string) {
		current := rc.Configurations[currentConfigName]
		client := eclient.NewEcomClient(current.FirebaseAPIKey, current.Endpoint, timeout)
		err := client.SetToken(&current)
		if err != nil {
			log.Fatal(err)
		}

		assocs, err := client.GetCatalogAssocs()
		if err != nil {
			log.Fatal(err)
		}

		for k, assoc := range assocs {
			fmt.Printf("%s:\n", k)
			for _, p := range assoc {
				fmt.Printf("\t%s\n", p.SKU)
			}
		}
		os.Exit(0)
	},
}

func init() {
	assocsCmd.AddCommand(assocsListCmd)
}
