package main

import "github.com/ecommerce-builder/ecom-cli-tool/cmd"

var version string

func main() {
	cmd.Version = version
	cmd.Execute()
}
