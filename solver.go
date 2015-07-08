package main

import (
	"github.com/codegangsta/cli"
	"os"
)

var Version string = "HEAD"
var Commands = []cli.Command{
	CommandPush,
}

func main() {
	app := cli.NewApp()
	app.Name = "solver"
	app.Usage = "Keep your Homebrew's Formula fresh"
	app.Version = Version
	app.Author = "Yasuaki Uechi"
	app.Email = "uetchy@randompaper.co"
	app.Commands = Commands
	app.Run(os.Args)
}
