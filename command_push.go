package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io"
	"os"
	"regexp"
	"strings"
	"text/template"
)

var CommandPush = cli.Command{
	Name: "push",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "github-token",
			Usage: "Github token",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "Github name",
		},
		cli.StringFlag{
			Name:  "owner",
			Usage: "Github owner",
		},
		cli.StringFlag{
			Name:  "message",
			Usage: "commit message",
		},
		cli.StringFlag{
			Name:  "committer",
			Usage: "committer name",
		},
		cli.StringFlag{
			Name:  "committer-email",
			Usage: "committer email",
		},
		cli.StringFlag{
			Name:  "version",
			Usage: "version",
		},
		cli.StringFlag{
			Name:  "tag",
			Usage: "release tag",
		},
		cli.StringFlag{
			Name:  "package-path-64",
			Usage: "package-path-64",
		},
		cli.StringFlag{
			Name:  "package-path-32",
			Usage: "package-path-32",
		},
	},
	Action: doPush,
}

func toCamelCase(str string) string {
	re, _ := regexp.Compile("[_-](.)")
	res := re.ReplaceAllStringFunc(str, func(j string) string {
		return strings.Title(j[1:])
	})
	return strings.Title(res)
}

func hashFileSum(path string) string {
	hasher := sha1.New()
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	if _, err := io.Copy(hasher, f); err != nil {
		fmt.Println(err)
	}
	sum := hex.EncodeToString(hasher.Sum(nil))
	return sum
}

func doPush(c *cli.Context) {
	githubToken := c.String("github-token")
	packageName := c.String("name")
	packageOwner := c.String("owner")
	packageVersion := c.String("version")
	binaryPath64 := c.String("package-path-64")
	binaryPath32 := c.String("package-path-32")
	commitMessage := c.String("message")
	commitAuthor := c.String("committer")
	commitAuthorEmail := c.String("committer-email")
	releaseTag := c.String("tag")
	if releaseTag == "" {
		releaseTag = "v" + packageVersion
	}

	// Generate formula
	githubRepoUrl := "https://github.com/" + packageOwner + "/" + packageName
	githubReleaseUrl := githubRepoUrl + "/releases/download/" + releaseTag
	inv := map[string]string{
		"BinName":       packageName,
		"ClassName":     toCamelCase(packageName),
		"Version":       packageVersion,
		"PackageUrl64":  githubReleaseUrl + "/" + binaryPath64,
		"PackageHash64": hashFileSum(binaryPath64),
		"PackageUrl32":  githubReleaseUrl + "/" + binaryPath32,
		"PackageHash32": hashFileSum(binaryPath32),
	}
	tmpl, _ := Asset("formula.tmpl")
	t := template.New("formula")
	template.Must(t.Parse(string(tmpl)))
	var buf bytes.Buffer
	t.Execute(&buf, inv)
	formula := buf.String()

	// Prepare for github API request
	formulaRepo := "homebrew-" + packageName
	formulaFile := packageName + ".rb"
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	// Fetch previous file's SHA hash
	stat, _, _, err := client.Repositories.GetContents(
		packageOwner,
		formulaRepo,
		formulaFile,
		&github.RepositoryContentGetOptions{},
	)
	if err != nil {
		fmt.Println("Error")
		os.Exit(1)
	}

	// Update formula
	content := &github.RepositoryContentFileOptions{
		Message: &commitMessage,
		Content: []byte(formula),
		SHA:     stat.SHA,
		Committer: &github.CommitAuthor{
			Name:  &commitAuthor,
			Email: &commitAuthorEmail,
		},
	}
	res, _, err := client.Repositories.UpdateFile(
		packageOwner,
		formulaRepo,
		formulaFile,
		content,
	)
	fmt.Println(res)
	if err != nil {
		fmt.Println("Error")
		os.Exit(1)
	}
}
