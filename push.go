package main

import (
  "os"
  "fmt"
  "regexp"
  "strings"
  "bytes"
  "io"
  "text/template"
  "crypto/sha1"
  "encoding/hex"
  "golang.org/x/oauth2"
  "github.com/google/go-github/github"
  "github.com/codegangsta/cli"
)

var CommandPush = cli.Command{
  Name: "push",
  Flags: []cli.Flag{
    cli.StringFlag{
      Name: "token",
      Usage: "Github token",
    },
    cli.StringFlag{
      Name: "name",
      Usage: "Github name",
    },
    cli.StringFlag{
      Name: "owner",
      Usage: "Github owner",
    },
    cli.StringFlag{
      Name: "version",
      Usage: "Github version",
    },
    cli.StringFlag{
      Name: "package-path-64",
      Usage: "package-path-64",
    },
    cli.StringFlag{
      Name: "package-path-32",
      Usage: "package-path-32",
    },
    cli.StringFlag{
      Name: "tag",
      Usage: "release tag",
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

func hashSum(path string) string {
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

  releaseTag := c.String("tag")
  if releaseTag == "" {
    releaseTag = "v" + packageVersion
  }
  packageUrl64 := "https://github.com/" + packageOwner + "/" + packageName + "/releases/download/" + releaseTag + "/" + binaryPath64
  packageHash64 := hashSum(binaryPath64)
  packageUrl32 := "https://github.com/" + packageOwner + "/" + packageName + "/releases/download/" + releaseTag + "/" + binaryPath32
  packageHash32 := hashSum(binaryPath32)

  inv := map[string]string {
    "BinName": packageName,
    "Version": packageVersion,
    "PackageUrl64": packageUrl64,
    "PackageHash64": packageHash64,
    "PackageUrl32": packageUrl32,
    "PackageHash32": packageHash32,
    "ClassName": toCamelCase(packageName),
  }

  // Load formula template
  tmpl, _ := Asset("formula.tmpl")
  t := template.New("formula")
  template.Must(t.Parse(string(tmpl)))
  var buf bytes.Buffer
  t.Execute(&buf, inv)
  formula := buf.String()
  fmt.Println(formula)

  // TODO: Upload to Github
  ts := oauth2.StaticTokenSource(
    &oauth2.Token{AccessToken: githubToken},
  )
  tc := oauth2.NewClient(oauth2.NoContext, ts)
  client := github.NewClient(tc)

  // release, _, _ := client.Repositories.GetRelease(packageOwner, packageName, 1473977)
  // fmt.Println(*release.AssetsURL)
}
