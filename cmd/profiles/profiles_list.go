package profiles

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/spf13/cobra"
)

// NewCmdProfilesList returns new initialized instance of list sub command
func NewCmdProfilesList() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "list",
		Short: "Display a list of available profiles",
		Run: func(cmd *cobra.Command, args []string) {
			if len(cfgs.Configurations) == 0 {
				fmt.Println("No profiles")
				os.Exit(0)
			}
			format := "%v\t%v\t%v\t%v\t%v\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "Active", "Endpoint", "Email", "Role", "Dev Key")
			fmt.Fprintf(tw, format, "------", "--------", "-----", "----", "-------")
			for k, v := range cfgs.Configurations {
				var active string
				if curCfg == k {
					active = "  *"
				} else {
					active = ""
				}
				fmt.Fprintf(tw, format, active, v.Endpoint, v.Customer.Email, v.Customer.Role, v.DevKey[0:5]+"********")
			}
			tw.Flush()
		},
	}
	return cmd
}
