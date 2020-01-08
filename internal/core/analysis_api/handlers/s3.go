package handlers

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
	"bytes"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/gateway/analysis/models"
)

// Delete one or more policies from S3.
//
// It is the caller's responsibility to ensure there are not more than 1000 policies in the request.
func s3BatchDelete(input *models.DeletePolicies) error {
	objects := make([]*s3.ObjectIdentifier, len(input.Policies))
	for i, entry := range input.Policies {
		objects[i] = &s3.ObjectIdentifier{Key: aws.String(string(entry.ID))}
	}

	_, err := s3Client.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: &env.Bucket,
		Delete: &s3.Delete{Objects: objects},
	})
	if err != nil {
		zap.L().Error("s3Client.DeleteObjects failed", zap.Error(err))
		return err
	}

	return nil
}

// Load a policy from the S3 bucket.
func s3Get(policyID models.ID, versionID models.VersionID) (*tableItem, error) {
	result, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket:    &env.Bucket,
		Key:       aws.String(string(policyID)),
		VersionId: aws.String(string(versionID)),
	})
	if err != nil {
		zap.L().Error("s3Client.GetObject failed", zap.Error(err))
		return nil, err
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		zap.L().Error("ioutil.ReadAll failed", zap.Error(err))
		return nil, err
	}

	var policy tableItem
	if err = jsoniter.Unmarshal(body, &policy); err != nil {
		zap.L().Error("policy unmarshal failed", zap.Error(err))
		return nil, err
	}
	policy.VersionID = versionID

	return &policy, nil
}

// Upload a policy to S3 and set the VersionID accordingly.
func s3Upload(policy *tableItem) error {
	// We don't need to store auto-generated fields - keep the S3 copy clean and minimal
	policy.LowerDisplayName = ""
	policy.LowerID = ""
	policy.LowerTags = nil
	policy.VersionID = ""

	body, err := jsoniter.Marshal(policy)
	if err != nil {
		zap.L().Error("policy marshal failed", zap.Error(err))
		return err
	}

	result, err := s3Client.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(body),
		Bucket: &env.Bucket,
		Key:    aws.String(string(policy.ID)),
	})
	if err != nil {
		zap.L().Error("s3Client.PutObject failed", zap.Error(err))
		return err
	}

	policy.VersionID = models.VersionID(*result.VersionId)
	return nil
}
