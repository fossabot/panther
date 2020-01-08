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
"""Integration test to invoke the live Lambda function."""
import json
import unittest
import boto3

_FUNCTION = 'panther-policy-engine'
_INPUT = {
    'policies':
        [
            {
                'body': 'def policy(resource): return resource["AnyExist"] is True',
                'id': 'cloudtrail//!#enabled',  # special characters in policy ID
                'resourceTypes': ['AWS.CloudTrailMeta']
            },
            {
                'body':
                    'from policyuniverse.arn import ARN\n' +
                    (  # test third-party import
                        'def policy(resource): return bool(resource["KmsKeyId"])'
                    ),
                'id': 'cloudtrail-encrypted',
                'resourceTypes': ['AWS.CloudTrail']
            },
            {
                'body': 'def... initely not valid Python',
                'id': 'syntax-error'
            },
            {
                'body': 'def policy(resource): return 0/0',
                'id': 'runtime-error'
            },
            {
                'body': 'import aws_globals\ndef policy(resource): return aws_globals.GLOBAL_TRUE',
                'id': 'import-test'
            },
            {
                'body': 'GLOBAL_TRUE=True\ndef policy(resource): return True',
                'id': 'aws_globals',
                'resourceTypes': ['AWS.Dummy.Type']
            }
        ],
    'resources':
        [
            {
                'attributes':
                    {
                        'CloudWatchLogsLogGroupArn': 'arn:aws:logs:us-west-2:123:log-group:Panther:*',
                        'CloudWatchLogsRoleArn': 'arn:aws:iam::123:role/PantherStreamCloudTrailCWL',
                        'HasCustomEventSelectors': True,
                        'HomeRegion': 'us-west-2',
                        'IncludeGlobalServiceEvents': True,
                        'IsMultiRegionTrail': True,
                        'IsOrganizationTrail': False,
                        'KmsKeyId': None,
                        'LogFileValidationEnabled': True,
                        'Name': 'PantherCloudTrail',
                        'S3BucketName': 'panther-test-cloudtrail',
                        'S3KeyPrefix': None,
                        'SnsTopicARN': None,
                        'SnsTopicName': None,
                        'TrailARN': 'arn:aws:cloudtrail:us-west-2:123:trail/PantherCloudTrail'
                    },
                'id': 'arn:aws:cloudtrail:us-west-2:123:trail/PantherCloudTrail',
                'type': 'AWS.CloudTrail'
            },
            {
                'attributes': {
                    'AnyExist': True
                },
                'id': 'arn:aws:cloudtrail:123:meta',
                'type': 'AWS.CloudTrailMeta'
            },
        ]
}


class IntegrationTest(unittest.TestCase):
    """Test the policy-engine by invoking the Lambda function."""

    def setUp(self) -> None:
        """Create the AWS client."""
        self._client = boto3.client('lambda')
        self.maxDiff = None  # pylint: disable=invalid-name

    def test_policy_engine(self) -> None:
        """Invoke the policy-engine and check the response."""
        response = self._client.invoke(FunctionName=_FUNCTION, Payload=json.dumps(_INPUT).encode('utf-8'))
        self.assertIsNone(response.get('FunctionError'))
        output = json.loads(response['Payload'].read())
        expected = {
            'resources':
                [
                    {
                        'id': 'arn:aws:cloudtrail:us-west-2:123:trail/PantherCloudTrail',
                        'errored':
                            [
                                {
                                    'id': 'syntax-error',
                                    'message': 'SyntaxError: invalid syntax (syntax-error.py, line 1)'
                                }, {
                                    'id': 'runtime-error',
                                    'message': 'ZeroDivisionError: division by zero'
                                }
                            ],
                        'failed': ['cloudtrail-encrypted'],
                        'passed': ['import-test']
                    },
                    {
                        'id': 'arn:aws:cloudtrail:123:meta',
                        'errored':
                            [
                                {
                                    'id': 'syntax-error',
                                    'message': 'SyntaxError: invalid syntax (syntax-error.py, line 1)'
                                }, {
                                    'id': 'runtime-error',
                                    'message': 'ZeroDivisionError: division by zero'
                                }
                            ],
                        'failed': [],
                        'passed': ['cloudtrail//!#enabled', 'import-test']
                    },
                ]
        }
        self.assertEqual(expected, output)


if __name__ == '__main__':
    # Ignore ResourceWarnings from unclosed boto connections
    unittest.main(warnings='ignore')
