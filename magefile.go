// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/deanishe/awgo/util"
	"github.com/deanishe/awgo/util/build"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

var (
	info     *build.Info
	env      map[string]string
	ldflags  string
	workDir  string
	buildDir = "./build"
	distDir  = "./dist"
	iconsDir = "./icons"
)

var (
	green  = "03ae03"
	blue   = "5484f3"
	red    = "b00000"
	yellow = "f8ac30"
)

func init() {
	var err error
	if info, err = build.NewInfo(); err != nil {
		panic(err)
	}
	if workDir, err = os.Getwd(); err != nil {
		panic(err)
	}
	env = info.Env()
	env["API_KEY"] = os.Getenv("GOODREADS_API_KEY")
	env["API_SECRET"] = os.Getenv("GOODREADS_API_SECRET")
	env["VERSION"] = info.Version
	env["PKG"] = "main"
	ldflags = `-X "$PKG.version=$VERSION"`
}

func mod(args ...string) error {
	argv := append([]string{"mod"}, args...)
	return sh.RunWith(env, "go", argv...)
}

// Aliases are mage command aliases.
var Aliases = map[string]interface{}{
	"b": Build,
	"c": Clean,
	"d": Dist,
	"l": Link,
}

// Build builds workflow in ./build
func Build() error {
	mg.Deps(cleanBuild)
	// mg.Deps(Deps)
	fmt.Println("building ...")

	err := sh.RunWith(env,
		"go", "build",
		// "-tags", "$TAGS",
		"-ldflags", ldflags,
		"-o", "./build/alfred-services",
		".",
	)
	if err != nil {
		return err
	}

	globs := build.Globs(
		"*.png",
		"info.plist",
		"*.html",
		"README.md",
		"LICENCE.txt",
		"icons/*.png",
		"*.js",
	)

	return build.SymlinkGlobs(buildDir, globs...)
}

// Run run workflow
func Run() error {
	mg.Deps(Build)
	fmt.Println("running ...")
	return sh.RunWith(env, buildDir+"/alfred-services", "-h")
}

// Dist build an .alfredworkflow file in ./dist
func Dist() error {
	mg.SerialDeps(Clean, Build)
	p, err := build.Export(buildDir, distDir)
	if err != nil {
		return err
	}

	fmt.Printf("built workflow file %q\n", p)
	return nil
}

// Config display configuration
func Config() {
	fmt.Println("     Workflow name:", info.Name)
	fmt.Println("         Bundle ID:", info.BundleID)
	fmt.Println("  Workflow version:", info.Version)
	fmt.Println("  Preferences file:", info.AlfredPrefsBundle)
	fmt.Println("       Sync folder:", info.AlfredSyncDir)
	fmt.Println("Workflow directory:", info.AlfredWorkflowDir)
	fmt.Println("    Data directory:", info.DataDir)
	fmt.Println("   Cache directory:", info.CacheDir)
}

// Link symlinks ./build directory to Alfred's workflow directory.
func Link() error {
	mg.Deps(Build)

	fmt.Println("linking ./build to workflow directory ...")
	target := filepath.Join(info.AlfredWorkflowDir, info.BundleID)
	// fmt.Printf("target: %s\n", target)

	if util.PathExists(target) {
		fmt.Println("removing existing workflow ...")
	}
	// try to remove it anyway, as dangling symlinks register as existing
	if err := os.RemoveAll(target); err != nil && !os.IsNotExist(err) {
		return err
	}

	src, err := filepath.Abs(buildDir)
	if err != nil {
		return err
	}
	return build.Symlink(target, src, true)
}

// Deps ensure dependencies
func Deps() error {
	mg.Deps(cleanDeps)
	fmt.Println("downloading deps ...")
	return mod("download")
}

// Vendor copy dependencies to ./vendor
func Vendor() error {
	mg.Deps(Deps)
	fmt.Println("vendoring deps ...")
	return mod("vendor")
}

// Clean remove build files
func Clean() {
	fmt.Println("cleaning ...")
	mg.Deps(cleanBuild, cleanMage, cleanDeps)
}

func cleanDeps() error {
	return mod("tidy", "-v")
}

// remove & recreate directory
func cleanDir(name string) error {
	if err := sh.Rm(name); err != nil {
		return err
	}
	return os.MkdirAll(name, 0755)
}

/*
func cleanDir(name string, exclude ...string) error {

	if _, err := os.Stat(name); err != nil {
		return nil
	}

	infos, err := ioutil.ReadDir(name)
	if err != nil {
		return err
	}
	for _, fi := range infos {

		var match bool
		for _, glob := range exclude {
			if match, err = doublestar.Match(glob, fi.Name()); err != nil {
				return err
			} else if match {
				break
			}
		}

		if match {
			fmt.Printf("excluded: %s\n", fi.Name())
			continue
		}

		p := filepath.Join(name, fi.Name())
		if err := os.RemoveAll(p); err != nil {
			return err
		}
	}
	return nil
}
*/

func cleanBuild() error {
	return cleanDir(buildDir)
}

func cleanMage() error {
	return sh.Run("mage", "-clean")
}

// CleanIcons delete all generated icons from ./icons
func CleanIcons() error {
	return cleanDir(iconsDir)
}

// func gitVersion() string {
// 	s, _ := sh.Output("git", "describe", "--tags", "--always", "--abbrev=10")
// 	return s
// }
