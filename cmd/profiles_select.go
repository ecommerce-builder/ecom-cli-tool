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

// projectsListCmd represents the projectsList command
var projectsSelectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select and change to a new configuration.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// build a slice of "Name (Endpoint)" strings
		pl := make([]string, 0, 8)
		for k, v := range rc.Configurations {
			pl = append(pl, fmt.Sprintf("%s (%s)", k, v.Endpoint))
		}

		sel := promptSelectProject(pl)
		name := sel[:strings.Index(sel, "(")-1]
		fmt.Fprintf(os.Stdout, "Project %q selected.\n", name)

		err := configmgr.WriteCurrentProject(name)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	projectsCmd.AddCommand(projectsSelectCmd)
}

func promptSelectProject(pl []string) string {
	proj := ""
	prompt := &survey.Select{
		Message: "Select a project:",
		Options: pl,
	}
	survey.AskOne(prompt, &proj, nil)
	return proj
}
