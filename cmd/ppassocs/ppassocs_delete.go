package ppassocs

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdPPAssocsDelete returns new initialized instance of the delete sub command
func NewCmdPPAssocsDelete() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "delete <pp_assocs_id>",
		Short: "Delete a product to product associations",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}

			// code
			ppAssocID := args[0]
			if !cmdvalidate.IsValidUUID(ppAssocID) {
				fmt.Fprintf(os.Stderr, "pp_assocs_id must be a valid v4 uuid\n")
				os.Exit(1)
			}

			ctx := context.Background()
			err = client.DeletePPAssoc(ctx, ppAssocID)
			if err == eclient.ErrBadRequest {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			if err == eclient.ErrPPAssocNotFound {
				fmt.Fprintf(os.Stderr, "product to product associations not found. Use ecom ppassocs list to check.\n")
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
