package main

import (
	"go/build"
	"os"
	"path/filepath"
)

func (f *Formula) GolangImportPath() string {
	cwd, _ := os.Getwd()
	rootPackage, _ := build.Import(".", cwd, 0)
	return rootPackage.ImportPath
}

func (f *Formula) GolangImportDir() string {
	return filepath.Dir(f.GolangImportPath())
}
