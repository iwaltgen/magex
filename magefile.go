//go:build mage

package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/iwaltgen/magex/dep"
	"github.com/iwaltgen/magex/script"
	"github.com/iwaltgen/magex/semver"
	"github.com/iwaltgen/magex/spinner"
)

const (
	packageName = "github.com/iwaltgen/magex"
	version     = "0.2.0"
)

type VERSION mg.Namespace

var gitcmd func(args ...string) error

func init() {
	gitcmd = sh.RunCmd("git")
}

// Run lint
func Lint() error {
	return sh.RunV("golangci-lint", "run")
}

// Run test cases
func Test() error {
	mg.Deps(Lint)

	return script.ExecStdout(
		"go test ./... -race -json -coverprofile codecov.out -covermode atomic",
		"tparse -all",
	)
}

// Show current version
func (VERSION) Show() {
	ver := "v" + version
	color.Green("version: %s", ver)
}

// Bump version
func (VERSION) Bump(typ string) error {
	current := version
	next, err := semver.Bump(current, typ)
	if err != nil {
		return err
	}

	files := []string{"magefile.go", "README.md"}
	for _, file := range files {
		if _, err := script.File(file).Replace(current, next).WriteFile(file); err != nil {
			return fmt.Errorf("failed to bump version `%s`: %w", file, err)
		}
	}

	for _, file := range files {
		if err := gitcmd("add", file); err != nil {
			return fmt.Errorf("failed to git add command `%s`: %w", file, err)
		}
	}

	color.Green("new version: %s", next)
	return gitcmd("commit", "-m", "chore: bump version")
}

// Create current version tag
func (VERSION) Tag() error {
	tag := "v" + version
	if err := gitcmd("tag", "-a", tag, "-m", tag+" release"); err != nil {
		return fmt.Errorf("failed to add git tag: %w", err)
	}
	return gitcmd("push", "origin", tag)
}

// Run install dependency tool
func Setup() error {
	defer spinner.Start(100 * time.Millisecond)()

	pkgs, err := dep.GlobImport("tool/deps.go")
	if err != nil {
		return fmt.Errorf("failed to load package import: %w", err)
	}

	args := []string{"install"}
	args = append(args, pkgs...)
	return sh.RunV(mg.GoCmd(), args...)
}
