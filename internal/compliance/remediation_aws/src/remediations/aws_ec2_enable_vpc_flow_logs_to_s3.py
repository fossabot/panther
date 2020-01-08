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
class AwsEc2EnableVpcFlowLogsToS3(RemediationBase):
    """Remediation that enables VPC Flow logs to S3 bucket"""

    @classmethod
    def _id(cls) -> str:
        return 'EC2.EnableVpcFlowLogsToS3'

    @classmethod
    def _parameters(cls) -> Dict[str, str]:
        return {'TargetBucketName': '', 'TargetPrefix': '', 'TrafficType': 'ALL'}

    @classmethod
    def _fix(cls, session: Session, resource: Dict[str, Any], parameters: Dict[str, str]) -> None:
        response = session.client('ec2').create_flow_logs(
            ResourceIds=[
                resource['Id'],
            ],
            ResourceType='VPC',
            TrafficType=parameters['TrafficType'],
            LogDestinationType='s3',
            LogDestination='arn:aws:s3:::{}/{}'.format(parameters['TargetBucketName'], parameters['TargetPrefix'])
        )
        if 'Unsuccessful' in response:
            raise Exception(response['Unsuccessful'][0])
