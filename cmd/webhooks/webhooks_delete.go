package webhooks

import (
	"context"
	"fmt"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdWebhooksDelete returns new initialized instance of the delete sub command
func NewCmdWebhooksDelete() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "delete <webhook_id>",
		Short: "Delete webhook",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			webhookID := args[0]
			if !cmdvalidate.IsValidUUID(webhookID) {
				fmt.Fprintf(os.Stderr, "webhook_id %s is not a valid v4 uuid\n",
					webhookID)
				os.Exit(1)
			}

			ctx := context.Background()
			err := client.DeleteWebhook(ctx, webhookID)
			if err == eclient.ErrWebhookNotFound {
				fmt.Fprintf(os.Stderr, "webhook_id %s not found\n",
					webhookID)
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "delete webhook failed: %+v\n", err)
			}
		},
	}
	return cmd
}
