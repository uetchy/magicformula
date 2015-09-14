package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"path/filepath"
)

var CommandBuild = cli.Command{
	Name:   "build",
	Flags:  globalFlags,
	Action: doBuild,
}

func doBuild(c *cli.Context) {
	name := c.String("name")
	if name == "" {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		name = filepath.Base(dir)
	}
	formula := Formula{
		Kind:        "golang",
		Name:        name,
		Description: "test",
		Version:     c.String("version"),
		URL:         c.String("url"),
		Homepage:    c.String("homepage"),
		TargetPath:  c.String("target-path"),
	}
	formulaData := formula.Format("templates/formula_golang.tmpl")

	fmt.Println(string(formulaData))
}
