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
		Name:  "kind",
		Usage: "Kind of packages",
		Value: "general",
	},
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
	cli.StringSliceFlag{
		Name:   "build-dep",
		Usage:  "Build dependencies",
		EnvVar: "BUILD_DEP",
	},
	cli.StringFlag{
		Name:  "zsh-completions",
		Usage: "Path to zsh completions",
	},
}

var CommandBuild = cli.Command{
	Name:   "build",
	Flags:  buildFlags,
	Action: doBuild,
}

func doBuild(c *cli.Context) {
	kind := c.String("kind")
	name := c.String("name")
	tag := c.String("tag")
	packagePath := c.String("path")
	packageURL := c.String("url")
	head := c.String("head")
	homepage := c.String("homepage")
	description := c.String("description")
	zshCompletionsPath := c.String("zsh-completions")
	buildDeps := make([]Dep, 0)
	for _, dep := range c.StringSlice("build-dep") {
		buildDeps = append(buildDeps, Dep{dep, "build"})
	}

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
		Kind:               kind,
		Name:               name,
		Description:        description,
		Tag:                tag,
		URL:                packageURL,
		Head:               head,
		Homepage:           homepage,
		PackagePath:        packagePath,
		Deps:               buildDeps,
		ZshCompletionsPath: zshCompletionsPath,
	}

	switch formula.Kind {
	case "golang":
		cwd, _ := os.Getwd()
		rootPackage, _ := build.Import(".", cwd, 0)
		goResources := make([]GoResource, 0)
		for _, goResource := range rootPackage.Imports {
			if d, _ := build.Import(goResource, cwd, 0); d.Goroot == false {
				goResourceName := d.Name
				goResourcePath := d.ImportPath
				goResourceURL, _ := url.Parse(goResourcePath)
				var properDepURL string
				switch goResourceURL.Host {
				case "golang.org":
					properDepURL = "https://github.com/golang/" + goResourceName + ".git"
				default:
					properDepURL = "https://" + goResourcePath + ".git"
				}
				repo, _ := git.InitRepository(d.Dir, false)
				headCommit, _ := repo.RevparseSingle("HEAD")
				headRevision := headCommit.Id().String()
				goResources = append(goResources, GoResource{
					Name:     goResourcePath,
					URL:      properDepURL,
					Revision: headRevision,
				})
			}
		}
		formula.GoResources = goResources
	}

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
