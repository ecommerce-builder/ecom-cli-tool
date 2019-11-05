package webhooks

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdWebhooksList returns new initialized instance of the list sub command
func NewCmdWebhooksList() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "list",
		Short: "list webhooks",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			ctx := context.Background()
			webhooks, err := client.GetWebhooks(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			format := "%s\t%s\t%s\t%v\t%v\t%v\t%v\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "Webhook ID", "Signing Key", "URL",
				"Events", "Enabled", "Created", "Modified")
			fmt.Fprintf(tw, format, "----------", "-----------", "---",
				"------", "-------", "-------", "--------")

			for _, w := range webhooks {
				fmt.Fprintf(tw, format,
					w.ID, w.SigningKey, w.URL,
					w.Events, w.Enabled,
					w.Created.In(location).Format(timeDisplayFormat),
					w.Modified.In(location).Format(timeDisplayFormat))
			}
			tw.Flush()
		},
	}
	return cmd
}
