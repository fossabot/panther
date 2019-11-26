package mage

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

var fmtTargets = []string{"api", "internal", "pkg", "tools", "magefile.go"}

// Fmt Format Go files
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
