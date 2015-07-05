package main

import (
  "os"
  "github.com/codegangsta/cli"
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
