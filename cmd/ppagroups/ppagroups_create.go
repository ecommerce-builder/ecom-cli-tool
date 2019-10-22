package ppagroups

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// NewCmdPPAGroupCreate returns new initialized instance of create sub command
func NewCmdPPAGroupCreate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new product to product associations group",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}

			// get the url and event list
			req, err := promptCreatePPAGroup()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			ctx := context.Background()
			ppaGroup, err := client.CreatePPAGroup(ctx, req)
			if err != nil {
				fmt.Printf("%+v\n", err)
				fmt.Fprintf(os.Stderr, "error creating developer key: %+v", errors.Unwrap(err))
				os.Exit(1)
			}

			format := "%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "Product to product group ID:", ppaGroup.ID)
			fmt.Fprintf(tw, format, "Code:", ppaGroup.Code)
			fmt.Fprintf(tw, format, "Name:", ppaGroup.Name)
			fmt.Fprintf(tw, format, "Created:", ppaGroup.Created)
			fmt.Fprintf(tw, format, "Modified:", ppaGroup.Modified)
			tw.Flush()
		},
	}
	return cmd
}

func promptCreatePPAGroup() (*eclient.CreatePAGroupRequest, error) {
	var req eclient.CreatePAGroupRequest

	// pp_assoc_group_code
	g := &survey.Input{
		Message: "Product to product association group code:",
	}
	survey.AskOne(g, &req.PPAssocGroupCode, nil)

	// name
	n := &survey.Input{
		Message: "Group name:",
	}
	survey.AskOne(n, &req.Name, nil)

	return &req, nil
}
