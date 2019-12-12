package analysis

// Config defines the file format when parsing a bulk upload.
//
// YAML tags required because the YAML unmarshaller needs them
// JSON tags not present because the JSON unmarshaller is easy
type Config struct {
	AnalysisType              string            `yaml:"AnalysisType"`
	AutoRemediationID         string            `yaml:"AutoRemediationID"`
	AutoRemediationParameters map[string]string `yaml:"AutoRemediationParameters"`
	Description               string            `yaml:"Description"`
	DisplayName               string            `yaml:"DisplayName"`
	Enabled                   bool              `yaml:"Enabled"`
	Filename                  string            `yaml:"Filename"`
	PolicyID                  string            `yaml:"PolicyID"`
	ResourceTypes             []string          `yaml:"ResourceTypes"`
	Reference                 string            `yaml:"Reference"`
	Runbook                   string            `yaml:"Runbook"`
	Severity                  string            `yaml:"Severity"`
	Suppressions              []string          `yaml:"Suppressions"`
	Tags                      []string          `yaml:"Tags"`
	Tests                     []Test            `yaml:"Tests"`
}

// Test is a unit test definition when parsing policies in a bulk upload.
type Test struct {
	ExpectedResult bool        `yaml:"ExpectedResult"`
	Name           string      `yaml:"Name"`
	Resource       interface{} `yaml:"Resource"`
	ResourceType   string      `yaml:"ResourceType"`
}
