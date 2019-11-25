// +build mage

package main

import (
	"errors"
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var fmtTargets = []string{"pkg", "magefile.go"}

// Format Go files
func Fmt() error {
	fmt.Println("fmt:", fmtTargets)

	// We use both goimports and gofmt because their feature sets do not completely overlap.
	// Goimports groups imports into 3 sections, gofmt has a code simplification (-s) flag.
	importArgs := append([]string{"-l", "-w", "-local=github.com/panther-labs/panther"}, fmtTargets...)
	if err := sh.Run("goimports", importArgs...); err != nil {
		return err
	}

	return sh.Run("gofmt", append([]string{"-l", "-s", "-w"}, fmtTargets...)...)
}

// Test contains targets for testing code syntax, style, and correctness.
type Test mg.Namespace

// LintErrors runs all lint checks and returns a combined list of errors.
func (t Test) LintErrors() []error {
	var errList []error

	if !mg.Verbose() {
		fmt.Println("test:lint: go vet")
	}
	if err := sh.RunV("go", "vet", "./..."); err != nil {
		errList = append(errList, fmt.Errorf("go vet failed: %s", err))
	}

	if !mg.Verbose() {
		fmt.Println("test:lint: golint")
	}
	if err := sh.RunV("golint", "-set_exit_status=1", "./..."); err != nil {
		errList = append(errList, fmt.Errorf("golint failed: %s", err))
	}

	// Validate Go formatting
	if !mg.Verbose() {
		fmt.Println("test:lint: formatting")
	}
	importArgs := append([]string{"-d", "-local=github.com/panther-labs/panther"}, fmtTargets...)
	output, err := sh.Output("goimports", importArgs...)
	if err != nil {
		errList = append(errList, fmt.Errorf("goimports failed: %s", err))
	} else if len(output) > 0 {
		errList = append(errList, fmt.Errorf("goimports diff: %d bytes", len(output)))
	}

	importArgs = append([]string{"-d", "-s"}, fmtTargets...)
	output, err = sh.Output("gofmt", importArgs...)
	if err != nil {
		errList = append(errList, fmt.Errorf("gofmt failed: %s", err))
	} else if len(output) > 0 {
		errList = append(errList, fmt.Errorf("gofmt diff: %d bytes", len(output)))
	}

	return errList
}

// JoinErrors formats multiple errors into a single error.
func JoinErrors(command string, errList []error) error {
	if len(errList) == 0 {
		return nil
	}

	errString := fmt.Sprintf("%s failed with %d error(s):", command, len(errList))
	for i, err := range errList {
		errString += fmt.Sprintf("\n\t[%d] %s", i+1, err)
	}
	return errors.New(errString)
}

// Check code style
func (t Test) Lint() error {
	return JoinErrors("lint", t.LintErrors())
}

// Run unit tests
func (Test) Unit() error {
	args := []string{"test", "-cover", "./..."}
	if mg.Verbose() {
		args = append(args, "-v")
	}

	return sh.RunV("go", args...)
}

// Run unit tests and view test coverage in HTML
func (t Test) Cover() error {
	if err := sh.RunV("go", "test", "-cover", "-coverprofile=.coverage", "./..."); err != nil {
		return err
	}

	return sh.Run("go", "tool", "cover", "-html=.coverage")
}

// Run all required checks
func (t Test) CI() error {
	if err := t.Unit(); err != nil {
		return err
	}
	return t.Lint()
}
