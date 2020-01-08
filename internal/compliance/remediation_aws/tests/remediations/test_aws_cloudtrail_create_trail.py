# Copyright 2020 Panther Labs Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from unittest import mock, TestCase
from boto3 import Session
from ...src.remediations.aws_cloudtrail_create_trail import AwsCloudTrailCreateTrail


class TestAwsCloudTrailCreateTrail(TestCase):

    @mock.patch.object(Session, 'client')
    def test_fix(self, mock_session: mock.MagicMock) -> None:
        mock_client = mock.Mock()
        mock_session.return_value = mock_client
        parameters = {
            'Name': 'TestTrailName',
            'TargetBucketName': 'TestTargetBucketName',
            'TargetPrefix': 'TestTargetPrefix',
            'SnsTopicName': 'TestSnsTopicName',
            'IncludeGlobalServiceEvents': 'True',
            'IsMultiRegionTrail': 'True',
            'KmsKeyId': 'TestKmsKeyId',
            'IsOrganizationTrail': 'True'
        }

        AwsCloudTrailCreateTrail()._fix(Session, {}, parameters)
        mock_session.assert_called_once_with('cloudtrail')
        mock_client.create_trail.assert_called_once_with(
            Name='TestTrailName',
            S3BucketName='TestTargetBucketName',
            S3KeyPrefix='TestTargetPrefix',
            SnsTopicName='TestSnsTopicName',
            IncludeGlobalServiceEvents=True,
            IsMultiRegionTrail=True,
            EnableLogFileValidation=True,
            KmsKeyId='TestKmsKeyId',
            IsOrganizationTrail=True
        )
        mock_client.start_logging.assert_called_once_with(Name='TestTrailName')
