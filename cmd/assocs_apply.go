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

var assocsApplyCmd = &cobra.Command{
	Use:   "apply <assocs.yaml>",
	Short: "Create or update an exising catalog association set.",
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

		var assocs eclient.Associations
		err = yaml.Unmarshal([]byte(data), &assocs)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		err = client.UpdateCatalogAssocs(assocs.Assocs)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	},
}

func init() {
	assocsCmd.AddCommand(assocsApplyCmd)
}
