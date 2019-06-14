package token

import (
	"fmt"
	"log"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// NewCmdTokenShow returns new initialized instance of token sub command
func NewCmdTokenShow() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "show",
		Short: "Show the current JSON Web Token",
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}

			file, err := configmgr.TokenFilename(&current)
			if err != nil {
				fmt.Fprintf(os.Stderr, "token file %q not found", file)
				os.Exit(1)
			}

			tar, err := configmgr.ReadTokenAndRefreshToken(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "tokenand refresh token cannot be read from %q: %v", file, err)
				os.Exit(1)
			}

			fmt.Printf("%s\n", tar.IDToken)
			os.Exit(0)
		},
	}
	return cmd
}
