package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/thedhejavu/errauditor/errauditor"
)

var (
	flagSet = flag.NewFlagSet("errauditor", flag.ContinueOnError)
	logger  *logrus.Logger
)

type app struct {
	excludeDirs     []string
	excludePatterns []*regexp.Regexp
}

func main() {

	a := &app{}
	logger := logrus.New()

	lvl, err := logrus.ParseLevel("info")
	logger.SetFormatter(&logrus.TextFormatter{})
	if err != nil {
		log.Panic(err)
	}

	logger.SetLevel(lvl)
	os.Exit(a.run(flagSet.Args()))
}

func (a *app) run(args []string) int {
	err := a.check(args)
	if err != nil {
		logger.Errorf("failed to run with: %s", err)
		return 1
	}
	return 0
}

func (a *app) check(args []string) error {
	// exclude directories or files
	a.excludePatterns = make([]*regexp.Regexp, 0, len(a.excludeDirs))
	for _, d := range a.excludeDirs {
		p, err := regexp.Compile(d)
		if err != nil {
			return fmt.Errorf("failed to parse exclude dir pattern: %v", err)
		}
		a.excludePatterns = append(a.excludePatterns, p)
	}

	// TODO: Reduce allocation.
	var files, dirs, pkgs []string
	// Check all files recursively when no args given.
	if len(args) == 0 {
		dirs = append(dirs, allPackagesInFS("./...")...)
	}
	for _, arg := range args {
		if strings.HasSuffix(arg, "/...") && isDir(arg[:len(arg)-len("/...")]) {
			dirs = append(dirs, allPackagesInFS(arg)...)
		} else if isDir(arg) {
			dirs = append(dirs, arg)
		} else if exists(arg) {
			files = append(files, arg)
		} else {
			pkgs = append(pkgs, arg)
		}
	}

	for _, f := range files {
		err := a.checkFile(f)
		if err != nil {
			logger.Debugf("failed to checkFile: %s", err)
			continue
		}
	}
	for _, d := range dirs {
		err := a.checkDir(d)
		if err != nil {
			logger.Debugf("failed to checkDir: %s", err)
			continue
		}
	}
	for _, p := range pkgs {
		err := a.checkPackage(p)
		if err != nil {
			logger.Debugf("failed to checkPackage: %s", err)
			continue
		}
	}
	return nil
}

func (a *app) checkFile(path string) error {
	dir := filepath.Dir(path)
	for _, p := range a.excludePatterns {
		if p.MatchString(dir) {
			return nil
		}
	}

	src, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return nil
	}
	if len(f.Comments) > 0 && isGenerated(src) {
		return fmt.Errorf("%s is a generated file", path)
	}

	return errauditor.Run(f, fset)
}

// Copyright (c) 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.
func (a *app) checkDir(dirname string) error {
	for _, p := range a.excludePatterns {
		if p.MatchString(dirname) {
			return nil
		}
	}
	pkg, err := build.ImportDir(dirname, 0)
	if err != nil {
		if _, nogo := err.(*build.NoGoError); nogo {
			// Don't complain if the failure is due to no Go source files.
			return nil
		}
		return nil
	}
	return a.checkImportedPackage(pkg)
}

func (a *app) checkPackage(pkgname string) error {
	pkg, err := build.Import(pkgname, ".", 0)
	if err != nil {
		if _, nogo := err.(*build.NoGoError); nogo {
			// Don't complain if the failure is due to no Go source files.
			return nil
		}
		return nil
	}

	return a.checkImportedPackage(pkg)
}

func (a *app) checkImportedPackage(pkg *build.Package) (err error) {
	var files []string
	files = append(files, pkg.GoFiles...)
	files = append(files, pkg.CgoFiles...)
	files = append(files, pkg.TestGoFiles...)

	// TODO: Reduce allocation.
	if pkg.Dir != "." {
		for _, f := range files {
			err := a.checkFile(filepath.Join(pkg.Dir, f))
			if err != nil {
				logger.Debugf("failed to checkImportedPackage: %s", err)
				continue
			}
		}
	}
	return
}

func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// isGenerated reports whether the source file is generated code
// according the rules from https://golang.org/s/generatedcode.
func isGenerated(src []byte) bool {
	var (
		genHdr = []byte("// Code generated ")
		genFtr = []byte(" DO NOT EDIT.")
	)
	sc := bufio.NewScanner(bytes.NewReader(src))
	for sc.Scan() {
		b := sc.Bytes()
		if bytes.HasPrefix(b, genHdr) && bytes.HasSuffix(b, genFtr) && len(b) >= len(genHdr)+len(genFtr) {
			return true
		}
	}
	return false
}
