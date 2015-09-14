package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

var CommandPush = cli.Command{
	Name:   "push",
	Flags:  globalFlags,
	Action: doPush,
}

func doPush(c *cli.Context) {
	err := assertContext(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Generate formula
	formula := Formula{
		Name:        c.String("name"),
		Description: "test",
		Version:     c.String("version"),
		URL:         c.String("url"),
		Homepage:    c.String("homepage"),
		TargetPath:  c.String("target-path"),
	}
	formulaData := formula.Format("templates/formula_go.tmpl")

	// githubReleaseUrl := "https://github.com/" + productOwner + "/" + productName + "/releases/download/" + releaseTag
	// githubReleaseUrl + filepath.Join("/", filepath.Base(targetPath))

	// Push to github
	err = PushToGithub(c.String("name"), c.String("owner"), formulaData, c.String("token"), c.String("committer"), c.String("committer-email"), c.String("commit-message"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
