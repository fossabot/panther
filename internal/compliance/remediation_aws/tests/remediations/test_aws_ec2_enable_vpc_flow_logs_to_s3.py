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
from ...src.remediations.aws_ec2_enable_vpc_flow_logs_to_s3 import AwsEc2EnableVpcFlowLogsToS3


class TestAwsEc2EnableVpcFlowLogsToS3(TestCase):

    @mock.patch.object(Session, 'client')
    def test_fix(self, mock_session: mock.MagicMock) -> None:
        mock_client = mock.Mock()
        mock_session.return_value = mock_client
        resource = {'Id': 'TestVpcId'}
        parameters = {'TargetBucketName': 'TestTargetBucketName', 'TargetPrefix': 'TestTargetPrefix', 'TrafficType': 'TestTrafficType'}
        mock_client.create_flow_logs.return_value = {}

        AwsEc2EnableVpcFlowLogsToS3()._fix(Session, resource, parameters)
        mock_session.assert_called_once_with('ec2')
        mock_client.create_flow_logs.assert_called_once_with(
            ResourceIds=['TestVpcId'],
            ResourceType='VPC',
            TrafficType='TestTrafficType',
            LogDestinationType='s3',
            LogDestination='arn:aws:s3:::TestTargetBucketName/TestTargetPrefix'
        )
