package mage

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	// CloudFormation templates + stacks
	backendStack    = "panther-backend"
	backendTemplate = "deployments/backend.yml"
	bucketStack     = "panther-buckets" // prereq stack with Panther S3 buckets
	bucketTemplate  = "deployments/shared/buckets.yml"
)

// Deploy defines mage targets for deploying Panther infrastructure.
type Deploy mg.Namespace

// Pre Deploy prerequisite S3 buckets with optional PARAMS
func (Deploy) Pre() error {
	return cfnDeploy(bucketTemplate, "", bucketStack)
}

// Backend Deploy backend infrastructure with optional PARAMS
func (Deploy) Backend() error {
	if err := Build.Lambda(Build{}); err != nil {
		return err
	}

	if err := embedSwaggerAll(); err != nil {
		return err
	}

	bucket, err := GetSourceBucket()
	if err != nil {
		return err
	}

	template, err := cfnPackage(backendTemplate, bucket, backendStack)
	if err != nil {
		return err
	}

	return cfnDeploy(template, bucket, backendStack)
}

// GetSourceBucket returns the name of the Panther source S3 bucket for CloudFormation uploads.
func GetSourceBucket() (string, error) {
	cfnClient := cloudformation.New(session.Must(session.NewSession()))
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
func cfnDeploy(templateFile, bucket, stack string) error {
	args := []string{
		"cloudformation", "deploy",
		"--capabilities", "CAPABILITY_AUTO_EXPAND", "CAPABILITY_IAM", "CAPABILITY_NAMED_IAM",
		"--stack-name", stack,
		"--template-file", templateFile,
	}
	if bucket != "" {
		args = append(args, "--s3-bucket", bucket, "--s3-prefix", stack)
	}
	if params := os.Getenv("PARAMS"); params != "" {
		args = append(args, "--parameter-overrides")
		args = append(args, strings.Split(params, " ")...)
	}

	if !mg.Verbose() {
		// Give some indication of progress for long-running commands if not in verbose mode
		fmt.Printf("deploy: cloudformation deploy %s => %s\n", templateFile, stack)
	}
	return sh.Run("aws", args...)
}
