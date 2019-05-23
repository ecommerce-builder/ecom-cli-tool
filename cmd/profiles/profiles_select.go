package profiles

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// NewCmdProfilesSelect returns new initialized instance of select sub command
func NewCmdProfilesSelect() *cobra.Command {
	cfgs, _, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "select",
		Short: "Select and change to a new profile",
		Run: func(cmd *cobra.Command, args []string) {
			// build a slice of "Name (Endpoint)" strings
			pl := make([]string, 0, 8)
			for k, v := range cfgs.Configurations {
				pl = append(pl, fmt.Sprintf("%s (%s %s %s)", k, v.Endpoint, v.Customer.Email, v.Customer.Role))
			}
			sel := promptSelectProfile(pl)
			name := sel[:strings.Index(sel, "(")-1]
			fmt.Fprintf(os.Stdout, "Profile %q selected.\n", name)
			if err := configmgr.WriteCurrentProject(name); err != nil {
				log.Fatal(err)
			}
		},
	}
	return cmd
}

func promptSelectProfile(pl []string) string {
	prompt := &survey.Select{
		Message: "Select a profile:",
		Options: pl,
	}
	var profile string
	survey.AskOne(prompt, &profile, nil)
	return profile
}
