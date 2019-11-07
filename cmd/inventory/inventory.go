package inventory

import (
	"github.com/spf13/cobra"
)

// NewCmdInventory returns new initialized instance of inventory sub command
func NewCmdInventory() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "inventory",
		Short: "Inventory management",
	}
	cmd.AddCommand(NewCmdInventoryGet())
	cmd.AddCommand(NewCmdInventoryList())
	cmd.AddCommand(NewCmdInventoryUpdate())
	cmd.AddCommand(NewCmdInventoryBatchUpdate())
	return cmd
}
