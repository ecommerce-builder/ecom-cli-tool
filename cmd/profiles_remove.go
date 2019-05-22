package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// profilesListCmd represents the profilesList command
var profilesRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a profile",
	Long:  `Removes a profile dropping the credentials`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(rc.Configurations) == 0 {
			fmt.Println("No profiles")
			os.Exit(0)
		}

		// build a slice of "Name (Endpoint)" strings
		pl := make([]string, 0, 8)
		for k, v := range rc.Configurations {
			pl = append(pl, fmt.Sprintf("%s (%s)", k, v.Endpoint))
		}
		fmt.Println(pl)
		sel := promptSelectProfile(pl)
		name := sel[:strings.Index(sel, "(")-1]
		fmt.Fprintf(os.Stdout, "Profile %q selected.\n", name)

		remove := confirm(fmt.Sprintf("Are you sure you want to remove %q", name))
		if remove {
			p := rc.Configurations[name]

			//client := eclient.New(p.Endpoint, timeout)
			hostname, err := configmgr.URLToHostName(p.Endpoint)
			//g, _ := client.GetConfig()

			fmt.Printf("%+v\n", p)
			filename := fmt.Sprintf("%s-%s", hostname, p.DevKey[:6])

			ok, err := configmgr.DeleteProject(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warn: remove profile failed: %+v\n", errors.Cause(err))
			}
			if !ok {
				fmt.Fprintf(os.Stderr, "Warn: remove profile failed: %+v\n", err)
			}

			// delete the configuration and write the new config to the filesystem.
			delete(rc.Configurations, name)
			err = configmgr.WriteConfig(rc)
			if err != nil {
				fmt.Fprintf(os.Stderr, "write config failed: %+v", err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stdout, "Project %q removed.\n", name)
			os.Exit(0)
		}
		fmt.Fprintf(os.Stdout, "Skipping removal\n")
		os.Exit(0)
	},
}

func confirm(msg string) bool {
	answer := false
	prompt := &survey.Confirm{
		Message: msg,
	}
	survey.AskOne(prompt, &answer, nil)
	return answer
}

func init() {
	profilesCmd.AddCommand(profilesRemoveCmd)
}
