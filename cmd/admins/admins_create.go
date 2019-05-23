package admins

import (
	"fmt"
	"log"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"

	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// NewCmdAdminsCreate returns new initialized instance of create sub command
func NewCmdAdminsCreate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		log.Fatal(err)
	}
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new administrator",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			email, passwd, first, last, err := promptCreateAdmin()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			client := eclient.New(current.Endpoint)
			err = client.SetToken(&current)
			if err != nil {
				log.Fatal(err)
			}
			customer, err := client.CreateAdmin(email, passwd, first, last)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			fmt.Println("Adminstrator created")
			fmt.Println("--------------------")
			fmt.Printf("UUID: %s\n", customer.UUID)
			fmt.Printf("UID: %s\n", customer.UID)
			fmt.Printf("Role: %s\n", customer.Role)
			fmt.Printf("Firstname: %s\n", customer.Firstname)
			fmt.Printf("Lastname: %s\n", customer.Lastname)
			fmt.Printf("Created: %s\n", customer.Created)
		},
	}
	return cmd
}

func promptCreateAdmin() (email, passwd, first, last string, err error) {
	ep := &survey.Input{
		Message: "Email:",
	}
	survey.AskOne(ep, &email, nil)

	pp := &survey.Password{
		Message: "Password:",
	}
	survey.AskOne(pp, &passwd, nil)

	pf := &survey.Input{
		Message: "Firstname:",
	}
	survey.AskOne(pf, &first, nil)

	pl := &survey.Input{
		Message: "Lastname:",
	}
	survey.AskOne(pl, &last, nil)
	return email, passwd, first, last, nil
}
