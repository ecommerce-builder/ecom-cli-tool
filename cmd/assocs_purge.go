package cmd

import (
	"log"
	"os"

	"bitbucket.org/andyfusniakteam/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

var assocsPurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge all catalog associations",
	Run: func(cmd *cobra.Command, args []string) {
		current := rc.Configurations[currentConfigName]
		client := eclient.New(current.Endpoint, timeout)
		err := client.SetToken(&current)
		if err != nil {
			log.Fatal(err)
		}
		err = client.PurgeCatalogAssocs()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	},
}

func init() {
	assocsCmd.AddCommand(assocsPurgeCmd)
}
