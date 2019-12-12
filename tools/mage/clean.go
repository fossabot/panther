package mage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Clean Remove auto-generated build artifacts
func Clean() error {
	dirs := []string{"out"} // directories to remove

	// Remove generated Swagger client/models
	pkgs, err := filepath.Glob("api/gateway/*/client")
	if err != nil {
		return err
	}
	dirs = append(dirs, pkgs...)

	pkgs, err = filepath.Glob("api/gateway/*/models")
	if err != nil {
		return err
	}
	dirs = append(dirs, pkgs...)

	// Remove __pycache__ folderrs
	for _, target := range pyTargets {
		err = filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, "__pycache__") {
				dirs = append(dirs, path)
			}
			return err
		})
		if err != nil {
			return err
		}
	}

	for _, pkg := range dirs {
		fmt.Println("clean: rm -r " + pkg)
		if err := os.RemoveAll(pkg); err != nil {
			return err
		}
	}

	return nil
}
