package ppagroups

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

const timeDisplayFormat = "2006-01-02 15:04"

var location *time.Location

func init() {
	var err error
	location, err = time.LoadLocation("Europe/London")
	if err != nil {
		fmt.Fprintf(os.Stderr, "time.LoadLocation(%q) failed: %+v", "Europe/London", err.Error())
		return
	}
}

// NewCmdPPAGroupsList returns new initialized instance of the list sub command
func NewCmdPPAGroupsList() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "list",
		Short: "list product to product associations",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			ctx := context.Background()
			ppaGroups, err := client.GetPPAGroups(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			fmt.Println(ppaGroups)

			format := "%s\t%s\t%s\t%v\t%v\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "PP Assoc Group ID", "Code", "Name",
				"Created", "Modified")
			fmt.Fprintf(tw, format, "-----------------", "----", "----",
				"-------", "--------")
			for _, g := range ppaGroups {
				fmt.Fprintf(tw, format,
					g.ID, g.Code, g.Name,
					g.Created.In(location).Format(timeDisplayFormat),
					g.Modified.In(location).Format(timeDisplayFormat))
			}
			tw.Flush()
		},
	}
	return cmd
}
