package ppagroups

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdPPAGroupsDelete returns new initialized instance of the delete sub command
func NewCmdPPAGroupsDelete() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "delete <code>",
		Short: "Delete a product to product associations group by code",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
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
				fmt.Fprintf(os.Stderr,
					"product to product associations group code %q not found\n", code)
				os.Exit(1)
			}

			err = client.DeletePPAGroup(ctx, ppaGroupID)
			if err == eclient.ErrBadRequest {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			if err == eclient.ErrPPAssocGroupNotFound {
				fmt.Fprintf(os.Stderr, "product to product associations group not found. Use ecom ppagroups list to check.\n")
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
		},
	}
	return cmd
}
