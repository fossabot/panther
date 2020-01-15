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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/magefile/mage/mg"
)

const sourceLicense = "docs/LICENSE_HEADER.txt"

var (
	// Base source paths where license headers will be added
	licensePaths = []string{
		"api", "deployments", "internal", "pkg", "tools", "web/scripts", "web/src", "magefile.go"}

	// Don't apply the license header to some files used in an integration test
	licenseExceptions = regexp.MustCompile(`internal/core/analysis_api/main/test_policies/.+`)
)

// Add a comment character in front of each line in a block of license text.
func commentEachLine(prefix, text string) string {
	lines := strings.Split(text, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			result = append(result, prefix)
		} else {
			result = append(result, prefix+" "+line)
		}
	}

	return strings.Join(result, "\n")
}

func addSourceLicenses(basePaths ...string) error {
	rawHeader, err := ioutil.ReadFile(sourceLicense)
	if err != nil {
		return err
	}
	header := strings.TrimSpace(string(rawHeader))

	asteriskLicense := "/**\n" + commentEachLine(" *", header) + "\n */"
	hashtagLicense := commentEachLine("#", header)

	for _, root := range basePaths {
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() || licenseExceptions.MatchString(path) {
				return nil // skip directories and excluded paths
			}

			return addFileLicense(path, asteriskLicense, hashtagLicense)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func addFileLicense(path, asteriskLicense, hashtagLicense string) error {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return modifyFile(path, func(contents string) string {
			return addGoLicense(contents, asteriskLicense)
		})
	case ".js", ".ts", ".tsx":
		return modifyFile(path, func(contents string) string {
			return prependHeader(contents, asteriskLicense)
		})
	case "dockerfile", ".py", ".sh", ".yml", ".yaml":
		return modifyFile(path, func(contents string) string {
			return prependHeader(contents, hashtagLicense)
		})
	default:
		return nil
	}
}

// Rewrite file contents on disk with the given modifier function.
func modifyFile(path string, modifier func(string) string) error {
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	contents := string(contentBytes)

	newContents := modifier(contents)
	if newContents == contents {
		return nil // no changes required
	}

	if mg.Verbose() {
		fmt.Println("fmt: license: " + path)
	}

	return ioutil.WriteFile(path, []byte(newContents), 0644)
}

// Add the license to the given Go file contents if necessary, returning the modified body.
func addGoLicense(contents, asteriskLicense string) string {
	if strings.Contains(contents, asteriskLicense) {
		return contents
	}

	// Loop over each line looking for the package declaration.
	// Comments before the package statement must be preserved for godoc and +build declarations.
	var result []string
	foundPackage := false
	for _, line := range strings.Split(contents, "\n") {
		result = append(result, line)
		if !foundPackage && strings.HasPrefix(strings.TrimSpace(line), "package ") {
			result = append(result, "\n"+asteriskLicense)
			foundPackage = true
		}
	}

	return strings.Join(result, "\n")
}

// Prepend a header if it doesn't already exist, returning the modified file contents.
func prependHeader(contents, header string) string {
	if strings.Contains(contents, header) {
		return contents
	}
	return header + "\n\n" + contents
}
