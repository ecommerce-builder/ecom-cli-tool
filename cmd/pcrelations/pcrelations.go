package pcrelations

import (
	"github.com/spf13/cobra"
)

// NewCmdPCRelations returns new initialized instance of assocs sub command
func NewCmdPCRelations() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "product-category-relations",
		Short: "product to category relations",
	}
	cmd.AddCommand(NewCmdPCRelationsApply())
	cmd.AddCommand(NewCmdPCRelationsList())
	cmd.AddCommand(NewCmdPCRelationsDelete())
	return cmd
}
