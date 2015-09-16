package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

type Formula struct {
	Kind        string // "golang"
	Name        string // "magicformula"
	Description string // "awesome app"
	Tag         string // "v1.2.0"
	Revision    string // "b6f7aadbeb21ae18972577173ce175af83ce239d"
	URL         string // "https://github.com/uetchy/magicformula-1.8.tar.gz"
	Head        string // "https://github.com/uetchy/magicformula.git"
	Homepage    string // "https://github.com/uetchy/magicformula"
	PackagePath  string // "/path/to/bin/magicformula"
	Deps        []Dep
}

type Dep struct {
	Name     string // "cli"
	URL      string // "github.com/codegangsta/cli"
	Revision string // "b6f7aadbeb21ae18972577173ce175af83ce239d"
}

func (f *Formula) CheckSum() string {
	hasher := sha1.New()
	fp, err := os.Open(f.PackagePath)
	if err != nil {
		fmt.Println(err)
	}
	defer fp.Close()
	if _, err := io.Copy(hasher, fp); err != nil {
		fmt.Println(err)
	}
	sum := hex.EncodeToString(hasher.Sum(nil))
	return sum
}

func (f *Formula) ClassName() string {
	re, _ := regexp.Compile("[_-](.)")
	res := re.ReplaceAllStringFunc(f.Name, func(j string) string {
		return strings.Title(j[1:])
	})
	return strings.Title(res)
}

func (f *Formula) Dir() string {
	return filepath.Dir(f.PackagePath)
}

func (f *Formula) Format(tmplName string) []byte {
	tmpl := template.Must(template.ParseFiles(tmplName))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return buf.Bytes()
}
