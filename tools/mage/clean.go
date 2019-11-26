package mage

import (
	"fmt"
	"os"
	"path/filepath"
)

// Clean Remove auto-generated build artifacts
func Clean() error {
	dirs := []string{"out"} // directories to remove

	// Remove generated Swagger client/models
	pkgs, err := filepath.Glob("api/*/client")
	if err != nil {
		return err
	}
	dirs = append(dirs, pkgs...)

	pkgs, err = filepath.Glob("api/*/models")
	if err != nil {
		return err
	}
	dirs = append(dirs, pkgs...)

	for _, pkg := range dirs {
		fmt.Println("clean: rm -r " + pkg)
		if err := os.RemoveAll(pkg); err != nil {
			return err
		}
	}

	return nil
}
