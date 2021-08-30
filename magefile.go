//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/iwaltgen/magex/pipe"
)

const (
	packageName = "github.com/iwaltgen/magex"
	version     = "0.1.0"
)

// Run lint
func Lint() error {
	return sh.RunV("golangci-lint", "run")
}

// Run test cases
func Test() error {
	mg.Deps(Lint)

	return pipe.ExecStdout(
		"go test ./... -timeout 10s -cover -json",
		"tparse -all",
	)
}
