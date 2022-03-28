//go:build mage

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/iwaltgen/magex/dep"
	"github.com/iwaltgen/magex/git"
	"github.com/iwaltgen/magex/script"
	"github.com/iwaltgen/magex/semver"
	"github.com/iwaltgen/magex/spinner"
)

var goCmd string

func init() {
	goCmd = mg.GoCmd()
}

// Run lint
func Lint() error {
	return sh.RunV("golangci-lint", "run")
}

// Run test cases
func Test() error {
	mg.Deps(Lint)

	return script.ExecStdout(
		goCmd+" test ./... -json -coverprofile codecov.out -covermode atomic",
		"tparse -all",
	)
}

// Show current version
func Version() error {
	version, err := currentVersion()
	if err != nil {
		return err
	}

	color.Green(version)
	return nil
}

// Release tag version [major, minor, patch]
func Release(typ string) error {
	current, err := currentVersion()
	if err != nil {
		return err
	}

	next, err := semver.Bump(current, typ)
	if err != nil {
		return err
	}

	return git.CreateTag(next,
		git.WithCreateTagMessage("release "+next),
		git.WithCreateTagProgress(os.Stdout),
	)
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
	return sh.RunV(goCmd, args...)
}

func currentVersion() (string, error) {
	tags, err := git.Tags(".")
	if err != nil {
		return "", err
	}

	latest, err := semver.Latest(tags)
	if err != nil {
		return "", err
	}

	return latest, nil
}
