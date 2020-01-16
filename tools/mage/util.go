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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gopkg.in/yaml.v2"
)

// Open and parse a yaml file.
func loadYamlFile(path string) (map[string]interface{}, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open '%s': %s", path, err)
	}

	var result map[string]interface{}
	if err := yaml.Unmarshal(contents, &result); err != nil {
		return nil, fmt.Errorf("failed to parse yaml file '%s': %s", path, err)
	}

	return result, nil
}

// Convert a parsed yaml map into a string map.
func stringMap(input map[interface{}]interface{}) map[string]string {
	result := make(map[string]string, len(input))
	for key, val := range input {
		result[key.(string)] = fmt.Sprintf("%v", val)
	}
	return result
}

// Get CloudFormation stack outputs as a map.
func getStackOutputs(awsSession *session.Session, name string) (map[string]string, error) {
	cfnClient := cloudformation.New(awsSession)
	input := &cloudformation.DescribeStacksInput{StackName: &name}
	response, err := cfnClient.DescribeStacks(input)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(response.Stacks[0].Outputs))
	for _, output := range response.Stacks[0].Outputs {
		result[aws.StringValue(output.OutputKey)] = aws.StringValue(output.OutputValue)
	}

	return result, nil
}

// Upload a local file to S3.
func uploadFileToS3(
	awsSession *session.Session, path, bucket, key string, meta map[string]*string) (*s3manager.UploadOutput, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	uploader := s3manager.NewUploader(awsSession)
	fmt.Printf("deploy: uploading %s to s3://%s/%s\n", path, bucket, key)
	return uploader.Upload(&s3manager.UploadInput{
		Body:     file,
		Bucket:   &bucket,
		Key:      &key,
		Metadata: meta,
	})
}

// Prompt the user for a string input.
func promptUser(prompt string, validator func(string) error) string {
	var result string

	for {
		fmt.Print(prompt)
		if _, err := fmt.Scanln(&result); err != nil {
			fmt.Println(err) // empty line, for example
			continue
		}

		result = strings.TrimSpace(result)
		if err := validator(result); err != nil {
			fmt.Println(err)
			continue
		}

		return result
	}
}

// Ensure non-empty strings.
func nonemptyValidator(input string) error {
	if len(input) == 0 {
		return errors.New("error: input is blank, please try again")
	}
	return nil
}

// Very simple email validation to prevent obvious mistakes.
func emailValidator(email string) error {
	if len(email) >= 4 && strings.Contains(email, "@") && strings.Contains(email, ".") {
		return nil
	}

	return errors.New("error: invalid email: must be at least 4 characters and contain '@' and '.'")
}

// A wrapper around the basic `exec.Command`, running it immediately & adding streaming output to STDOUT/STDERR
func runCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg ...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

