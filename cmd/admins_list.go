package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"bitbucket.org/andyfusniakteam/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

var adminsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all administrators",
	Run: func(cmd *cobra.Command, args []string) {
		current := rc.Configurations[currentConfigName]
		client := eclient.New(current.Endpoint, timeout)
		err := client.SetToken(&current)
		if err != nil {
			log.Fatal(err)
		}
		admins, err := client.ListAdmins()
		if err != nil {
			log.Fatal(err)
		}
		format := "%v\t%v\t%v\t%v\t%v\n"
		tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprintf(tw, format, "Name", "Email", "UID", "UUID", "Created")
		fmt.Fprintf(tw, format, "----", "-----", "---", "----", "-------")
		for _, admin := range admins {
			fmt.Fprintf(tw, format, admin.Firstname+" "+admin.Lastname, admin.Email, admin.UID, admin.UUID, admin.Created)
		}
		tw.Flush()
		os.Exit(0)
	},
}

func init() {
	adminsCmd.AddCommand(adminsListCmd)
}
