package profiles

import (
	"github.com/spf13/cobra"
)

// NewCmdProfiles returns new initialized instance of profiles sub command
func NewCmdProfiles() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "profiles",
		Short: "Profile management",
	}
	cmd.AddCommand(NewCmdProfilesCreate())
	cmd.AddCommand(NewCmdProfilesList())
	cmd.AddCommand(NewCmdProfilesRemove())
	cmd.AddCommand(NewCmdProfilesSelect())
	return cmd
}
