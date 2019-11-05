package users

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// NewCmdUsersCreate returns new initialized instance of create sub command
func NewCmdUsersCreate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create a user",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			req, err := promptCreateUser(client)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			ctx := context.Background()
			user, err := client.CreateUser(ctx, req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error creating user: %v\n", err.Error())
				os.Exit(1)
			}

			format := "%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, format, "User ID:", user.ID)
			fmt.Fprintf(tw, format, "UID:", user.UID)
			fmt.Fprintf(tw, format, "Role:", user.Role)
			fmt.Fprintf(tw, format, "Price List ID:", user.PriceListID)
			fmt.Fprintf(tw, format, "Email:", user.Email)
			fmt.Fprintf(tw, format, "Firstname:", user.Firstname)
			fmt.Fprintf(tw, format, "Lastname:", user.Lastname)
			fmt.Fprintf(tw, format, "Created:",
				user.Created.In(location).Format(timeDisplayFormat))
			fmt.Fprintf(tw, format, "Modified:",
				user.Modified.In(location).Format(timeDisplayFormat))
			tw.Flush()
		},
	}
	return cmd
}

func promptCreateUser(client *eclient.EcomClient) (*eclient.CreateUserRequest, error) {
	var req eclient.CreateUserRequest

	// role
	r := &survey.Select{
		Message: "Role:",
		Options: []string{"user", "admin"},
	}
	survey.AskOne(r, &req.Role, nil)

	// email
	e := &survey.Input{
		Message: "Email:",
	}
	survey.AskOne(e, &req.Email, nil)

	// password
	t := &survey.Password{
		Message: "Password:",
	}
	survey.AskOne(t, &req.Password, nil)

	// firstname
	f := &survey.Input{
		Message: "Firstname:",
	}
	survey.AskOne(f, &req.Firstname, nil)

	// lastname
	l := &survey.Input{
		Message: "Lastname:",
	}
	survey.AskOne(l, &req.Lastname, nil)

	return &req, nil
}
