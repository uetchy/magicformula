package main

import (
	"bytes"
	"crypto/sha256"
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
	Kind               string // "golang"
	Name               string // "magicformula"
	Description        string // "awesome app"
	Tag                string // "v1.2.0"
	Revision           string // "b6f7aadbeb21ae18972577173ce175af83ce239d"
	URL                string // "https://~~/magicformula.tar.gz" or "https://~~/magicformula.git"
	Head               string // "https://github.com/uetchy/magicformula.git"
	Homepage           string // "https://github.com/uetchy/magicformula"
	PackagePath        string // "/path/to/bin/magicformula"
	ZshCompletionsPath string
	Deps               []Dep
	GoResources        []GoResource
}

type Dep struct {
	Name  string
	Extra string
}

// Returns kind of package url
func (f *Formula) URLScheme() string {
	switch filepath.Ext(f.URL) {
	case ".gz":
		return "archive"
	case ".git":
		return "scm"
	default:
		return "binary"
	}
}

func (f *Formula) TemplatePath(phase string) string {
	return f.Kind + "_" + phase
}

// Generate sha256 checksum from PackagePath.
func (f *Formula) CheckSum() string {
	hasher := sha256.New()
	fp, err := os.Open(f.PackagePath)
	if err != nil {
		return ("")
	}
	defer fp.Close()
	if _, err := io.Copy(hasher, fp); err != nil {
		return ("")
	}
	sum := hex.EncodeToString(hasher.Sum(nil))
	return sum
}

// Returns capitalized camel case of Name.
func (f *Formula) ClassName() string {
	re, _ := regexp.Compile("[_-](.)")
	res := re.ReplaceAllStringFunc(f.Name, func(j string) string {
		return strings.Title(j[1:])
	})
	return strings.Title(res)
}

// Generate formula file
func (f *Formula) Format() []byte {
	t := template.New("formula.tmpl")
	t = template.Must(t.ParseGlob("templates/*.tmpl"))
	var buf bytes.Buffer
	err := t.Execute(&buf, f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return buf.Bytes()
}
