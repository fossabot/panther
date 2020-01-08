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

from typing import Any, Dict

from boto3 import Session

from ..app import Remediation
from ..app.remediation_base import RemediationBase
from ..app.exceptions import InvalidParameterException


@Remediation
class AwsS3EnableBucketEncryption(RemediationBase):
    """Remediation that enables encryption for an S3 bucket"""

    @classmethod
    def _id(cls) -> str:
        return 'S3.EnableBucketEncryption'

    @classmethod
    def _parameters(cls) -> Dict[str, str]:
        return {'SSEAlgorithm': 'AES256', 'KMSMasterKeyID': ''}

    @classmethod
    def _fix(cls, session: Session, resource: Dict[str, Any], parameters: Dict[str, str]) -> None:
        if parameters['SSEAlgorithm'] == 'AES256':
            session.client('s3').put_bucket_encryption(
                Bucket=resource['Name'],
                ServerSideEncryptionConfiguration={
                    'Rules': [{
                        'ApplyServerSideEncryptionByDefault': {
                            'SSEAlgorithm': 'AES256'
                        },
                    },],
                },
            )
        elif parameters['SSEAlgorithm'] == 'aws:kms':
            session.client('s3').put_bucket_encryption(
                Bucket=resource['Name'],
                ServerSideEncryptionConfiguration={
                    'Rules':
                        [
                            {
                                'ApplyServerSideEncryptionByDefault':
                                    {
                                        'SSEAlgorithm': "aws:kms",
                                        'KMSMasterKeyID': parameters['KMSMasterKeyID']
                                    },
                            },
                        ],
                },
            )
        else:
            raise InvalidParameterException("Invalid value {} for parameter {}".format(parameters['SSEAlgorithm'], 'SSEAlgorithm'))
