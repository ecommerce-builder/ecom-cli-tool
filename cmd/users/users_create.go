package users

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

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
				log.Fatal(err)
			}

			req, err := promptCreateUser(client)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			fmt.Println(req.Role)
			fmt.Println(req.Email)
			fmt.Println(req.Password)
			fmt.Println(req.Firstname)
			fmt.Println(req.Lastname)

			ctx := context.Background()
			user, err := client.CreateUser(ctx, req)
			if err != nil {
				fmt.Printf("%+v\n", err)
				fmt.Fprintf(os.Stderr, "error creating user: %+v", errors.Unwrap(err))
			}

			fmt.Println(user)
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
