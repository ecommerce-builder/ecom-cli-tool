package cmd

import (
	"fmt"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

var profilesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new profile",
	Run: func(cmd *cobra.Command, args []string) {
		endpoint, devKey, err := promptAddProfile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		client := eclient.New(endpoint, timeout)
		g, err := client.GetConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		ecomClient := eclient.New(endpoint, timeout)
		customToken, customer, err := ecomClient.SignInWithDevKey(devKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		tar, err := ecomClient.ExchangeCustomTokenForIDAndRefreshToken(g.WebAPIKey, customToken)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}

		hostname, err := configmgr.URLToHostName(endpoint)
		filename := fmt.Sprintf("%s-%s", hostname, devKey[:6])
		configmgr.WriteTokenAndRefreshToken(filename, tar)
		if rc.Configurations == nil {
			rc.Configurations = make(map[string]configmgr.EcomConfigEntry)
		}

		c := configmgr.Customer{
			UUID:      customer.UUID,
			UID:       customer.UID,
			Role:      customer.Role,
			Email:     customer.Email,
			Firstname: customer.Firstname,
			Lastname:  customer.Lastname,
		}
		rc.Configurations[filename] = configmgr.EcomConfigEntry{
			DevKey:   devKey,
			Endpoint: endpoint,
			Customer: c,
		}

		err = configmgr.WriteConfig(rc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "write config failed: %+v", err)
			os.Exit(1)
		}
	},
}

func promptAddProfile() (endpoint, devKey string, err error) {
	e := &survey.Input{
		Message: "Endpoint:",
	}
	survey.AskOne(e, &endpoint, nil)
	d := &survey.Input{
		Message: "Developer Key:",
	}
	survey.AskOne(d, &devKey, nil)
	return endpoint, devKey, nil
}

func init() {
	profilesCmd.AddCommand(profilesCreateCmd)
}
