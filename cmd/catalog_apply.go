package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func isValidEndpoint(ep string, valid []string) (bool, error) {
	if len(valid) == 0 {
		return true, nil
	}
	url, err := url.Parse(ep)
	if err != nil {
		return false, err
	}
	ephost := url.Hostname()
	for _, s := range valid {
		if s == ephost {
			return true, nil
		}
	}
	return false, nil
}

var catalogApplyCmd = &cobra.Command{
	Use:   "apply <catalog.yaml>",
	Short: "Create or update the shop catalog.",

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		current := rc.Configurations[currentConfigName]
		client := eclient.New(current.Endpoint, timeout)
		if err := client.SetToken(&current); err != nil {
			log.Fatal(err)
		}

		data, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		var catalog eclient.Catalog
		if err = yaml.Unmarshal([]byte(data), &catalog); err != nil {
			log.Fatalf("error: %v", err)
		}

		// disallow applying the catalog.yaml files with endpoints: ['host1', 'host2']
		// guards to the system.
		ok, err := isValidEndpoint(current.Endpoint, catalog.Endpoints)
		if err != nil {
			log.Fatal(err)
		}
		if !ok {
			fmt.Fprintf(os.Stderr, "The catalog.yaml file has endpoint guards for %v only. Your current profile endpoint is %q. Either switch profiles using 'ecom profiles select' or adjust the catalog.yaml file.\n", catalog.Endpoints, current.Endpoint)
			os.Exit(1)
		}

		if err = client.UpdateCatalog(catalog.Category); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	catalogCmd.AddCommand(catalogApplyCmd)
}
