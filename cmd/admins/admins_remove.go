package admins

import (
	"log"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdAdminsRemove returns new initialized instance of remove sub command
func NewCmdAdminsRemove() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		log.Fatal(err)
	}
	var cmd = &cobra.Command{
		Use:   "remove <uuid>",
		Short: "Remove an administrator",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}
			uuid := args[0]
			if err := client.DeleteAdmin(uuid); err != nil {
				log.Fatal(err)
			}
		},
	}
	return cmd
}
