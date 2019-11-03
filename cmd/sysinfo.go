package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdSysInfo returns new initialized instance of the sysinfo sub command
func NewCmdSysInfo() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "sysinfo",
		Short: "Prints system information from the running API service.",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			ecomClient := eclient.New(current.Endpoint)
			if err := ecomClient.SetToken(&current); err != nil {
				log.Fatal(err)
			}
			sysInfo, err := ecomClient.SysInfo()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			format := "%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)

			fmt.Fprintf(tw, format, "Ecom CLI Tool", "")
			fmt.Fprintf(tw, format, "-------------", "")
			fmt.Fprintf(tw, format, "Client Version", Version)
			fmt.Fprintf(tw, format, "Endpoint", current.Endpoint)
			fmt.Fprintf(tw, format, "", "")

			fmt.Fprintf(tw, format, "API Service", "")
			fmt.Fprintf(tw, format, "-----------", "")
			fmt.Fprintf(tw, format, "Service Version", sysInfo.APIVersion)
			fmt.Fprintf(tw, format, "", "")
			fmt.Fprintf(tw, format, "Postgres", "")
			fmt.Fprintf(tw, format, "--------", "")
			fmt.Fprintf(tw, format, "Schema Version", sysInfo.Env.Pg.SchemaVersion)
			fmt.Fprintf(tw, format, "Host", sysInfo.Env.Pg.Host)
			fmt.Fprintf(tw, format, "Port", sysInfo.Env.Pg.Port)
			fmt.Fprintf(tw, format, "Database", sysInfo.Env.Pg.Database)
			fmt.Fprintf(tw, format, "User", sysInfo.Env.Pg.User)
			fmt.Fprintf(tw, format, "SSLMode", sysInfo.Env.Pg.SslMode)
			fmt.Fprintf(tw, format, "", "")
			fmt.Fprintf(tw, format, "Google", "")
			fmt.Fprintf(tw, format, "------", "")
			fmt.Fprintf(tw, format, "GAE Project ID", sysInfo.Env.Goog.ProjectID)
			fmt.Fprintf(tw, format, "", "")
			fmt.Fprintf(tw, format, "Firebase", "")
			fmt.Fprintf(tw, format, "--------", "")
			fmt.Fprintf(tw, format, "Firebase Project ID", sysInfo.Env.Firebase.ProjectID)
			fmt.Fprintf(tw, format, "Firebase Web API Key", sysInfo.Env.Firebase.APIKEY)
			fmt.Fprintf(tw, format, "", "")
			fmt.Fprintf(tw, format, "Stripe", "")
			fmt.Fprintf(tw, format, "--------", "")
			fmt.Fprintf(tw, format, "Stripe Success URL", sysInfo.Env.Stripe.StripeSuccessURL)
			fmt.Fprintf(tw, format, "Stripe Cancel URL", sysInfo.Env.Stripe.StripeCancelURL)
			fmt.Fprintf(tw, format, "", "")
			fmt.Fprintf(tw, format, "App", "")
			fmt.Fprintf(tw, format, "---", "")
			fmt.Fprintf(tw, format, "HTTP Port", sysInfo.Env.App.AppPort)
			fmt.Fprintf(tw, format, "Root Email", sysInfo.Env.App.AppRootEmail)
			fmt.Fprintf(tw, "%v\t%t\t\n", "Stackdriver Logging Enabled:", sysInfo.Env.App.AppEnableStackDriverLogging)
			fmt.Fprintf(tw, format, "Endpoint:", sysInfo.Env.App.AppEndpoint)
			tw.Flush()
		},
	}
	return cmd
}
