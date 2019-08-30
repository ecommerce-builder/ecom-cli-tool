package categoriestree

import (
	"github.com/spf13/cobra"
)

// NewCmdCategoriesTree returns new initialized instance of the catalog sub command
func NewCmdCategoriesTree() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "categories-tree",
		Short: "Categories Tree management",
	}
	cmd.AddCommand(NewCmdCategoriesTreeApply())
	cmd.AddCommand(NewCmdCategoriesTreeGet())
	cmd.AddCommand(NewCmdCategoriesTreeDelete())
	return cmd
}
