package mage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

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

	// Lint CloudFormation
	if !mg.Verbose() {
		fmt.Println("test:lint: cfn-lint deployments/")
	}
	templates, err := filepath.Glob("deployments/*/*.yml")
	if err != nil {
		errList = append(errList, fmt.Errorf("cfn-lint failed: glob: %s", err))
	}
	if err := sh.RunV("cfn-lint", templates...); err != nil {
		errList = append(errList, fmt.Errorf("cfn-lint failed: %s", err))
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

// Lint Check code style
func (t Test) Lint() error {
	return JoinErrors("lint", t.LintErrors())
}

// Unit Run unit tests
func (Test) Unit() error {
	args := []string{"test", "-cover", "./..."}
	if mg.Verbose() {
		args = append(args, "-v")
	}

	return sh.RunV("go", args...)
}

// Cover Run unit tests and view test coverage in HTML
func (t Test) Cover() error {
	if err := os.MkdirAll("out/", 0755); err != nil {
		return err
	}

	if err := sh.RunV("go", "test", "-cover", "-coverprofile=out/.coverage", "./..."); err != nil {
		return err
	}

	return sh.Run("go", "tool", "cover", "-html=out/.coverage")
}

// CI Run all required checks
func (t Test) CI() error {
	if err := Build.Lambda(Build{}); err != nil {
		return err
	}
	if err := t.Unit(); err != nil {
		return err
	}
	return t.Lint()
}

// Integration Run TestIntegration* for PKG (default: ./...)
func (t Test) Integration() error {
	if err := sh.Run("go", "clean", "-testcache"); err != nil {
		return err
	}

	pkg := os.Getenv("PKG")
	if pkg == "" {
		pkg = "./..."
	}
	testArgs := []string{"test", pkg, "-run=TestIntegration*"}
	if mg.Verbose() {
		testArgs = append(testArgs, "-v")
	} else {
		fmt.Println("test:integration: go test " + pkg + " -run=TestIntegration*")
	}

	if err := os.Setenv("INTEGRATION_TEST", "True"); err != nil {
		return err
	}
	defer os.Unsetenv("INTEGRATION_TEST")
	return sh.RunV("go", testArgs...)
}
