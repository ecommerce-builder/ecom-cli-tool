package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// profilesListCmd represents the profilesList command
var profilesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Display a list of available profiles",
	Run: func(cmd *cobra.Command, args []string) {
		format := "%v\t%v\t%v\n"
		tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprintf(tw, format, "Active", "Endpoint", "Dev Key")
		fmt.Fprintf(tw, format, "------", "--------", "-------")
		for k, v := range rc.Configurations {
			var active string
			if currentConfigName == k {
				active = "  *"
			} else {
				active = ""
			}
			fmt.Fprintf(tw, format, active, v.Endpoint, v.DevKey[0:5]+"********")
		}
		tw.Flush()
	},
}

func init() {
	profilesCmd.AddCommand(profilesListCmd)
}
