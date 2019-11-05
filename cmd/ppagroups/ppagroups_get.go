package ppagroups

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdPPAGroupsGet returns new initialized instance of the get sub command
func NewCmdPPAGroupsGet() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "get <code>",
		Short: "Get an individual product to product associations group",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			// code
			code := args[0]
			ctx := context.Background()
			groups, err := client.GetPPAGroups(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			var ppaGroupID string
			for _, g := range groups {
				if g.Code == code {
					ppaGroupID = g.ID
					break
				}
			}
			if ppaGroupID == "" {
				fmt.Fprintf(os.Stderr, "product to product associations group with code %q not found\n", code)
				os.Exit(1)
			}

			group, err := client.GetPPAGroup(ctx, ppaGroupID)
			if err == eclient.ErrPPAssocGroupNotFound {
				fmt.Fprintf(os.Stderr,
					"Product to product associations group %q not found.\n", ppaGroupID)
				os.Exit(0)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			format := "%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "PPA Group ID:", group.ID)
			fmt.Fprintf(tw, format, "Code:", group.Code)
			fmt.Fprintf(tw, format, "Name:", group.Name)
			fmt.Fprintf(tw, format, "Created:",
				group.Created.In(location).Format(timeDisplayFormat))
			fmt.Fprintf(tw, format, "Modified:",
				group.Modified.In(location).Format(timeDisplayFormat))
			tw.Flush()
		},
	}
	return cmd
}
