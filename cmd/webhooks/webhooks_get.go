package webhooks

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
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

// NewCmdWebhooksGet returns new initialized instance of the get sub command
func NewCmdWebhooksGet() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "get <webhook_id>",
		Short: "Get a webhook",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}

			webhookID := args[0]
			if !cmdvalidate.IsValidUUID(webhookID) {
				fmt.Fprintf(os.Stderr, "webhook_id %s is not a valid v4 uuid\n", webhookID)
				os.Exit(1)
			}

			ctx := context.Background()
			webhook, err := client.GetWebhook(ctx, webhookID)
			if err == eclient.ErrWebhookNotFound {
				fmt.Fprintf(os.Stderr, "webhook %s not found\n", webhookID)
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			format := "%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "Webhook ID:", webhook.ID)
			fmt.Fprintf(tw, format, "Signing Key:", webhook.SigningKey)
			fmt.Fprintf(tw, format, "URL:", webhook.URL)
			fmt.Fprintf(tw, format, "Events:", webhook.Events)
			fmt.Fprintf(tw, format, "Enabled:", webhook.Enabled)
			fmt.Fprintf(tw, format, "Created:", webhook.Created.In(location).Format(timeDisplayFormat))
			fmt.Fprintf(tw, format, "Modified:", webhook.Modified.In(location).Format(timeDisplayFormat))
			tw.Flush()
		},
	}
	return cmd
}
