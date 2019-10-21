package devkeys

import (
	"github.com/spf13/cobra"
)

// NewCmdDevKeys returns new initialized instance of devkeys sub command
func NewCmdDevKeys() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "devkeys",
		Short: "Developer Keys Management",
	}
	cmd.AddCommand(NewCmdDevKeysCreate())
	cmd.AddCommand(NewCmdDevKeysList())
	return cmd
}
