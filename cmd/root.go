package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/spf13/cobra"
)

var timeout = time.Duration(10 * time.Second)

var rc *configmgr.EcomConfigurations
var currentConfigName string

var rootCmd = &cobra.Command{
	Use:   "ecom",
	Short: "ecom is a CLI tool for administering ecommerce systems",
	Long:  `See the user guide for more details.`,
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	var err error
	rc, err = configmgr.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	currentConfigName, err = configmgr.ReadCurrentConfigName()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}
}

// Execute the command line tool
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
