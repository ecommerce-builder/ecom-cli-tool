package webhooks

import (
	"context"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

const (
	// EventServiceStarted event
	EventServiceStarted string = "service.started"

	// EventAddressCreated event
	EventAddressCreated string = "address.created"

	// EventAddressUpdated event
	EventAddressUpdated string = "address.updated"

	// EventUserCreated event
	EventUserCreated string = "user.created"

	// EventOrderCreated triggerred after an order has been placed.
	EventOrderCreated string = "order.created"

	// EventOrderUpdated event
	EventOrderUpdated string = "order.updated"
)

// NewCmdWebhooksCreate returns new initialized instance of create sub command
func NewCmdWebhooksCreate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create a webhook",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			// get the url and event list
			req, err := promptCreateWebhook()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			// attempt to create the webhook
			ctx := context.Background()
			webhook, err := client.CreateWebhook(ctx, req)
			if err == eclient.ErrEventTypeNotFound {
				fmt.Fprint(os.Stderr, "one or more of the events are not known\n")
				os.Exit(1)
			}
			if err == eclient.ErrWebhookExists {
				fmt.Fprintf(os.Stderr, "webhook with this URL of %s already exists\n",
					req.URL)
				os.Exit(1)
			}
			if err != nil {
				fmt.Printf("%+v\n", err)
				fmt.Fprintf(os.Stderr, "error creating user: %+v\n",
					errors.Unwrap(err))
				os.Exit(1)
			}

			format := "%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "ID:", webhook.ID)
			fmt.Fprintf(tw, format, "Signing Key:", webhook.SigningKey)
			fmt.Fprintf(tw, format, "URL:", webhook.URL)
			fmt.Fprintf(tw, format, "Events:", webhook.Events)
			fmt.Fprintf(tw, format, "Enabled:", webhook.Enabled)
			fmt.Fprintf(tw, format, "Created:", webhook.Created)
			fmt.Fprintf(tw, format, "Modified:", webhook.Modified)
			tw.Flush()
		},
	}
	return cmd
}

func promptCreateWebhook() (*eclient.CreateWebhookRequest, error) {
	var req eclient.CreateWebhookRequest

	// url
	u := &survey.Input{
		Message: "URL:",
	}
	survey.AskOne(u, &req.URL, nil)

	// events
	prompt := &survey.MultiSelect{
		Message: "Events:",
		Options: []string{
			EventServiceStarted,
			EventAddressCreated,
			EventAddressUpdated,
			EventUserCreated,
			EventOrderCreated,
		},
	}
	survey.AskOne(prompt, &req.Events.Data, nil)

	return &req, nil
}
