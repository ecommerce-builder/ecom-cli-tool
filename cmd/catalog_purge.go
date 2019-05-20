package cmd

import (
	"log"
	"os"

	"bitbucket.org/andyfusniakteam/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

var catalogPurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge the entire catalog",
	Run: func(cmd *cobra.Command, args []string) {
		current := rc.Configurations[currentConfigName]
		client := eclient.New(current.Endpoint, timeout)
		err := client.SetToken(&current)
		if err != nil {
			log.Fatal(err)
		}

		err = client.PurgeCatalog()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	},
}

func init() {
	catalogCmd.AddCommand(catalogPurgeCmd)
}
