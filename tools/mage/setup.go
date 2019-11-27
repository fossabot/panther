package mage

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/magefile/mage/sh"
)

const (
	cfnlintVersion  = "0.25.5"
	golangciVersion = "1.21.0"
	swaggerVersion  = "0.21.0"
)

var (
	golangciPkg   = fmt.Sprintf("golangci-lint-%s-linux-amd64", golangciVersion)
	golangciLinux = fmt.Sprintf(
		"https://github.com/golangci/golangci-lint/releases/download/v%s/%s.tar.gz",
		golangciVersion, golangciPkg)
	swaggerLinux = fmt.Sprintf(
		"https://github.com/go-swagger/go-swagger/releases/download/v%s/swagger_linux_amd64",
		swaggerVersion)
)

// Setup Install development dependencies
func Setup() error {
	if err := sh.RunV("pip3", "install", "--user", "cfn-lint=="+cfnlintVersion); err != nil {
		return err
	}

	// TODO - Update to CircleCI
	// Some libraries are only needed for development, not for CI
	if os.Getenv("CODEBUILD_CI") == "" {
		if err := sh.RunV("go", "get", "-u", "golang.org/x/tools/cmd/goimports"); err != nil {
			return err
		}
		if err := sh.RunV("pip3", "install", "--user", "--upgrade", "awscli"); err != nil {
			return err
		}
	}

	env, err := sh.Output("uname")
	if err != nil {
		return err
	}

	// Install swagger and golang-ci
	switch env {
	case "Darwin":
		if err := sh.RunV("brew", "tap", "go-swagger/go-swagger"); err != nil {
			return err
		}
		return sh.RunV("brew", "reinstall", "go-swagger", "golangci-lint")

	case "Linux":
		if err := sh.RunV("curl", "-o", "/usr/local/bin/swagger", "-fL", swaggerLinux); err != nil {
			return err
		}
		if err := sh.RunV("chmod", "+x", "/usr/local/bin/swagger"); err != nil {
			return err
		}

		// golang-ci
		if err := os.MkdirAll("/tmp/golangci", 0755); err != nil {
			return err
		}
		if err := sh.RunV("curl", "-o", "/tmp/golangci/ci.tar.gz", "-fL", golangciLinux); err != nil {
			return err
		}
		if err := sh.RunV("tar", "-xzvf", "/tmp/golangci/ci.tar.gz", "-C", "/tmp/golangci/"); err != nil {
			return err
		}
		return sh.RunV("mv", path.Join("/tmp/golangci", golangciPkg, "golangci-lint"), "/usr/local/bin")

	default:
		return errors.New("unknown environment: " + env)
	}
}
