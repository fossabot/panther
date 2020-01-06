package mage

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/panther-labs/panther/pkg/shutil"
)

const (
	// CloudFormation templates + stacks
	applicationStack    = "panther-app"
	applicationTemplate = "deployments/template.yml"
	bucketStack         = "panther-buckets" // prereq stack with Panther S3 buckets
	bucketTemplate      = "deployments/core/buckets.yml"

	// Python layer
	layerSourceDir   = "out/pip/analysis/python"
	layerZipfile     = "out/layer.zip"
	layerS3ObjectKey = "layers/python-analysis.zip"
)

// NOTE: Mage ignores the first word of the comment if it matches the function name.
// So the comment below is intentionally "Deploy Deploy"

// Deploy Deploy application infrastructure
func Deploy() error {
	config, err := loadYamlFile(configFile)
	if err != nil {
		return err
	}

	bucketParams := flattenParameterValues(config["BucketsParameterValues"])
	if err = cfnDeploy(bucketTemplate, "", bucketStack, bucketParams); err != nil {
		return err
	}

	if err = Build.Lambda(Build{}); err != nil {
		return err
	}

	if err = generateGlueTables(); err != nil {
		return err
	}

	if err = embedAPISpecs(); err != nil {
		return err
	}

	awsSession, err := session.NewSession()
	if err != nil {
		return err
	}

	bucket, err := getSourceBucket(awsSession)
	if err != nil {
		return err
	}

	template, err := cfnPackage(applicationTemplate, bucket, applicationStack)
	if err != nil {
		return err
	}

	deployParams, err := getDeployParams(awsSession, config, bucket)
	if err != nil {
		return err
	}

	if err := cfnDeploy(template, bucket, applicationStack, deployParams); err != nil {
		return err
	}

	loadBalancer, err := getLoadBalancerURL(awsSession)
	if err != nil {
		return err
	}

	fmt.Printf("deploy: Panther URL = http://%s\n", loadBalancer)
	return nil

	// TODO - install the initial rule sets here
}

// Get the name of the source bucket from the buckets stack outputs.
func getSourceBucket(awsSession *session.Session) (string, error) {
	cfnClient := cloudformation.New(awsSession)
	input := &cloudformation.DescribeStacksInput{StackName: aws.String(bucketStack)}
	response, err := cfnClient.DescribeStacks(input)
	if err != nil {
		return "", err
	}

	for _, output := range response.Stacks[0].Outputs {
		if aws.StringValue(output.OutputKey) == "SourceBucketName" {
			return *output.OutputValue, nil
		}
	}

	return "", errors.New("SourceBucketName output not found in stack " + bucketStack)
}

// Generate the set of deploy parameters.
//
// This will prompt the user for required parameters if they are not set and
// also upload the layer zipfile unless a custom layer is specified.
func getDeployParams(awsSession *session.Session, config map[string]interface{}, bucket string) ([]string, error) {
	params := config["AppParameterValues"].(map[interface{}]interface{})

	// If no email is specified in the config file, prompt the user.
	if !validEmail(params["UserEmail"].(string)) {
		if err := promptRequiredValues(params); err != nil {
			return nil, err
		}
	}

	// If no custom Python layer is defined, then we need to build the default one.
	if params["PythonLayerVersionArn"].(string) == "" {
		// Convert libs from []interface{} to []string
		rawLibs := config["PipLayer"].([]interface{})
		libs := make([]string, len(rawLibs))
		for i, lib := range rawLibs {
			libs[i] = lib.(string)
		}

		version, err := uploadLayer(awsSession, libs, bucket, layerS3ObjectKey)
		if err != nil {
			return nil, err
		}
		params["PythonLayerKey"] = layerS3ObjectKey
		params["PythonLayerObjectVersion"] = version
	}

	return flattenParameterValues(params), nil
}

// Upload custom Python analysis layer to S3 (if it isn't already), returning version ID
func uploadLayer(awsSession *session.Session, libs []string, bucket, key string) (string, error) {
	s3Client := s3.New(awsSession)
	head, err := s3Client.HeadObject(&s3.HeadObjectInput{Bucket: &bucket, Key: &key})

	sort.Strings(libs)
	libString := strings.Join(libs, ",")
	if err == nil && aws.StringValue(head.Metadata["Libs"]) == libString {
		fmt.Printf("deploy: s3://%s/%s exists and is up to date\n", bucket, key)
		return *head.VersionId, nil
	}

	// The layer is re-uploaded only if it doesn't exist yet or the library versions changed.
	fmt.Println("deploy: downloading " + libString)
	if err := os.RemoveAll(layerSourceDir); err != nil {
		return "", err
	}
	if err := os.MkdirAll(layerSourceDir, 0755); err != nil {
		return "", err
	}
	args := append([]string{"install", "-t", layerSourceDir}, libs...)
	if err := sh.Run("pip3", args...); err != nil {
		return "", err
	}

	// The package structure needs to be:
	//
	// layer.zip
	// │ python/policyuniverse/
	// └ python/policyuniverse-VERSION.dist-info/
	//
	// https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html#configuration-layers-path
	if err := shutil.ZipDirectory(path.Dir(layerSourceDir), layerZipfile); err != nil {
		return "", err
	}

	// Upload to S3
	fmt.Printf("deploy: uploading %s to s3://%s/%s\n", layerZipfile, bucket, key)
	uploader := s3manager.NewUploader(awsSession)
	zipFile, err := os.Open(layerZipfile)
	if err != nil {
		return "", err
	}

	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:     zipFile,
		Bucket:   &bucket,
		Key:      &key,
		Metadata: map[string]*string{"Libs": &libString},
	})
	if err != nil {
		return "", err
	}
	return *result.VersionID, nil
}

// Upload resources to S3 and return the path to the modified CloudFormation template.
func cfnPackage(templateFile, bucket, stack string) (string, error) {
	outputDir := path.Join("out", path.Dir(templateFile))
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	// There is no equivalent to this command in the AWS Go SDK.
	pkgOut := path.Join(outputDir, "package."+path.Base(templateFile))
	args := []string{"cloudformation", "package",
		"--output-template-file", pkgOut,
		"--s3-bucket", bucket,
		"--s3-prefix", stack,
		"--template-file", templateFile,
	}

	if mg.Verbose() {
		return pkgOut, sh.Run("aws", args...)
	}

	// By default, just print a single progress message instead of several lines of explanation
	fmt.Printf("deploy: cloudformation package %s => %s\n", templateFile, pkgOut)
	_, err := sh.Output("aws", args...)
	return pkgOut, err
}

// Deploy the final CloudFormation template.
func cfnDeploy(templateFile, bucket, stack string, params []string) error {
	args := []string{
		"cloudformation", "deploy",
		"--capabilities", "CAPABILITY_AUTO_EXPAND", "CAPABILITY_IAM", "CAPABILITY_NAMED_IAM",
		"--no-fail-on-empty-changeset",
		"--stack-name", stack,
		"--template-file", templateFile,
	}
	if bucket != "" {
		args = append(args, "--s3-bucket", bucket, "--s3-prefix", stack)
	}
	if len(params) > 0 {
		args = append(args, "--parameter-overrides")
		args = append(args, params...)
	}

	if !mg.Verbose() {
		// Give some indication of progress for long-running commands if not in verbose mode
		fmt.Printf("deploy: cloudformation deploy %s => %s\n", templateFile, stack)
	}
	return sh.Run("aws", args...)
}

func getLoadBalancerURL(awsSession *session.Session) (string, error) {
	cfnClient := cloudformation.New(awsSession)
	input := &cloudformation.DescribeStacksInput{StackName: aws.String(applicationStack)}
	response, err := cfnClient.DescribeStacks(input)
	if err != nil {
		return "", err
	}

	for _, output := range response.Stacks[0].Outputs {
		if aws.StringValue(output.OutputKey) == "LoadBalancerUrl" {
			return *output.OutputValue, nil
		}
	}

	return "", errors.New("LoadBalancerUrl output not found in stack " + applicationStack)
}
