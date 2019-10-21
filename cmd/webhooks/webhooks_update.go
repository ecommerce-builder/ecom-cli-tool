package webhooks

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// NewCmdWebhooksUpdate returns new initialized instance of the update sub command
func NewCmdWebhooksUpdate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "update <webhook_id>",
		Short: "Get a webhook",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}

			// promo_rule_code to id
			webhookID := args[0]
			if !cmdvalidate.IsValidUUID(webhookID) {
				fmt.Fprintf(os.Stderr, "webhook_id must be a valid v4 uuid")
				os.Exit(1)
			}

			ctx := context.Background()
			existingWebhook, err := client.GetWebhook(ctx, webhookID)
			if err == eclient.ErrWebhookNotFound {
				fmt.Fprintf(os.Stderr, "webhook %s not found", webhookID)
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			req, err := promptUpdateWebhook(existingWebhook.URL, existingWebhook.Events, existingWebhook.Enabled)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v", err)
				os.Exit(1)
			}

			// attempt to create the webhook
			webhook, err := client.UpdateWebhook(ctx, webhookID, req)
			if err == eclient.ErrEventTypeNotFound {
				fmt.Fprint(os.Stderr, "one or more of the events are not known")
				os.Exit(1)
			}
			if err == eclient.ErrWebhookExists {
				fmt.Fprint(os.Stderr, "webhook with this URL already exists")
				os.Exit(1)
			}
			if err != nil {
				fmt.Printf("%+v\n", err)
				fmt.Fprintf(os.Stderr, "error creating user: %+v", errors.Unwrap(err))
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

func promptUpdateWebhook(url string, events []string, enabled bool) (*eclient.UpdateWebhookRequest, error) {
	var req eclient.UpdateWebhookRequest

	// url
	u := &survey.Input{
		Message: "URL:",
		Default: url,
	}
	survey.AskOne(u, &req.URL, nil)

	// events
	e := &survey.MultiSelect{
		Message: "Events:",
		Options: []string{
			EventServiceStarted,
			EventUserCreated,
			EventOrderCreated,
		},
		Default: events,
	}
	survey.AskOne(e, &req.Events.Data, nil)

	// enabled
	n := &survey.Confirm{
		Message: "Enabled this webhook?",
		Default: enabled,
	}
	survey.AskOne(n, &req.Enabled, nil)

	return &req, nil
}
