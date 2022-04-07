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

var goCmd = mg.GoCmd()

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
	cv, err := semver.LatestTag(".")
	if err != nil {
		return err
	}

	color.Green(cv)
	return nil
}

// Release tag version [major, minor, patch]
func Release(typ string) error {
	cv, err := semver.LatestTag(".")
	if err != nil {
		return err
	}

	nv, err := semver.Bump(cv, semver.ParseBumpType(typ))
	if err != nil {
		return err
	}

	err = git.CreateTag(nv,
		git.WithCreateTagMessage("release "+nv),
		git.WithCreateTagPushProgress(os.Stdout),
	)
	if err == nil {
		color.Green(nv)
	}
	return err
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
