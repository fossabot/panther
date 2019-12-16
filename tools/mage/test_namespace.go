package mage

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Test contains targets for testing code syntax, style, and correctness.
type Test mg.Namespace

var build = Build{}

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
	mg.Deps(build.API)
	var errs []error

	// go metalinting
	fmt.Println("test:lint: golang")
	args := []string{"run", "--timeout", "10m"}
	if mg.Verbose() {
		args = append(args, "-v")
	}
	if err := sh.RunV("golangci-lint", args...); err != nil {
		errs = append(errs, err)
	}

	// python yapf
	fmt.Println("test:lint: python")
	args = []string{"--diff", "--parallel", "--recursive"}
	if mg.Verbose() {
		args = append(args, "--verbose")
	}
	if output, err := sh.Output("venv/bin/yapf", append(args, pyTargets...)...); err != nil {
		errs = append(errs, fmt.Errorf("yapf diff: %d bytes (err: %v)", len(output), err))
	}

	// python bandit (security linting)
	args = []string{"--recursive"}
	if mg.Verbose() {
		args = append(args, "--verbose")
	} else {
		args = append(args, "--quiet")
	}
	if err := sh.RunV("venv/bin/bandit", append(args, pyTargets...)...); err != nil {
		errs = append(errs, err)
	}

	// python lint
	args = []string{"-j", "0", "--disable", "duplicate-code,fixme,too-few-public-methods", "--max-line-length", "140"}
	if mg.Verbose() {
		args = append(args, "--verbose")
	}
	if err := sh.RunV("venv/bin/pylint", append(args, pyTargets...)...); err != nil {
		errs = append(errs, err)
	}

	// python mypy (type check)
	args = []string{"--cache-dir", "out/.mypy_cache", "--disallow-untyped-defs", "--ignore-missing-imports", "--warn-unused-ignores"}
	if mg.Verbose() {
		args = append(args, "--verbose")
	}
	if err := sh.RunV("venv/bin/mypy", append(args, pyTargets...)...); err != nil {
		errs = append(errs, err)
	}

	// Lint CloudFormation
	fmt.Println("test:lint: CloudFormation")
	var templates []string
	err := filepath.Walk("deployments", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".yml") {
			templates = append(templates, path)
		}
		return err
	})
	if err != nil {
		errs = append(errs, fmt.Errorf("filepath.Walk(deployments) failed: %v", err))
	}
	if err := sh.RunV("venv/bin/cfn-lint", templates...); err != nil {
		errs = append(errs, err)
	}

	return JoinErrors("test:lint", errs)
}

// Unit Run unit tests
func (Test) Unit() error {
	mg.Deps(build.API)
	args := []string{"test", "-cover", "./..."}
	if mg.Verbose() {
		args = append(args, "-v")
	}

	fmt.Println("test:unit: go test")
	if err := sh.RunV("go", args...); err != nil {
		return err
	}

	args = []string{"-m", "unittest", "discover"}
	if mg.Verbose() {
		args = append(args, "--verbose")
	}
	for _, target := range pyTargets {
		args = append(args, path.Dir(target))
	}

	fmt.Println("test:unit python unittest")
	return sh.RunV("python3", args...)
}

// Cover Run Go unit tests and view test coverage in HTML
func (t Test) Cover() error {
	mg.Deps(build.API)
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
	if err := build.Lambda(); err != nil {
		return err
	}
	if err := t.Unit(); err != nil {
		return err
	}
	return t.Lint()
}

// Integration Run TestIntegration* for PKG (default: ./...)
func (t Test) Integration() error {
	mg.Deps(build.API)
	if err := sh.Run("go", "clean", "-testcache"); err != nil {
		return err
	}

	pkg := os.Getenv("PKG")
	if pkg == "" {
		pkg = "./..."
	}
	// Note: We do NOT run integration tests in parallel
	testArgs := []string{"test", pkg, "-run=TestIntegration*", "-p", "1"}
	if mg.Verbose() {
		testArgs = append(testArgs, "-v")
	} else {
		fmt.Println("test:integration: go test " + pkg + " -run=TestIntegration*")
	}

	if err := os.Setenv("INTEGRATION_TEST", "True"); err != nil {
		return err
	}
	defer os.Unsetenv("INTEGRATION_TEST")
	if err := sh.RunV("go", testArgs...); err != nil {
		return err
	}

	// Run Python integration tests unless a Go pkg is specified
	if os.Getenv("PKG") == "" {
		fmt.Println("test:integration: python engine")
		return sh.RunV("python3", "internal/core/analysis_engine/tests/integration.py")
	}
	return nil
}
