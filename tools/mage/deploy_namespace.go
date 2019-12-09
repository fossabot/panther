package mage

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/panther-labs/panther/pkg/shutil"
)

const (
	// CloudFormation templates + stacks
	backendStack    = "panther-backend"
	backendTemplate = "deployments/backend.yml"
	bucketStack     = "panther-buckets" // prereq stack with Panther S3 buckets
	bucketTemplate  = "deployments/shared/buckets.yml"

	// Python layer
	layerSourceDir = "out/pip/analysis/python"
	layerZipfile   = "out/layer.zip"
)

// Modify this to update the pip libraries in the Python analysis layer
var pipLibs = []string{
	"policyuniverse==1.3.2.1",
}

// Deploy defines mage targets for deploying Panther infrastructure.
type Deploy mg.Namespace

// Pre Deploy prerequisite S3 buckets with optional PARAMS
func (Deploy) Pre() error {
	return cfnDeploy(bucketTemplate, "", bucketStack, nil)
}

// Backend Deploy backend infrastructure with optional PARAMS
func (Deploy) Backend() error {
	if err := Build.Lambda(Build{}); err != nil {
		return err
	}

	if err := embedSwaggerAll(); err != nil {
		return err
	}

	awsSession, err := session.NewSession()
	if err != nil {
		return err
	}

	bucket, err := GetSourceBucket(awsSession)
	if err != nil {
		return err
	}

	version, err := uploadLayer(awsSession, bucket, "layers/python-analysis.zip")
	if err != nil {
		return err
	}

	template, err := cfnPackage(backendTemplate, bucket, backendStack)
	if err != nil {
		return err
	}

	return cfnDeploy(template, bucket, backendStack, []string{"PythonLayerObjectVersion=" + version})
}

// GetSourceBucket returns the name of the Panther source S3 bucket for CloudFormation uploads.
func GetSourceBucket(awsSession *session.Session) (string, error) {
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

	return "", errors.New("SourceBucketName output not found in " + bucketStack)
}

// Upload custom Python analysis layer to S3 (if it isn't already), returning version ID
func uploadLayer(awsSession *session.Session, bucket, key string) (string, error) {
	s3Client := s3.New(awsSession)
	head, err := s3Client.HeadObject(&s3.HeadObjectInput{Bucket: &bucket, Key: &key})

	sort.Strings(pipLibs)
	libString := strings.Join(pipLibs, ",")
	if err == nil && aws.StringValue(head.Metadata["Libs"]) == libString {
		fmt.Printf("deploy:backend: s3://%s/%s exists and is up to date\n", bucket, key)
		return *head.VersionId, nil
	}

	// The layer is re-uploaded only if it doesn't exist yet or the library versions changed.
	fmt.Println("deploy:backend: downloading " + libString)
	if err := os.RemoveAll(layerSourceDir); err != nil {
		return "", err
	}
	if err := os.MkdirAll(layerSourceDir, 0755); err != nil {
		return "", err
	}
	args := append([]string{"install", "-t", layerSourceDir}, pipLibs...)
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
	fmt.Printf("deploy:backend: uploading %s to s3://%s/%s\n", layerZipfile, bucket, key)
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

// Deploy the final CloudFormation template to a dev account.
func cfnDeploy(templateFile, bucket, stack string, params []string) error {
	args := []string{
		"cloudformation", "deploy",
		"--capabilities", "CAPABILITY_AUTO_EXPAND", "CAPABILITY_IAM", "CAPABILITY_NAMED_IAM",
		"--stack-name", stack,
		"--template-file", templateFile,
	}
	if bucket != "" {
		args = append(args, "--s3-bucket", bucket, "--s3-prefix", stack)
	}

	// Combine parameters from PARAMS env variable with those passed to this function.
	if env := os.Getenv("PARAMS"); env != "" {
		params = append(params, strings.Split(env, " ")...)
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
