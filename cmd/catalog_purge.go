package cmd

import (
	"log"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

var catalogPurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge the entire catalog",
	Run: func(cmd *cobra.Command, args []string) {
		current := rc.Configurations[currentConfigName]
		client := eclient.New(current.Endpoint, timeout)
		if err := client.SetToken(&current); err != nil {
			log.Fatal(err)
		}
		if err := client.PurgeCatalog(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	},
}

func init() {
	catalogCmd.AddCommand(catalogPurgeCmd)
}
