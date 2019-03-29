package main

import "bitbucket.org/andyfusniakteam/ecom-cli-tool/cmd"

var version string

func main() {
	cmd.Version = version
	cmd.Execute()
}
