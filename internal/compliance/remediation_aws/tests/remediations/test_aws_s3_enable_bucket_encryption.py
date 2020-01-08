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

from ...src.app.exceptions import InvalidParameterException
from ...src.remediations.aws_s3_enable_bucket_encryption import AwsS3EnableBucketEncryption


class TestAwsS3EnableBucketEncryption(TestCase):

    @mock.patch.object(Session, 'client')
    def test_fix_aes256(self, mock_session: mock.MagicMock) -> None:
        mock_client = mock.Mock()
        mock_session.return_value = mock_client
        resource = {'Name': 'TestName'}
        parameters = {'SSEAlgorithm': 'AES256', 'KMSMasterKeyID': ''}
        AwsS3EnableBucketEncryption()._fix(Session, resource, parameters)
        mock_session.assert_called_once_with('s3')

        mock_client.put_bucket_encryption.assert_called_with(
            Bucket='TestName',
            ServerSideEncryptionConfiguration={
                'Rules': [{
                    'ApplyServerSideEncryptionByDefault': {
                        'SSEAlgorithm': 'AES256'
                    },
                },],
            },
        )

    @mock.patch.object(Session, 'client')
    def test_fix_kms(self, mock_session: mock.MagicMock) -> None:
        mock_client = mock.Mock()
        mock_session.return_value = mock_client
        resource = {'Name': 'TestName'}
        parameters = {'SSEAlgorithm': 'aws:kms', 'KMSMasterKeyID': '313e6a3d-57c7-4544-ba59-0fecaabaf7b2'}
        AwsS3EnableBucketEncryption()._fix(Session, resource, parameters)
        mock_session.assert_called_once_with('s3')

        mock_client.put_bucket_encryption.assert_called_with(
            Bucket='TestName',
            ServerSideEncryptionConfiguration={
                'Rules':
                    [
                        {
                            'ApplyServerSideEncryptionByDefault':
                                {
                                    'SSEAlgorithm': 'aws:kms',
                                    'KMSMasterKeyID': '313e6a3d-57c7-4544-ba59-0fecaabaf7b2'
                                },
                        },
                    ],
            },
        )

    @mock.patch.object(Session, 'client')
    def test_fix_unknown_algorithm(self, mock_session: mock.MagicMock) -> None:
        mock_client = mock.Mock()
        mock_session.return_value = mock_client
        resource = {'Name': 'TestName'}
        parameters = {'SSEAlgorithm': 'unknown'}

        self.assertRaises(InvalidParameterException, AwsS3EnableBucketEncryption()._fix, Session, resource, parameters)
        mock_session.assert_not_called()
