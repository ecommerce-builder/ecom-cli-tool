package cmd

import (
	"fmt"
	"log"
	"os"

	"bitbucket.org/andyfusniakteam/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// customersCmd represents the customers command
var customersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List customers",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		current := rc.Configurations[currentConfigName]
		client := eclient.NewEcomClient(current.FirebaseAPIKey, current.Endpoint, timeout)
		err := client.SetToken(&current)
		if err != nil {
			log.Fatal(err)
		}

		customers, err := client.ListCustomers()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		fmt.Println(customers)
	},
}

func init() {
	customersCmd.AddCommand(customersListCmd)
}
