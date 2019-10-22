package devkeys

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdDevKeysCreate returns new initialized instance of create sub command
func NewCmdDevKeysCreate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "create <email>",
		Short: "Create a new developer key for a given user",
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

			req := &eclient.DevKeyRequest{
				UserID: userMap[email],
			}
			devKey, err := client.CreateDeveloperKey(ctx, req)
			if err != nil {
				fmt.Printf("%+v\n", err)
				fmt.Fprintf(os.Stderr, "error creating developer key: %+v", errors.Unwrap(err))
			}

			format := "%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "Developer Key ID", devKey.ID)
			fmt.Fprintf(tw, format, "User ID", devKey.UserID)
			fmt.Fprintf(tw, format, "Private Key", devKey.Key)
			fmt.Fprintf(tw, format, "Created", devKey.Created)
			fmt.Fprintf(tw, format, "Modified", devKey.Modified)
			fmt.Fprintf(tw, format, "", "")
			tw.Flush()
		},
	}
	return cmd
}
