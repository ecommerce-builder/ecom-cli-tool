package catalog

import (
	"github.com/spf13/cobra"
)

// NewCmdCatalog returns new initialized instance of the catalog sub command
func NewCmdCatalog() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "catalog",
		Short: "Catalog management",
	}
	cmd.AddCommand(NewCmdCatalogApply())
	cmd.AddCommand(NewCmdCatalogGet())
	cmd.AddCommand(NewCmdCatalogPurge())
	return cmd
}
