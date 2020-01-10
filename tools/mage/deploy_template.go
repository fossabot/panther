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
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

const (
	pollInterval = 5 * time.Second // How long to wait in between requests to the CloudFormation service
	pollTimeout  = time.Hour       // Give up if CreateChangeSet or ExecuteChangeSet takes longer than this
)

// Deploy a CloudFormation template.
//
// This is our own implementation of "cloudformation deploy" from the AWS CLI.
// Here we have more control over the output and waiters.
func deployTemplate(awsSession *session.Session, templateFile, stack string, params map[string]string) error {
	changeSet, err := createChangeSet(awsSession, templateFile, stack, params)
	if err != nil {
		return err
	}
	if changeSet == "" {
		return nil // nothing to do
	}
	return executeChangeSet(awsSession, changeSet, stack)
}

// Create a CloudFormation change set, returning its name.
//
// If there are no pending changes, the change set is deleted and a blank name is returned.
func createChangeSet(awsSession *session.Session, templateFile, stack string, params map[string]string) (string, error) {
	// Change set name - username + unix time (must be unique)
	changeSetName := fmt.Sprintf("panther-%d", time.Now().UnixNano())

	// Change set type - CREATE if a new stack otherwise UPDATE
	client := cloudformation.New(awsSession)
	response, err := client.DescribeStacks(&cloudformation.DescribeStacksInput{StackName: &stack})
	changeSetType := "CREATE"
	if err == nil && len(response.Stacks) > 0 {
		changeSetType = "UPDATE"
	}

	parameters := make([]*cloudformation.Parameter, 0, len(params))
	for key, val := range params {
		parameters = append(parameters, &cloudformation.Parameter{
			ParameterKey:   aws.String(key),
			ParameterValue: aws.String(val),
		})
	}

	templateBody, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return "", err
	}

	createInput := &cloudformation.CreateChangeSetInput{
		Capabilities: []*string{
			aws.String("CAPABILITY_AUTO_EXPAND"),
			aws.String("CAPABILITY_IAM"),
			aws.String("CAPABILITY_NAMED_IAM"),
		},
		ChangeSetName: &changeSetName,
		ChangeSetType: &changeSetType,
		Parameters:    parameters,
		StackName:     &stack,
		Tags: []*cloudformation.Tag{
			// Tags are propagated to every supported resource in the stack
			{
				Key:   aws.String("Application"),
				Value: aws.String("Panther"),
			},
		},
		TemplateBody: aws.String(string(templateBody)),
	}

	if _, err = client.CreateChangeSet(createInput); err != nil {
		return "", err
	}

	// Wait for change set creation to finish
	describeInput := &cloudformation.DescribeChangeSetInput{ChangeSetName: &changeSetName, StackName: &stack}
	prevStatus := ""
	for start := time.Now(); time.Since(start) < pollTimeout; {
		response, err := client.DescribeChangeSet(describeInput)
		if err != nil {
			return "", err
		}

		status := aws.StringValue(response.Status)
		reason := aws.StringValue(response.StatusReason)
		if status == "FAILED" && strings.HasPrefix(reason, "The submitted information didn't contain changes") {
			fmt.Printf("deploy: %s: no changes needed\n", stack)
			_, err := client.DeleteChangeSet(&cloudformation.DeleteChangeSetInput{
				ChangeSetName: &changeSetName,
				StackName:     &stack,
			})
			return "", err
		}

		if status != prevStatus {
			fmt.Printf("deploy: %s: CreateChangeSet: %s\n", stack, status)
			prevStatus = status
		}

		switch status {
		case "CREATE_COMPLETE":
			return changeSetName, nil // success!
		case "FAILED":
			return "", fmt.Errorf("create-change-set failed: " + reason)
		default:
			time.Sleep(pollInterval)
		}
	}

	return "", fmt.Errorf("create-change-set failed: timeout %s", pollTimeout)
}

func executeChangeSet(awsSession *session.Session, changeSet, stack string) error {
	client := cloudformation.New(awsSession)
	_, err := client.ExecuteChangeSet(&cloudformation.ExecuteChangeSetInput{
		ChangeSetName: &changeSet,
		StackName:     &stack,
	})
	if err != nil {
		return err
	}

	// Wait for change set to finish.
	// We build our own waiter to handle both update + create and to show status progress.
	input := &cloudformation.DescribeStacksInput{StackName: &stack}
	prevStatus := ""
	for start := time.Now(); time.Since(start) < pollTimeout; {
		response, err := client.DescribeStacks(input)
		if err != nil || len(response.Stacks) == 0 {
			// Stack may not exist yet
			time.Sleep(pollInterval)
			continue
		}

		status := *response.Stacks[0].StackStatus
		if status != prevStatus {
			fmt.Printf("deploy: %s: ExecuteChangeSet: %s\n", stack, status)
			prevStatus = status
		}

		if status == "CREATE_COMPLETE" || status == "UPDATE_COMPLETE" {
			return nil // success!
		} else if strings.Contains(status, "IN_PROGRESS") {
			// TODO - show progress of nested stacks (e.g. % updated)
			time.Sleep(pollInterval)
		} else {
			return fmt.Errorf("execute-change-set failed: %s: %s",
				status, aws.StringValue(response.Stacks[0].StackStatusReason))
		}
	}

	return fmt.Errorf("execute-change-set failed: timeout %s", pollTimeout)
}
