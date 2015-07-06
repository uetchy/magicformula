package main

import (
	"github.com/codegangsta/cli"
	"os"
)

var Commands = []cli.Command{
	CommandPush,
}

func main() {
	app := cli.NewApp()
	app.Name = "solver"
	app.Commands = Commands
	app.Run(os.Args)
}
