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
	Short: "Select and change to a project.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// build a slice of "Name (Endpoint)" strings
		pl := make([]string, 0, 8)
		for k, v := range rc.Configurations {
			fmt.Println(v)
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
