package devkeys

import (
	"context"
	"fmt"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdDevKeysDelete returns new initialized instance of the delete sub command
func NewCmdDevKeysDelete() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "delete <coupon_code>",
		Short: "Delete a coupon",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			// developer_key_id
			devKeyID := args[0]
			if !cmdvalidate.IsValidUUID(devKeyID) {
				fmt.Fprintf(os.Stderr, "developer_key_id %q is not a valid v4 uuid\n", devKeyID)
				os.Exit(1)
			}

			ctx := context.Background()
			err = client.DeleteDeveloperKey(ctx, devKeyID)
			if err == eclient.ErrDeveloperKeyNotFound {
				fmt.Fprintf(os.Stderr, "developer key not found. Use ecom devkeys list to check.\n")
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
