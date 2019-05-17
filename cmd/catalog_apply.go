package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"bitbucket.org/andyfusniakteam/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

var catalogApplyCmd = &cobra.Command{
	Use:   "apply <catalog.yaml>",
	Short: "Create or update the shop catalog.",
	Long:  ``,

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		current := rc.Configurations[currentConfigName]
		client := eclient.New(current.Endpoint, timeout)
		err := client.SetToken(&current)
		if err != nil {
			log.Fatal(err)
		}

		data, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		var catalog eclient.Catalog
		err = yaml.Unmarshal([]byte(data), &catalog)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		err = client.UpdateCatalog(catalog.Category)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	catalogCmd.AddCommand(catalogApplyCmd)
}
