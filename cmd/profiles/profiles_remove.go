package profiles

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// NewCmdProfilesRemove returns new initialized instance of remove sub command
func NewCmdProfilesRemove() *cobra.Command {
	cfgs, _, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove a profile",
		Long:  `Removes a profile dropping the credentials`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(cfgs.Configurations) == 0 {
				fmt.Println("No profiles")
				os.Exit(0)
			}
			// build a slice of "Name (Endpoint)" strings
			pl := make([]string, 0, 8)
			for k, v := range cfgs.Configurations {
				pl = append(pl, fmt.Sprintf("%s (%s)", k, v.Endpoint))
			}

			sel := promptSelectProfile(pl)
			name := sel[:strings.Index(sel, "(")-1]
			fmt.Fprintf(os.Stdout, "Profile %q selected.\n", name)

			remove := confirm(fmt.Sprintf("Are you sure you want to remove %q", name))
			if remove {
				p := cfgs.Configurations[name]
				hostname, err := configmgr.URLToHostName(p.Endpoint)
				filename := fmt.Sprintf("%s-%s", hostname, p.DevKey[:6])
				ok, err := configmgr.DeleteProject(filename)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warn: remove profile failed: %+v\n", errors.Unwrap(err))
				}
				if !ok {
					fmt.Fprintf(os.Stderr, "Warn: remove profile failed: %+v\n", err)
				}
				// delete the configuration and write the new config to the filesystem.
				delete(cfgs.Configurations, name)
				if err = configmgr.WriteConfig(cfgs); err != nil {
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
	return cmd
}

func confirm(msg string) bool {
	prompt := &survey.Confirm{
		Message: msg,
	}
	var answer bool
	survey.AskOne(prompt, &answer, nil)
	return answer
}
