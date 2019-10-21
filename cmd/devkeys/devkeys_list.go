package devkeys

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

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

// NewCmdDevKeysList returns new initialized instance of list sub command
func NewCmdDevKeysList() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "list <user_id>",
		Short: "List developer keys",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}

			email := args[0]
			ctx := context.Background()
			users, err := client.GetUsers(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			userMap := make(map[string]string, 0)
			var userOpts []string
			for _, user := range users {
				userOpts = append(userOpts, user.Email)
				userMap[user.Email] = user.ID
			}

			if _, ok := userMap[email]; !ok {
				fmt.Fprintf(os.Stderr, "Email %s did not match any users\n", email)
				os.Exit(1)
			}

			devKeys, err := client.GetDeveloperKeys(ctx, userMap[email])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			format := "%s\t%s\t%s\t%v\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "Developer Key ID", "Key", "Created")
			fmt.Fprintf(tw, format, "----------------", "---", "-------")

			for _, devKey := range devKeys {
				fmt.Fprintf(tw, format,
					devKey.ID, devKey.Key,
					devKey.Created.In(location).Format(timeDisplayFormat))
			}

			tw.Flush()
		},
	}
	return cmd
}
