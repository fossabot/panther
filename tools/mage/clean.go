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
