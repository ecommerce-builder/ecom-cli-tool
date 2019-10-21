package users

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

// NewCmdUsersList returns new initialized instance of list sub command
func NewCmdUsersList() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List users",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}

			ctx := context.Background()
			users, err := client.ListUsers(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			format := "%s\t%s\t%s\t%s\t%s\t%s\t%v\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "User ID", "UID", "Role", "Email", "Firstname", "Lastname", "Created")
			fmt.Fprintf(tw, format, "-------", "---", "----", "-----", "---------", "--------", "-------")

			for _, user := range users {
				fmt.Fprintf(tw, format,
					user.ID, user.UID, user.Role, user.Email, user.Firstname, user.Lastname, user.Created.In(location).Format(timeDisplayFormat))
			}

			tw.Flush()
		},
	}
	return cmd
}
