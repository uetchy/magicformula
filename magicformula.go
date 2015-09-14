package main

import (
	"errors"
	"github.com/codegangsta/cli"
	"os"
)

var Version string = "HEAD"
var Commands = []cli.Command{
	CommandBuild,
	CommandPush,
}
var globalFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "github-token",
		Usage:  "Github token",
		EnvVar: "GITHUB_TOKEN",
	},
	cli.StringFlag{
		Name:   "github-user",
		Usage:  "Repo owner",
		EnvVar: "GITHUB_USER",
	},
	cli.StringFlag{
		Name:   "package-name",
		Usage:  "Package name",
		EnvVar: "MF_PACKAGE_NAME",
	},
	cli.StringFlag{
		Name:   "commit-message",
		Usage:  "Commit message",
		EnvVar: "MF_COMMIT_MESSAGE",
	},
	cli.StringFlag{
		Name:   "committer",
		Usage:  "Committer name",
		EnvVar: "MF_COMMITTER",
	},
	cli.StringFlag{
		Name:   "committer-email",
		Usage:  "Committer email",
		EnvVar: "MF_COMMITTER_EMAIL",
	},
	cli.StringFlag{
		Name:   "release-tag",
		Usage:  "Release tag",
		EnvVar: "RELEASE_TAG",
	},
	cli.StringFlag{
		Name:   "package-path",
		Usage:  "Package path for 64bit arch",
		EnvVar: "MF_PACKAGE_PATH",
	},
	cli.StringFlag{
		Name:   "homepage",
		Usage:  "Homepage",
		EnvVar: "MF_HOMEPAGE",
	},
}

func assertContext(c *cli.Context) error {
	if c.String("token") == "" {
		return errors.New("Missing 'token'")
	}

	if c.String("name") == "" {
		return errors.New("Missing 'name'")
	}

	if c.String("owner") == "" {
		return errors.New("Missing 'owner'")
	}

	if c.String("version") == "" {
		return errors.New("Missing 'version'")
	}

	if c.String("target-path") == "" {
		return errors.New("Missing 'target-path'")
	}

	if c.String("committer") == "" {
		return errors.New("Missing 'committer'")
	}

	if c.String("committer-email") == "" {
		return errors.New("Missing 'committer-email'")
	}

	return nil
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
