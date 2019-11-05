package address

import (
	"context"
	"fmt"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdAddressDelete returns new initialized instance of the delete sub command
func NewCmdAddressDelete() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "delete <address_id>",
		Short: "Delete an address by id",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			addrID := args[0]
			if !cmdvalidate.IsValidUUID(addrID) {
				fmt.Fprintf(os.Stderr, "address_id %q is not a valid v4 uuid\n", addrID)
				os.Exit(1)
			}

			ctx := context.Background()
			err = client.DeleteAddress(ctx, addrID)
			if err == eclient.ErrAddressNotFound {
				fmt.Fprintf(os.Stderr, "address %q not found\n", addrID)
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
