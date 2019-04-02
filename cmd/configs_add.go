package cmd

import (
	"fmt"
	"os"

	"bitbucket.org/andyfusniakteam/ecom-cli-tool/configmgr"
	"bitbucket.org/andyfusniakteam/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// projectsListCmd represents the projectsList command
var projectsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new configuration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		name, webKey, endpoint, devKey, err := promptAddConfiguration()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		ecomClient := eclient.NewEcomClient(webKey, endpoint, timeout)
		customToken, err := ecomClient.SignInWithDevKey(devKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		tar, err := ecomClient.ExchangeCustomTokenForIDAndRefreshToken(customToken)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}

		configmgr.WriteTokenAndRefreshToken(webKey, endpoint, tar)

		if rc.Configurations == nil {
			rc.Configurations = make(map[string]configmgr.EcomConfigEntry)
		}

		rc.Configurations[name] = configmgr.EcomConfigEntry{
			FirebaseAPIKey: webKey,
			DevKey:         devKey,
			Endpoint:       endpoint,
		}

		err = configmgr.WriteConfig(rc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "write config failed: %+v", err)
			os.Exit(1)
		}
	},
}

func promptAddConfiguration() (name, webKey, endpoint, devKey string, err error) {
	n := &survey.Input{
		Help:    "Choose a name for this configuration",
		Message: "Name:",
	}
	survey.AskOne(n, &name, nil)

	k := &survey.Input{
		Message: "Web API Key:",
	}
	survey.AskOne(k, &webKey, nil)

	e := &survey.Input{
		Message: "Endpoint:",
	}
	survey.AskOne(e, &endpoint, nil)

	d := &survey.Input{
		Message: "Developer Key:",
	}
	survey.AskOne(d, &devKey, nil)
	return name, webKey, endpoint, devKey, nil
}

func init() {
	projectsCmd.AddCommand(projectsAddCmd)
}
