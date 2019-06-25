package main

import (
	"fmt"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/cmd"
)

var version string

func main() {
	cmd.Version = version
	root := cmd.NewEcomCmd()
	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}
}
