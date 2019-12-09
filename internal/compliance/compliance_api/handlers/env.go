package handlers

type envConfig struct {
	ComplianceTable string `required:"true" split_words:"true"`
	IndexName       string `required:"true" split_words:"true"`
}

// Env is the parsed environment variables
var Env envConfig
