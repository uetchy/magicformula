package main


import (
  "os"
  "github.com/codegangsta/cli"
)

func main() {
  app := cli.NewApp()
  app.Name = "solver"
  app.Commands = []cli.Command{
    CommandPush,
  }

  app.Run(os.Args)
}
