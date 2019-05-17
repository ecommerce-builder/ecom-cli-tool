package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"bitbucket.org/andyfusniakteam/ecom-cli-tool/configmgr"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// profilessListCmd represents the profilesList command
var profilesSelectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select and change to a new profile",
	Run: func(cmd *cobra.Command, args []string) {
		// build a slice of "Name (Endpoint)" strings
		pl := make([]string, 0, 8)
		for k, v := range rc.Configurations {
			pl = append(pl, fmt.Sprintf("%s (%s)", k, v.Endpoint))
		}

		sel := promptSelectProfile(pl)
		name := sel[:strings.Index(sel, "(")-1]
		fmt.Fprintf(os.Stdout, "Profile %q selected.\n", name)

		err := configmgr.WriteCurrentProject(name)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	profilesCmd.AddCommand(profilesSelectCmd)
}

func promptSelectProfile(pl []string) string {
	profile := ""
	prompt := &survey.Select{
		Message: "Select a profile:",
		Options: pl,
	}
	survey.AskOne(prompt, &profile, nil)
	return profile
}
