package main

import (
	"github.com/codegangsta/cli"
	"os"
)

var Version string = "HEAD"
var Commands = []cli.Command{
	CommandBuild,
}

func main() {
	app := cli.NewApp()
	app.Name = "magicformula"
	app.Usage = "Generate and upload Homebrew Formula"
	app.Version = Version
	app.Author = "Yasuaki Uechi"
	app.Email = "uetchy@randompaper.co"
	app.Commands = Commands
	app.Run(os.Args)
}
