package mage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Test contains targets for testing code syntax, style, and correctness.
type Test mg.Namespace

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
	// go metalinting
	args := []string{"run"}
	if mg.Verbose() {
		args = append(args, "-v")
	} else {
		fmt.Println("test:lint: golangci-lint run")
	}

	if err := sh.RunV("golangci-lint", args...); err != nil {
		return err
	}

	// Lint CloudFormation
	if !mg.Verbose() {
		fmt.Println("test:lint: cfn-lint deployments/")
	}
	var templates []string
	err := filepath.Walk("deployments", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".yml") {
			templates = append(templates, path)
		}
		return nil
	})

	if err != nil {
		return err
	}
	if err := sh.RunV("cfn-lint", templates...); err != nil {
		return err
	}

	return nil
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
