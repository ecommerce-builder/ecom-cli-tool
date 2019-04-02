package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// signoutCmd represents the signout command
var signoutCmd = &cobra.Command{
	Use:   "signout",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// webKey, endpoint, err := promptSignOut()
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "%v", err)
		// 	os.Exit(1)
		// }

		// exists, err := configmgr.Exists(filename)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "%v\n", err)
		// 	os.Exit(1)
		// }

		// if !exists {
		// 	fmt.Fprintf(os.Stderr, "No configuration exists for the given Web API Key and Endpoint.")
		// 	os.Exit(1)
		// }
	},
}

func promptSignOut() (webKey, endpoint string, err error) {
	k := &survey.Input{
		Message: "Web API Key:",
	}
	survey.AskOne(k, &webKey, nil)

	e := &survey.Input{
		Message: "Endpoint:",
	}
	survey.AskOne(e, &endpoint, nil)

	return webKey, endpoint, nil
}

func init() {
	rootCmd.AddCommand(signoutCmd)
}
