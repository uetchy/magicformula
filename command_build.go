package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/libgit2/git2go"
	"github.com/motemen/go-gitconfig"
	"go/build"
	"net/url"
	"os"
	"path/filepath"
)

var buildFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "description",
		Usage:  "Package description",
		EnvVar: "DESCRIPTION",
	},
	cli.StringFlag{
		Name:   "name",
		Usage:  "Package name",
		EnvVar: "PACKAGE_NAME",
	},
	cli.StringFlag{
		Name:   "tag",
		Usage:  "Release tag",
		EnvVar: "RELEASE_TAG",
	},
	cli.StringFlag{
		Name:   "homepage",
		Usage:  "Homepage",
		EnvVar: "HOMEPAGE",
	},
	cli.StringFlag{
		Name:   "path",
		Usage:  "Package path",
		EnvVar: "PACKAGE_PATH",
	},
	cli.StringFlag{
		Name:   "url",
		Usage:  "Package url",
		EnvVar: "PACKAGE_URL",
	},
}

var CommandBuild = cli.Command{
	Name:   "build",
	Flags:  buildFlags,
	Action: doBuild,
}

func doBuild(c *cli.Context) {
	name := c.String("name")
	tag := c.String("tag")
	packagePath := c.String("path")
	packageURL := c.String("url")
	head := c.String("head")
	homepage := c.String("homepage")
	description := c.String("description")
	kind := "golang"

	if name == "" {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		name = filepath.Base(dir)
	}
	if head == "" {
		head, _ = gitconfig.GetString("remote.origin.url")
	}
	if packageURL == "" {
		fmt.Println(errors.New("Missing 'url'"))
		os.Exit(1)
	}
	if tag == "" {
		fmt.Println(errors.New("Missing 'tag'"))
		os.Exit(1)
	}

	formula := Formula{
		Kind:        kind,
		Name:        name,
		Description: description,
		Tag:         tag,
		URL:         packageURL,
		Head:        head,
		Homepage:    homepage,
		PackagePath: packagePath,
	}

	cwd, _ := os.Getwd()
	rootPackage, _ := build.Import(".", cwd, 0)
	deps := make([]Dep, 0)
	for _, dep := range rootPackage.Imports {
		if d, _ := build.Import(dep, cwd, 0); d.Goroot == false {
			depName := d.Name
			depPath := d.ImportPath
			depURL, _ := url.Parse(depPath)
			var properDepURL string
			switch depURL.Host {
			case "golang.org":
				properDepURL = "https://github.com/golang/" + depName + ".git"
			default:
				properDepURL = "https://" + depPath + ".git"
			}
			repo, _ := git.InitRepository(d.Dir, false)
			headCommit, _ := repo.RevparseSingle("HEAD")
			headRevision := headCommit.Id().String()
			deps = append(deps, Dep{
				Name:     depPath,
				URL:      properDepURL,
				Revision: headRevision,
			})
		}
	}
	formula.Deps = deps
	formulaData := formula.Format()

	if c.Bool("push") == false {
		fmt.Println(string(formulaData))
		os.Exit(0)
	}
}

func hashCommit(header string, content []byte) string {
	hasher := sha1.New()
	hasher.Write([]byte(header))
	hasher.Write([]byte("\x00"))
	hasher.Write(content)
	sum := hex.EncodeToString(hasher.Sum(nil))
	return sum
}
