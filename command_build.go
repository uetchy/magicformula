package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
	"github.com/motemen/go-gitconfig"
	"golang.org/x/oauth2"
	"os"
	"path/filepath"
)

var buildFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "push",
		Usage: "Push formula to GitHub",
	},
	cli.StringFlag{
		Name:   "description",
		Usage:  "Package description",
		EnvVar: "DESCRIPTION",
	},
	cli.StringFlag{
		Name:   "token",
		Usage:  "Github token",
		EnvVar: "GITHUB_TOKEN",
	},
	cli.StringFlag{
		Name:   "owner",
		Usage:  "Repo owner",
		EnvVar: "GITHUB_USER",
	},
	cli.StringFlag{
		Name:   "name",
		Usage:  "Package name",
		EnvVar: "PACKAGE_NAME",
	},
	cli.StringFlag{
		Name:   "commit-message",
		Usage:  "Commit message",
		EnvVar: "COMMIT_MESSAGE",
	},
	cli.StringFlag{
		Name:   "committer",
		Usage:  "Committer name",
		EnvVar: "COMMITTER",
	},
	cli.StringFlag{
		Name:   "committer-email",
		Usage:  "Committer email",
		EnvVar: "COMMITTER_EMAIL",
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
		Name:   "package-path",
		Usage:  "Package path for 64bit arch",
		EnvVar: "PACKAGE_PATH",
	},
	cli.StringFlag{
		Name:   "homepage",
		Usage:  "Homepage",
		EnvVar: "HOMEPAGE",
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
	packagePath := c.String("package-path")
	url := c.String("url")
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
	if tag == "" {
		fmt.Println(errors.New("Missing 'tag'"))
		os.Exit(1)
	}
	if packagePath == "" {
		fmt.Println(errors.New("Missing 'package-path'"))
		os.Exit(1)
	}

	formula := Formula{
		Kind:        kind,
		Name:        name,
		Description: description,
		Tag:         tag,
		URL:         url,
		Head:        head,
		Homepage:    homepage,
		PackagePath:  packagePath,
	}
	formulaData := formula.Format("templates/formula_golang.tmpl")

	if c.Bool("push") == false {
		fmt.Println(string(formulaData))
		os.Exit(0)
	}

	// Push to GitHub
	token := c.String("token")
	owner := c.String("owner")
	committer := c.String("committer")
	committerEmail := c.String("committer-email")
	commitMessage := c.String("commit-message")

	if token == "" {
		fmt.Println(errors.New("Missing 'token'"))
		os.Exit(1)
	}
	if owner == "" {
		githubUser, err := gitconfig.GetString("github.user")
		if err != nil {
			owner = githubUser
		} else {
			fmt.Println(errors.New("Missing 'owner'"))
			os.Exit(1)
		}
	}
	if committer == "" {
		userName, err := gitconfig.GetString("user.name")
		if err != nil {
			committer = userName
		} else {
			fmt.Println(errors.New("Missing 'committer'"))
			os.Exit(1)
		}
	}
	if committerEmail == "" {
		userEmail, err := gitconfig.GetString("user.email")
		if err != nil {
			committerEmail = userEmail
		} else {
			fmt.Println(errors.New("Missing 'committer-email'"))
			os.Exit(1)
		}
	}

	content := &github.RepositoryContentFileOptions{
		Message: &commitMessage,
		Content: formulaData,
		Committer: &github.CommitAuthor{
			Name:  &committer,
			Email: &committerEmail,
		},
	}

	// Prepare for github API request
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	repo := "homebrew-" + name
	filename := name + ".rb"

	// Fetch previous file's SHA hash
	stat, _, _, _ := client.Repositories.GetContents(
		owner,
		repo,
		filename,
		&github.RepositoryContentGetOptions{},
	)

	var fileFunc func(string, string, string, *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error)
	if stat != nil {
		// Avoid no-change commit
		header := "blob " + fmt.Sprintf("%v", len(formulaData))
		sha := hashCommit(header, formulaData)
		if *stat.SHA == sha {
			fmt.Println(errors.New("No changes"))
			os.Exit(0)
		}
		content.SHA = stat.SHA

		// Upload changes
		fileFunc = client.Repositories.UpdateFile
	} else {
		fileFunc = client.Repositories.CreateFile
	}
	res, _, err := fileFunc(
		owner,
		repo,
		filename,
		content,
	)
	fmt.Println(*res.SHA) // DEBUG:

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
