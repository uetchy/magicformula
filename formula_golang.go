package main

import (
	"go/build"
	"os"
	"path/filepath"
)

type GoResource struct {
	Name     string // "cli"
	URL      string // "github.com/codegangsta/cli"
	Revision string // "b6f7aadbeb21ae18972577173ce175af83ce239d"
}

func (f *Formula) GolangImportPath() string {
	cwd, _ := os.Getwd()
	rootPackage, _ := build.Import(".", cwd, 0)
	return rootPackage.ImportPath
}

func (f *Formula) GolangImportDir() string {
	return filepath.Dir(f.GolangImportPath())
}
