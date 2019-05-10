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
		client := eclient.NewEcomClient(current.FirebaseAPIKey, current.Endpoint, timeout)
		err := client.SetToken(&current)
		if err != nil {
			log.Fatal(err)
		}

		data, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		type Category struct {
			URL        string     `yaml:"url"`
			Name       string     `yaml:"name"`
			Categories []Category `yaml:"categories"`
		}
		
		type Categories struct {
			Categories []Category `yaml:"categories"`
		}

		type Catalog struct {
			Catalog Categories `yaml:"catalog"`
		}

		var cat Catalog
		err = yaml.Unmarshal([]byte(data), &cat)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		for _, n := range cat.Catalog.Categories {
			fmt.Println(n)
		}
	},
}

func init() {
	catalogCmd.AddCommand(catalogApplyCmd)
}
