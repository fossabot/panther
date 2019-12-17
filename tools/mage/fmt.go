package mage

import (
	"fmt"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	goTargets = []string{"api", "internal", "pkg", "tools", "magefile.go"}
	pyTargets = []string{"internal/compliance/remediation_aws", "internal/core/analysis_engine"}
)

// Fmt Format source files
func Fmt() error {
	fmt.Println("fmt: gofmt", strings.Join(goTargets, " "))

	// We use both goimports and gofmt because their feature sets do not completely overlap.
	// Goimports groups imports into 3 sections, gofmt has a code simplification (-s) flag.
	args := append([]string{"-l", "-w", "-local=github.com/panther-labs/panther"}, goTargets...)
	if err := sh.Run("goimports", args...); err != nil {
		return err
	}
	if err := sh.RunV("gofmt", append([]string{"-l", "-s", "-w"}, goTargets...)...); err != nil {
		return err
	}

	fmt.Println("fmt: yapf", strings.Join(pyTargets, " "))
	args = []string{"--in-place", "--parallel", "--recursive"}
	if mg.Verbose() {
		args = append(args, "--verbose")
	}
	return sh.Run("venv/bin/yapf", append(args, pyTargets...)...)
}
