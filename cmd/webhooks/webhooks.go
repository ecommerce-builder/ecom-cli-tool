package webhooks

import (
	"github.com/spf13/cobra"
)

// NewCmdWebhooks returns new initialized instance of the webhooks sub command
func NewCmdWebhooks() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "webhooks",
		Short: "Webhooks Management",
	}
	cmd.AddCommand(NewCmdWebhooksCreate())
	cmd.AddCommand(NewCmdWebhooksList())
	cmd.AddCommand(NewCmdWebhooksGet())
	cmd.AddCommand(NewCmdWebhooksUpdate())
	cmd.AddCommand(NewCmdWebhooksDelete())
	return cmd
}
