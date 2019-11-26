package mage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"gopkg.in/yaml.v2"
)

const (
	pantherLambdaKey = "x-panther-lambda-cfn-resource" // top-level key in Swagger file
	space8           = "        "
)

var (
	apiDefinitionPattern = regexp.MustCompile(`DefinitionBody: api/[^#]+\.yml.*`)
)

// Environment variables used in the deploy command
type deployEnv struct {
	Bucket   string `required:"true"`
	Params   string
	Template string `required:"true"`
}

// Deploy Upload to BUCKET in AWS_REGION and deploy TEMPLATE with optional PARAMS
func Deploy() error {
	// Parse config environment variables.
	var config deployEnv
	if err := envconfig.Process("", &config); err != nil {
		return fmt.Errorf("invalid configuration: %s", err)
	}

	if err := Clean(); err != nil {
		return err
	}
	if err := Build.Lambda(Build{}); err != nil {
		return err
	}

	stack := "panther-" + strings.TrimSuffix(path.Base(config.Template), ".yml")
	cfn, err := transformTemplate("deploy", config.Template, config.Bucket, stack)
	if err != nil {
		return err
	}

	return cfnDeploy(cfn, config.Bucket, stack, config.Params)
}

// Upload resources to S3 and return the path to the final CloudFormation template.
func transformTemplate(cmd, inputFile, bucket, prefix string) (string, error) {
	outputDir := path.Join("out", path.Dir(inputFile))
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	// Transformation 1: Embed Swagger definitions (a no-op if no apis are defined)
	swaggerOut := path.Join(outputDir, "swagger."+path.Base(inputFile))
	if err := embedSwagger(cmd, inputFile, swaggerOut); err != nil {
		return "", err
	}

	// Transformation 2: Upload Lambda source
	// Note: "sam package" is just an alias for "aws cloudformation package"
	// There is no equivalent to this command in the AWS Go SDK.
	pkgOut := path.Join(outputDir, "package."+path.Base(inputFile))
	args := []string{"cloudformation", "package",
		"--output-template-file", pkgOut,
		"--s3-bucket", bucket,
		"--s3-prefix", prefix,
		"--template-file", swaggerOut,
	}

	if mg.Verbose() {
		return pkgOut, sh.Run("aws", args...)
	}

	// By default, just print a single progress message instead of several lines of explanation
	fmt.Printf("%s: cloudformation package %s => %s\n", cmd, inputFile, pkgOut)
	_, err := sh.Output("aws", args...)
	return pkgOut, err
}

// Deploy the final CloudFormation template to a dev account.
func cfnDeploy(templateFile, bucket, stack, params string) error {
	// Note: "sam deploy" is just an alias for "aws cloudformation deploy"
	if !mg.Verbose() {
		// Give some indication of progress for long-running commands if not in verbose mode
		fmt.Printf("deploy: cloudformation deploy %s => %s\n", templateFile, stack)
	}
	flags := []string{
		"cloudformation", "deploy",
		"--capabilities", "CAPABILITY_AUTO_EXPAND", "CAPABILITY_IAM", "CAPABILITY_NAMED_IAM",
		"--s3-bucket", bucket,
		"--stack-name", stack,
		"--template-file", templateFile,
	}
	if params != "" {
		flags = append(flags, "--parameter-overrides")
		flags = append(flags, strings.Split(params, " ")...)
	}

	return sh.Run("aws", flags...)
}

// Transform a CloudFormation template by embedding API Swagger definitions.
//
// TODO - add unit tests for this function
func embedSwagger(cmd, cfnSource, cfnDest string) error {
	cfn, err := ioutil.ReadFile(cfnSource)
	if err != nil {
		return fmt.Errorf("failed to open CloudFormation template %s: %s", cfnSource, err)
	}

	var errList []error
	newCfn := apiDefinitionPattern.ReplaceAllFunc(cfn, func(match []byte) []byte {
		filename := strings.TrimSpace(strings.Split(string(match), " ")[1])
		fmt.Printf("%s: embedding swagger DefinitionBody: %s\n", cmd, filename)

		body, err := loadSwagger(filename)
		if err != nil {
			errList = append(errList, err)
			return match // return the original string unmodified
		}

		return []byte("DefinitionBody:\n" + *body)
	})

	if err := JoinErrors("swagger-embed", errList); err != nil {
		return err
	}

	if err := ioutil.WriteFile(cfnDest, newCfn, 0644); err != nil {
		return fmt.Errorf("failed to write swaggered CloudFormation template %s: %s", cfnDest, err)
	}

	return nil
}

// Load and transform a Swagger file for embedding in CloudFormation.
func loadSwagger(filename string) (*string, error) {
	swagger, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open '%s': %s", filename, err)
	}

	var apiBody map[string]interface{}
	if err := yaml.Unmarshal(swagger, &apiBody); err != nil {
		return nil, fmt.Errorf("failed to parse swagger yaml '%s': %s", filename, err)
	}

	// Allow AWS_IAM authorization (i.e. AWS SIGv4 signatures).
	apiBody["securityDefinitions"] = map[string]interface{}{
		"sigv4": map[string]string{
			"type":                         "apiKey",
			"name":                         "Authorization",
			"in":                           "header",
			"x-amazon-apigateway-authtype": "awsSigv4",
		},
	}

	// API Gateway will validate all requests to the maximum possible extent.
	apiBody["x-amazon-apigateway-request-validators"] = map[string]interface{}{
		"validate-all": map[string]bool{
			"validateRequestParameters": true,
			"validateRequestBody":       true,
		},
	}

	functionResource := apiBody[pantherLambdaKey].(string)
	if functionResource == "" {
		return nil, fmt.Errorf("%s is required in '%s'", pantherLambdaKey, filename)
	}
	delete(apiBody, pantherLambdaKey)

	// Every method requires the same boilerplate settings: validation, sigv4, lambda integration
	for _, endpoints := range apiBody["paths"].(map[interface{}]interface{}) {
		for _, definition := range endpoints.(map[interface{}]interface{}) {
			def := definition.(map[interface{}]interface{})
			def["x-amazon-apigateway-integration"] = map[string]interface{}{
				"httpMethod":          "POST",
				"passthroughBehavior": "never",
				"type":                "aws_proxy",
				"uri": map[string]interface{}{
					"Fn::Sub": strings.Join([]string{
						"arn:aws:apigateway:${AWS::Region}:lambda:path",
						"2015-03-31",
						"functions",
						"arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:${" + functionResource + "}",
						"invocations",
					}, "/"),
				},
			}
			def["x-amazon-apigateway-request-validator"] = "validate-all"
			def["security"] = []map[string]interface{}{
				{"sigv4": []string{}},
			}

			// Replace integer response codes with strings (cfn doesn't support non-string keys).
			responses := def["responses"].(map[interface{}]interface{})
			for code, val := range responses {
				if intcode, ok := code.(int); ok {
					responses[strconv.Itoa(intcode)] = val
					delete(responses, code)
				}
			}
		}
	}

	newBody, err := yaml.Marshal(apiBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal modified yaml: %s", err)
	}

	// Add spaces for the correct indentation when embedding.
	result := space8 + strings.ReplaceAll(string(newBody), "\n", "\n"+space8)
	return &result, nil
}
