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
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var CommandPush = cli.Command{
	Name: "push",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "token",
			Usage:  "Github token",
			EnvVar: "GITHUB_TOKEN",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "repo name",
		},
		cli.StringFlag{
			Name:  "owner",
			Usage: "repo owner",
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
			Usage: "Version",
		},
		cli.StringFlag{
			Name:   "tag",
			Usage:  "release tag",
			EnvVar: "RELEASE_TAG",
		},
		cli.StringFlag{
			Name:  "target-path-64",
			Usage: "64bit",
		},
		cli.StringFlag{
			Name:  "target-path-32",
			Usage: "32bit",
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

func hashCommit(header string, content []byte) string {
	hasher := sha1.New()
	hasher.Write([]byte(header))
	hasher.Write([]byte("\x00"))
	hasher.Write(content)
	sum := hex.EncodeToString(hasher.Sum(nil))
	return sum
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
	githubToken := c.String("token")
	if githubToken == "" {
		fmt.Println("Missing 'token'")
		os.Exit(1)
	}
	productName := c.String("name")
	if productName == "" {
		fmt.Println("Missing 'name'")
		os.Exit(1)
	}
	productOwner := c.String("owner")
	if productOwner == "" {
		fmt.Println("Missing 'owner'")
		os.Exit(1)
	}
	productVersion := c.String("version")
	if productVersion == "" {
		fmt.Println("Missing 'version'")
		os.Exit(1)
	}
	targetPath64 := c.String("target-path-64")
	if targetPath64 == "" {
		fmt.Println("Missing 'target-path-64'")
		os.Exit(1)
	}
	targetPath32 := c.String("target-path-32")
	commitMessage := c.String("message")
	if commitMessage == "" {
		commitMessage = "Deploy from " + productOwner
	}
	commitAuthor := c.String("committer")
	if commitAuthor == "" {
		fmt.Println("Missing 'committer'")
		os.Exit(1)
	}
	commitAuthorEmail := c.String("committer-email")
	if commitAuthorEmail == "" {
		fmt.Println("Missing 'committer-email'")
		os.Exit(1)
	}
	releaseTag := c.String("tag")
	if releaseTag == "" {
		releaseTag = "v" + productVersion
	}

	// Generate formula
	githubReleaseUrl := "https://github.com/" + productOwner + "/" + productName + "/releases/download/" + releaseTag
	inv := map[string]string{
		"BinName":      productName,
		"ClassName":    toCamelCase(productName),
		"Version":      productVersion,
		"TargetUrl64":  githubReleaseUrl + filepath.Join("/", filepath.Base(targetPath64)),
		"TargetHash64": hashFileSum(targetPath64),
	}
	if targetPath32 != "" {
		inv["TargetUrl32"] = githubReleaseUrl + filepath.Join("/", filepath.Base(targetPath32))
		inv["TargetHash32"] = hashFileSum(targetPath32)
	}
	tmpl, _ := Asset("templates/formula.tmpl")
	t := template.New("formula")
	template.Must(t.Parse(string(tmpl)))
	var buf bytes.Buffer
	err := t.Execute(&buf, inv)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	formula := buf.Bytes()

	content := &github.RepositoryContentFileOptions{
		Message: &commitMessage,
		Content: formula,
		Committer: &github.CommitAuthor{
			Name:  &commitAuthor,
			Email: &commitAuthorEmail,
		},
	}

	// Prepare for github API request
	formulaRepo := "homebrew-" + productName
	formulaFile := productName + ".rb"
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	// Fetch previous file's SHA hash
	stat, _, _, _ := client.Repositories.GetContents(
		productOwner,
		formulaRepo,
		formulaFile,
		&github.RepositoryContentGetOptions{},
	)

	if stat != nil {
		// Avoid no-change commit
		header := "blob " + fmt.Sprintf("%v", len(formula))
		sha := hashCommit(header, formula)
		if *stat.SHA == sha {
			fmt.Println("No changes")
			os.Exit(0)
		}
		content.SHA = stat.SHA

		// Upload changes
		res, _, err := client.Repositories.UpdateFile(
			productOwner,
			formulaRepo,
			formulaFile,
			content,
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(*res.SHA)
	} else {
		// Create file
		res, _, err := client.Repositories.CreateFile(
			productOwner,
			formulaRepo,
			formulaFile,
			content,
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(*res.SHA)
	}
}
