package cmd

import (
	"fmt"
	"os"
	"strings"

	"bitbucket.org/andyfusniakteam/ecom-cli-tool/configmgr"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// projectsListCmd represents the projectsList command
var projectsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a project.",
	Long:  `Removes a project and drops all credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		// build a slice of "Name (Endpoint)" strings
		pl := make([]string, 0, 8)
		for k, v := range rc.Configurations {
			pl = append(pl, fmt.Sprintf("%s (%s)", k, v.Endpoint))
		}

		sel := promptSelectProject(pl)
		name := sel[:strings.Index(sel, "(")-1]
		fmt.Fprintf(os.Stdout, "Project %q selected.\n", name)

		remove := confirm(fmt.Sprintf("Are you sure you want to remove %q", name))
		if remove {
			p := rc.Configurations[name]

			hostname, err := configmgr.URLToHostName(p.Endpoint)
			filename := fmt.Sprintf("%s-%s", p.FirebaseAPIKey, hostname)

			ok, err := configmgr.DeleteProject(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warn: remove project failed: %+v\n", errors.Cause(err))
			}
			if !ok {
				fmt.Fprintf(os.Stderr, "Warn: remove project failed: %+v\n", err)
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
	projectsCmd.AddCommand(projectsRemoveCmd)
}
