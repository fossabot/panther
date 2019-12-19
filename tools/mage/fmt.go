package mage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

	// 1) gofmt to standardize the syntax formatting with code simplification (-s) flag
	args := append([]string{"-l", "-w", "-local=github.com/panther-labs/panther"}, goTargets...)
	if err := sh.Run("gofmt", append([]string{"-l", "-s", "-w"}, goTargets...)...); err != nil {
		return err
	}

	// 2) Remove empty newlines from import groups
	if err := removeAllImportNewlines(); err != nil {
		return err
	}

	// 3) Goimports to group imports into 3 sections
	if err := sh.Run("goimports", args...); err != nil {
		return err
	}

	fmt.Println("fmt: yapf", strings.Join(pyTargets, " "))
	args = []string{"--in-place", "--parallel", "--recursive"}
	if mg.Verbose() {
		args = append(args, "--verbose")
	}
	return sh.Run("venv/bin/yapf", append(args, pyTargets...)...)
}

func removeAllImportNewlines() error {
	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			return removeImportNewlines(path)
		}
		return nil
	})
}

// Remove empty newlines from formatted import groups so goimports will correctly group them.
func removeImportNewlines(path string) error {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var newLines []string
	inImport := false
	for _, line := range strings.Split(string(contents), "\n") {
		if inImport {
			if len(line) == 0 {
				continue // skip empty newlines in import groups
			}
			if line[0] == ')' { // gofmt always puts the ending paren on its own line
				inImport = false
			}
		} else if strings.HasPrefix(line, "import (") {
			inImport = true
		}

		newLines = append(newLines, line)
	}

	return ioutil.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0644)
}
