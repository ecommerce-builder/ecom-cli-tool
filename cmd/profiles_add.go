package cmd

import (
	"fmt"
	"os"

	"bitbucket.org/andyfusniakteam/ecom-cli-tool/configmgr"
	"bitbucket.org/andyfusniakteam/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// profilesListCmd represents the profilesList command
var profilesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new profile",
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
		customToken, err := ecomClient.SignInWithDevKey(devKey)
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
		rc.Configurations[filename] = configmgr.EcomConfigEntry{
			DevKey:   devKey,
			Endpoint: endpoint,
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
	profilesCmd.AddCommand(profilesAddCmd)
}
