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
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	jsoniter "github.com/json-iterator/go"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/pkg/shutil"
)

const (
	configFile = "deployments/panther_config.yml"

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

	awsSession, err := session.NewSession()
	if err != nil {
		return err
	}

	bucketParams := config["BucketsParameterValues"].(map[interface{}]interface{})
	if err = deployTemplate(awsSession, bucketTemplate, bucketStack, stringMap(bucketParams)); err != nil {
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

	outputs, err := getStackOutputs(awsSession, bucketStack)
	if err != nil {
		return err
	}
	bucket := outputs["SourceBucketName"]

	template, err := cfnPackage(applicationTemplate, bucket, applicationStack)
	if err != nil {
		return err
	}

	deployParams, err := getDeployParams(awsSession, config, bucket)
	if err != nil {
		return err
	}

	if err = deployTemplate(awsSession, template, applicationStack, deployParams); err != nil {
		return err
	}

	outputs, err = getStackOutputs(awsSession, applicationStack)
	if err != nil {
		return err
	}

	if err = buildAndPushImageFromSource(awsSession, outputs); err != nil {
		return err
	}

	if err := enableTOTP(awsSession, outputs["UserPoolId"]); err != nil {
		return err
	}

	if err := inviteFirstUser(awsSession, outputs["UserPoolId"]); err != nil {
		return err
	}

	// TODO - underline link
	fmt.Printf("\nPanther URL = https://%s\n", outputs["LoadBalancerUrl"])
	return nil

	// TODO - install initial rule sets
}

// Generate the set of deploy parameters for the main application stack.
//
// This will first upload the layer zipfile unless a custom layer is specified.
func getDeployParams(awsSession *session.Session, config map[string]interface{}, bucket string) (map[string]string, error) {
	result := stringMap(config["AppParameterValues"].(map[interface{}]interface{}))

	// If no custom Python layer is defined, then we need to build the default one.
	if result["PythonLayerVersionArn"] == "" {
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
		result["PythonLayerKey"] = layerS3ObjectKey
		result["PythonLayerObjectVersion"] = version
	}

	if result["WebApplicationCertificateArn"] == "" {
		certificateArn, err := uploadLocalCertificate(awsSession)
		if err != nil {
			return nil, err
		}
		result["WebApplicationCertificateArn"] = certificateArn
	}

	// If you need to dynamically add some CloudFormation parameters before deploying, you would do
	// that here

	return result, nil
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
	result, err := uploadFileToS3(awsSession, layerZipfile, bucket, key, map[string]*string{"Libs": &libString})
	if err != nil {
		return "", err
	}
	return *result.VersionID, nil
}

// Upload resources to S3 and return the path to the modified CloudFormation template.
// TODO - replace this with our own to avoid relying on the aws cli
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

// Enable software 2FA for the Cognito user pool - this is not yet supported in CloudFormation.
func enableTOTP(awsSession *session.Session, userPoolID string) error {
	if mg.Verbose() {
		fmt.Printf("deploy: enabling TOTP for user pool %s\n", userPoolID)
	}

	client := cognitoidentityprovider.New(awsSession)
	_, err := client.SetUserPoolMfaConfig(&cognitoidentityprovider.SetUserPoolMfaConfigInput{
		MfaConfiguration: aws.String("ON"),
		SoftwareTokenMfaConfiguration: &cognitoidentityprovider.SoftwareTokenMfaConfigType{
			Enabled: aws.Bool(true),
		},
		UserPoolId: &userPoolID,
	})
	return err
}

// If the Admin group is empty (e.g. on the initial deploy), create the initial admin user.
func inviteFirstUser(awsSession *session.Session, userPoolID string) error {
	cognitoClient := cognitoidentityprovider.New(awsSession)
	group, err := cognitoClient.ListUsersInGroup(&cognitoidentityprovider.ListUsersInGroupInput{
		GroupName:  aws.String("Admin"),
		UserPoolId: &userPoolID,
	})
	if err != nil {
		return err
	}
	if len(group.Users) > 0 {
		return nil // an admin already exists - nothing to do
	}

	// Prompt the user for email + first/last name
	fmt.Println("\nSetting up initial Panther admin user...")
	firstName := promptUser("First name: ", nonemptyValidator)
	lastName := promptUser("Last name: ", nonemptyValidator)
	email := promptUser("Email: ", emailValidator)

	// Hit users-api.InviteUser to invite a new user to the admin group
	input := &models.LambdaInput{
		InviteUser: &models.InviteUserInput{
			GivenName:  &firstName,
			FamilyName: &lastName,
			Email:      &email,
			UserPoolID: &userPoolID,
			Role:       aws.String("Admin"),
		},
	}
	payload, err := jsoniter.Marshal(input)
	if err != nil {
		return err
	}

	lambdaClient := lambda.New(awsSession)
	response, err := lambdaClient.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String("panther-users-api"),
		Payload:      payload,
	})
	if err != nil {
		return err
	}

	if response.FunctionError != nil {
		return fmt.Errorf("failed to invoke panther-users-api: %s error: %s",
			*response.FunctionError, string(response.Payload))
	}

	return nil
}

// Functions that build a personalized docker image from source, while pushing it to the private image repo of the user
func buildAndPushImageFromSource(awsSession *session.Session, outputs map[string]string) error {

	fmt.Println("docker: Building docker image from source...")
	if err := runCommand("docker", "build",
		"--file", "deployments/web/Dockerfile",
		"--tag", outputs["WebApplicationImage"],
		"--build-arg", fmt.Sprintf("AWS_ACCOUNT_ID=%s", outputs["AWSAccountId"]),
		"--build-arg", fmt.Sprintf("AWS_REGION=%s", outputs["AWSRegion"]),
		"--build-arg", fmt.Sprintf("GRAPHQL_ENDPOINT=%s", outputs["WebApplicationGraphQLApiEndpoint"]),
		"--build-arg", fmt.Sprintf("AWS_COGNITO_USER_POOL_ID=%s", outputs["UserPoolId"]),
		"--build-arg", fmt.Sprintf("AWS_COGNITO_APP_CLIENT_ID=%s", outputs["UserPoolClientId"]),
		".",
	); err != nil {
		return err
	}

	fmt.Println("Requesting access to remote image repo...")
	ecrClient := ecr.New(awsSession)
	req, resp := ecrClient.GetAuthorizationTokenRequest(&ecr.GetAuthorizationTokenInput{})
	if err := req.Send(); err != nil {
		return err
	}

	ecrAuthorizationToken := *resp.AuthorizationData[0].AuthorizationToken
	ecrServer := *resp.AuthorizationData[0].ProxyEndpoint

	decodedCredentialsInBytes, _ := base64.StdEncoding.DecodeString(ecrAuthorizationToken)
	credentials := strings.Split(string(decodedCredentialsInBytes), ":")

	fmt.Println("Logging in to remote image repo...")
	if err := runCommand("docker", "login", "-u", credentials[0], "-p", credentials[1], ecrServer); err != nil {
		return err
	}

	fmt.Println("Begin image push to remote repo...")
	if err := runCommand("docker", "push", outputs["WebApplicationImage"]); err != nil {
		return err
	}

	fmt.Println("Image pushed successfully!")
	return nil
}
