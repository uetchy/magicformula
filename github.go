package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func hashCommit(header string, content []byte) string {
	hasher := sha1.New()
	hasher.Write([]byte(header))
	hasher.Write([]byte("\x00"))
	hasher.Write(content)
	sum := hex.EncodeToString(hasher.Sum(nil))
	return sum
}

func PushToGithub(name string, owner string, formulaData []byte, token string, committer string, committerEmail string, message string) error {
	content := &github.RepositoryContentFileOptions{
		Message: &message,
		Content: formulaData,
		Committer: &github.CommitAuthor{
			Name:  &committer,
			Email: &committerEmail,
		},
	}

	// Prepare for github API request
	repo := "homebrew-" + name
	filename := name + ".rb"
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	// Fetch previous file's SHA hash
	stat, _, _, _ := client.Repositories.GetContents(
		owner,
		repo,
		filename,
		&github.RepositoryContentGetOptions{},
	)

	if stat != nil {
		// Avoid no-change commit
		header := "blob " + fmt.Sprintf("%v", len(formulaData))
		sha := hashCommit(header, formulaData)
		if *stat.SHA == sha {
			return errors.New("No changes")
		}
		content.SHA = stat.SHA

		// Upload changes
		res, _, err := client.Repositories.UpdateFile(
			owner,
			repo,
			filename,
			content,
		)
		if err != nil {
			return err
		}

		fmt.Println(*res.SHA)
	} else {
		// Create file
		res, _, err := client.Repositories.CreateFile(
			owner,
			repo,
			filename,
			content,
		)
		if err != nil {
			return err
		}

		fmt.Println(*res.SHA)
	}

	return nil
}
