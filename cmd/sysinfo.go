package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"bitbucket.org/andyfusniakteam/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// sysinfoCmd represents the sysinfo command
var sysinfoCmd = &cobra.Command{
	Use:   "sysinfo",
	Short: "Prints system information from the running API service.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		current := rc.Configurations[currentConfigName]
		ecomClient := eclient.NewEcomClient(current.FirebaseAPIKey, current.Endpoint, timeout)
		err := ecomClient.SetToken(&current)
		if err != nil {
			log.Fatal(err)
		}

		sysInfo, err := ecomClient.SysInfo()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		format := "%v\t%v\t\n"
		tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprintf(tw, format, "API Service", "")
		fmt.Fprintf(tw, format, "-----------", "")
		fmt.Fprintf(tw, format, "Version", sysInfo.APIVersion)
		fmt.Fprintf(tw, format, "", "")
		fmt.Fprintf(tw, format, "Postgres", "")
		fmt.Fprintf(tw, format, "--------", "")
		fmt.Fprintf(tw, format, "Host", sysInfo.Env.Pg.Host)
		fmt.Fprintf(tw, format, "Port", sysInfo.Env.Pg.Port)
		fmt.Fprintf(tw, format, "Database", sysInfo.Env.Pg.Database)
		fmt.Fprintf(tw, format, "User", sysInfo.Env.Pg.User)
		fmt.Fprintf(tw, format, "SSLMode", sysInfo.Env.Pg.SslMode)
		fmt.Fprintf(tw, format, "", "")
		fmt.Fprintf(tw, format, "Google", "")
		fmt.Fprintf(tw, format, "------", "")
		fmt.Fprintf(tw, format, "Project ID", sysInfo.Env.Goog.ProjectID)
		fmt.Fprintf(tw, format, "Web API Key", sysInfo.Env.Goog.WebAPIKey)
		fmt.Fprintf(tw, format, "", "")
		fmt.Fprintf(tw, format, "App", "")
		fmt.Fprintf(tw, format, "---", "")
		fmt.Fprintf(tw, format, "HTTP Port", sysInfo.Env.App.HTTPPort)
		fmt.Fprintf(tw, format, "Root Email", sysInfo.Env.App.RootEmail)
		tw.Flush()
	},
}

func init() {
	rootCmd.AddCommand(sysinfoCmd)
}
