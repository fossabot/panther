package mage

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	golangciVersion = "1.21.0"
	swaggerVersion  = "0.21.0"

	binDir = "/usr/local/bin"
)

// Setup Install development dependencies
func Setup() error {
	fmt.Println("setup: installing python3 venv")
	if err := os.RemoveAll("venv"); err != nil {
		return err
	}
	if err := sh.RunV("python3", "-m", "venv", "venv"); err != nil {
		return err
	}
	args := []string{"install", "-r", "requirements.txt"}
	if !mg.Verbose() {
		args = append(args, "--quiet")
	}
	if err := sh.RunV("venv/bin/pip3", args...); err != nil {
		return err
	}

	// TODO - This will change when we switch to Circle CI
	// Some libraries are only needed for development, not for CI
	if os.Getenv("CODEBUILD_CI") == "" {
		fmt.Println("setup: installing goimports and awscli for local development")
		if err := sh.RunV("go", "get", "golang.org/x/tools/cmd/goimports"); err != nil {
			return err
		}
		args := []string{"install", "awscli"}
		if !mg.Verbose() {
			args = append(args, "--quiet")
		}
		if err := sh.RunV("pip3", args...); err != nil {
			return err
		}
	}

	env, err := sh.Output("uname")
	if err != nil {
		return err
	}

	if err := installSwagger(env); err != nil {
		return err
	}
	return installGolangCiLint(env)
}

func installSwagger(uname string) error {
	fmt.Println("setup: installing go-swagger")
	url := fmt.Sprintf("https://github.com/go-swagger/go-swagger/releases/download/v%s/swagger_%s_amd64",
		swaggerVersion, strings.ToLower(uname))
	binary := path.Join(binDir, "swagger")
	if err := sh.RunV("curl", "-o", binary, "-fL", url); err != nil {
		return err
	}
	return sh.RunV("chmod", "+x", binary)
}

func installGolangCiLint(uname string) error {
	fmt.Println("setup: installing golangci-lint")
	downloadDir := filepath.Join(os.TempDir(), "golangci")
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return err
	}

	pkg := fmt.Sprintf("golangci-lint-%s-%s-amd64", golangciVersion, strings.ToLower(uname))
	url := fmt.Sprintf("https://github.com/golangci/golangci-lint/releases/download/v%s/%s.tar.gz",
		golangciVersion, pkg)
	if err := sh.RunV("curl", "-o", path.Join(downloadDir, "ci.tar.gz"), "-fL", url); err != nil {
		return err
	}

	if err := sh.RunV("tar", "-xzvf", path.Join(downloadDir, "ci.tar.gz"), "-C", downloadDir); err != nil {
		return err
	}
	return sh.RunV("mv", path.Join(downloadDir, pkg, "golangci-lint"), binDir)
}
