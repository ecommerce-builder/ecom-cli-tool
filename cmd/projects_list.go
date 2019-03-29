package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// projectsListCmd represents the projectsList command
var projectsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Display a list of projects available to this command line tool.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println(rc.Configurations)
		fmt.Println(currentConfigName)

		format := "%v\t%v\t%v\t%v\t%v\t\n"
		tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprintf(tw, format, "Name", "Active", "Endpoint", "Firebase API Key", "Dev Key")
		fmt.Fprintf(tw, format, "----", "------", "--------", "----------------", "-------")
		for k, v := range rc.Configurations {
			var active string
			if currentConfigName == k {
				active = "*"
			} else {
				active = ""
			}
			fmt.Fprintf(tw, format, k, active, v.Endpoint, v.FirebaseAPIKey, v.DevKey[0:5]+"********")
		}
		tw.Flush()
	},
}

func init() {
	projectsCmd.AddCommand(projectsListCmd)
}
