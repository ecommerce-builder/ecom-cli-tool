package assocs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdAssocsApply returns new initialized instance of apply sub command
func NewCmdAssocsApply() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "apply <assocs.yaml>",
		Short: "Create or update an exising catalog association set.",
		Long:  ``,

		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
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
	return cmd
}
