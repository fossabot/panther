package mage

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
