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
	"os"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/registry"
	"github.com/panther-labs/panther/tools/cfngen/gluecf"
)

// Generate Glue tables for log processor output as CloudFormation
func generateGlueTables() (err error) {
	outDir := "out/deployments/log_analysis"
	err = os.MkdirAll(outDir, 0755)
	if err != nil {
		return
	}
	glueCfFileName := outDir + "/gluetables.json"

	glueCfFile, err := os.Create(glueCfFileName)
	if err != nil {
		return
	}
	defer func() {
		glueCfFile.Close()
	}()

	cf, err := gluecf.GenerateCloudFormation(registry.AvailableTables())
	if err != nil {
		return
	}
	_, err = glueCfFile.Write(cf)
	return
}
