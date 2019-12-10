package mage

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/panther-labs/panther/pkg/shutil"
)

var buildEnv = map[string]string{"GOARCH": "amd64", "GOOS": "linux"}

// Build contains targets for compiling source code.
type Build mg.Namespace

// API Generate Go client/models from Swagger specs in api/
func (b Build) API() error {
	specs, err := filepath.Glob("api/*/api.yml")
	if err != nil {
		return err
	}

	for _, spec := range specs {
		if !mg.Verbose() {
			fmt.Println("build:api: swagger generate " + spec)
		}
		args := []string{"generate", "client", "-q", "-t", path.Dir(spec), "-f", spec}
		if err := sh.Run("swagger", args...); err != nil {
			return err
		}
	}

	return nil
}

// Lambda Compile all Lambda function source
func (b Build) Lambda() error {
	mg.Deps(b.API)

	var packages []string
	err := filepath.Walk("internal", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && strings.HasSuffix(path, "main") {
			packages = append(packages, path)
		}
		return err
	})
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		if err := buildPackage(pkg); err != nil {
			return err
		}
	}

	return nil
}

func buildPackage(pkg string) error {
	targetDir := path.Join("out", "bin", pkg)
	binary := path.Join(targetDir, "main")
	oldInfo, statErr := os.Stat(binary)
	oldHash, hashErr := shutil.SHA256(binary)

	if !mg.Verbose() {
		fmt.Println("build:lambda: go build " + targetDir)
	}
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}
	if err := sh.RunWith(buildEnv, "go", "build", "-ldflags", "-s -w", "-o", targetDir, "./"+pkg); err != nil {
		return err
	}

	if statErr == nil && hashErr == nil {
		if hash, err := shutil.SHA256(binary); err == nil && hash == oldHash {
			// Optimization - if the binary contents haven't changed, reset the last modified time.
			// "aws cloudformation package" re-uploads any binary whose modification time has changed,
			// even if the contents are identical. So this lets us skip any unmodified binaries, which can
			// significantly reduce the total deployment time if only one or two functions changed.
			//
			// With 5 unmodified Lambda functions, deploy:backend went from 146s => 109s with this fix.
			if mg.Verbose() {
				fmt.Printf("build:lambda: %s unchanged, reverting timestamp\n", binary)
			}
			modTime := oldInfo.ModTime()
			return os.Chtimes(binary, modTime, modTime)
		}
	}

	return nil
}
