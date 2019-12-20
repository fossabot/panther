package mage

import (
	"os"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/registry"
)

// Generate Glue tables for log processor output as CloudFormation
func generateGlueTables() error {
	outDir := "out/deployments/log_analysis"
	err := os.MkdirAll(outDir, 0755)
	if err != nil {
		return err
	}
	return registry.GenerateGlueCloudFormation(outDir + "/gluetables.json")
}
