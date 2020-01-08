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


@Remediation
class AwsCloudTrailCreateTrail(RemediationBase):
    """Remediation that creates a new CloudTrail trail to S3"""

    @classmethod
    def _id(cls) -> str:
        return 'CloudTrail.CreateTrail'

    @classmethod
    def _parameters(cls) -> Dict[str, str]:
        return {
            'Name': 'AutoRemediationTrail',
            'TargetBucketName': '',
            'TargetPrefix': '',
            'SnsTopicName': '',
            'IsMultiRegionTrail': 'true',
            'KmsKeyId': '',
            'IncludeGlobalServiceEvents': 'true',
            'IsOrganizationTrail': 'false'
        }

    @classmethod
    def _fix(cls, session: Session, resource: Dict[str, Any], parameters: Dict[str, str]) -> None:
        client = session.client('cloudtrail')
        client.create_trail(
            Name=parameters['Name'],
            S3BucketName=parameters['TargetBucketName'],
            S3KeyPrefix=parameters['TargetPrefix'],
            SnsTopicName=parameters['SnsTopicName'],
            IncludeGlobalServiceEvents=parameters['IncludeGlobalServiceEvents'].lower() == 'true',
            IsMultiRegionTrail=parameters['IsMultiRegionTrail'].lower() == 'true',
            EnableLogFileValidation=True,
            KmsKeyId=parameters['KmsKeyId'],
            IsOrganizationTrail=parameters['IsOrganizationTrail'].lower() == 'true'
        )
        client.start_logging(Name=parameters['Name'])
