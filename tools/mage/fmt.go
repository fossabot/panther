package mage

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"bytes"
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
	pyTargets = []string{"internal/compliance/remediation_aws", "internal/compliance/policy_engine"}
)

// Fmt Format source files
func Fmt() error {
	fmt.Println("fmt: license")
	if err := addSourceLicenses(licensePaths...); err != nil {
		return err
	}

	fmt.Println("fmt: gofmt", strings.Join(goTargets, " "))

	// 1) gofmt to standardize the syntax formatting with code simplification (-s) flag
	if err := sh.Run("gofmt", append([]string{"-l", "-s", "-w"}, goTargets...)...); err != nil {
		return err
	}

	// 2) Remove empty newlines from import groups
	if err := removeAllImportNewlines(); err != nil {
		return err
	}

	// 3) Goimports to group imports into 3 sections
	args := append([]string{"-w", "-local=github.com/panther-labs/panther"}, goTargets...)
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

	var newLines [][]byte
	inImport := false
	for _, line := range bytes.Split(contents, []byte("\n")) {
		if inImport {
			if len(line) == 0 {
				continue // skip empty newlines in import groups
			}
			if line[0] == ')' { // gofmt always puts the ending paren on its own line
				inImport = false
			}
		} else if bytes.HasPrefix(line, []byte("import (")) {
			inImport = true
		}

		newLines = append(newLines, line)
	}

	return ioutil.WriteFile(path, bytes.Join(newLines, []byte("\n")), 0644)
}
