package admins

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdAdminsList returns new initialized instance of list sub command
func NewCmdAdminsList() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List all administrators",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			admins, err := client.ListAdmins()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			format := "%v\t%v\t%v\t%v\t%v\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "Name", "Email", "UID", "UUID", "Created")
			fmt.Fprintf(tw, format, "----", "-----", "---", "----", "-------")
			for _, admin := range admins {
				fmt.Fprintf(tw, format, admin.Firstname+" "+admin.Lastname, admin.Email, admin.UID, admin.UUID, admin.Created)
			}
			tw.Flush()
			os.Exit(0)
		},
	}
	return cmd
}
