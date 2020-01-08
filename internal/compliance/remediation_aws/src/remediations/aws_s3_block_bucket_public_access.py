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
class AwsS3BlockBucketPublicAccess(RemediationBase):
    """Remediation that puts an S3 bucket block public access configuration"""

    @classmethod
    def _id(cls) -> str:
        return 'S3.BlockBucketPublicAccess'

    @classmethod
    def _parameters(cls) -> Dict[str, str]:
        return {'BlockPublicAcls': 'true', 'IgnorePublicAcls': 'true', 'BlockPublicPolicy': 'true', 'RestrictPublicBuckets': 'true'}

    @classmethod
    def _fix(cls, session: Session, resource: Dict[str, Any], parameters: Dict[str, str]) -> None:
        session.client('s3').put_public_access_block(
            Bucket=resource['Name'],
            PublicAccessBlockConfiguration={
                'BlockPublicAcls': parameters['BlockPublicAcls'].lower() == 'true',
                'IgnorePublicAcls': parameters['IgnorePublicAcls'].lower() == 'true',
                'BlockPublicPolicy': parameters['BlockPublicPolicy'].lower() == 'true',
                'RestrictPublicBuckets': parameters['RestrictPublicBuckets'].lower() == 'true',
            },
        )
