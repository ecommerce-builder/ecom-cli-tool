package cmd

import (
	"log"

	"bitbucket.org/andyfusniakteam/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

var adminDeleteCmd = &cobra.Command{
	Use:   "remove <uuid>",
	Short: "Remove an administrator",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		current := rc.Configurations[currentConfigName]
		client := eclient.New(current.Endpoint, timeout)
		err := client.SetToken(&current)
		if err != nil {
			log.Fatal(err)
		}

		uuid := args[0]
		err = client.DeleteAdmin(uuid)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	adminsCmd.AddCommand(adminDeleteCmd)
}
