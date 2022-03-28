//go:build mage

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/iwaltgen/magex/dep"
	"github.com/iwaltgen/magex/script"
	"github.com/iwaltgen/magex/semver"
	"github.com/iwaltgen/magex/spinner"
)

const (
	packageName = "github.com/iwaltgen/magex"
	version     = "0.5.4"
)

type VERSION mg.Namespace

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
func (VERSION) Show() {
	ver := "v" + version
	color.Green("version: %s", ver)
}

// Bump version
func (ns VERSION) Bump(typ string) error {
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

	worktree, err := ns.worktree()
	if err != nil {
		return err
	}

	for _, file := range files {
		if _, err := worktree.Add(file); err != nil {
			return fmt.Errorf("failed to git add command `%s`: %w", file, err)
		}
	}

	hash, err := worktree.Commit("chore: bump version", &git.CommitOptions{})
	color.Green("new version: %s [%s]", next, hash.String())
	return err
}

// Create current version tag
func (ns VERSION) Tag() error {
	tag := "v" + version
	repo, err := git.PlainOpen(".")
	if err != nil {
		return err
	}

	head, err := repo.Head()
	if err != nil {
		return err
	}

	_, err = repo.CreateTag(tag, head.Hash(), &git.CreateTagOptions{
		Message: tag + " release",
	})
	if err != nil {
		return fmt.Errorf("failed to add git tag: %w", err)
	}

	return repo.Push(&git.PushOptions{
		Progress: os.Stdout,
		RefSpecs: []config.RefSpec{
			config.RefSpec(fmt.Sprintf("refs/tags/%[1]s:refs/tags/%[1]s", tag)),
		},
	})
}

// Show latest version tag
func (ns VERSION) Tags() error {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return err
	}

	refs, err := repo.Tags()
	if err != nil {
		return err
	}

	return refs.ForEach(func(r *plumbing.Reference) error {
		fmt.Println(r.Name().Short())
		return nil
	})
}

func (VERSION) worktree() (*git.Worktree, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return nil, err
	}

	return repo.Worktree()
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
