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
"""Unit tests for src/engine.py"""
import json
import os
import tempfile
import unittest
from unittest import mock

from ..src import engine

_TMP = tempfile.gettempdir()


@mock.patch.object(engine, 'print')
class TestEngine(unittest.TestCase):
    """Unit tests for engine subprocess."""

    @mock.patch.object(engine.sys.stdin, 'read', return_value='{"policies": [], "resources": []}')
    def test_empty(self, mock_read: mock.MagicMock, mock_print: mock.MagicMock) -> None:
        """Run with no policies and no resources."""
        engine.main()
        mock_read.assert_called_once()
        mock_print.assert_called_once()

        output = json.loads(mock_print.call_args[0][0])
        self.assertEqual({'resources': []}, output)

    @mock.patch.object(
        engine.sys.stdin,
        'read',
        return_value=json.dumps(
            {
                'policies':
                    [
                        {
                            # Import failure
                            'body': os.path.join(_TMP, 'panther-1.py'),
                            'id': 'panther-1',
                            'resourceTypes': ['AWS.CloudTrail']
                        },
                        {
                            # Global policy: return True if resource["key"] == "value"
                            'body': os.path.join(_TMP, 'panther-2.py'),
                            'id': 'panther-2',
                        },
                        {
                            # Return True for all AWS.CloudTrail resources (and imports a library)
                            'body': os.path.join(_TMP, 'panther-3.py'),
                            'id': 'panther-3',
                            'resourceTypes': ['AWS.CloudTrail']
                        },
                        {
                            # Return True if data contains {"hello": "world"}
                            'body': os.path.join(_TMP, 'panther-4.py'),
                            'id': 'panther-4',
                            'resourceTypes': ['AWS.CloudTrail', 'Github.Repo']
                        },
                        {
                            # Runtime exception
                            'body': os.path.join(_TMP, 'panther-5.py'),
                            'id': 'panther-5',
                            'resourceTypes': ['AWS.CloudTrail']
                        }
                    ],
                'resources':
                    [
                        {
                            'attributes': {
                                'key': 'value'
                            },
                            'id': 'my-trail',
                            'type': 'AWS.CloudTrail',
                        }, {
                            'attributes': {
                                'hello': 'world'
                            },
                            'id': 'github.com/panther-labs/policy-engine',
                            'type': 'Github.Repo'
                        }, {
                            'attributes': {
                                'hello': 'panther'
                            },
                            'id': 'github.com/panther-labs/policy-api',
                            'type': 'Github.Repo'
                        }
                    ],
            }
        )
    )
    def test_all(self, mock_read: mock.MagicMock, mock_print: mock.MagicMock) -> None:
        """Test a variety of policy conditions."""
        with open(os.path.join(_TMP, 'panther-1.py'), 'w') as policy_file:
            policy_file.write('def... initely not valid Python')
        with open(os.path.join(_TMP, 'panther-2.py'), 'w') as policy_file:
            policy_file.write('def policy(resource): return resource.get("key") == "value"')
        with open(os.path.join(_TMP, 'panther-3.py'), 'w') as policy_file:
            # random.random() is always less than 1, this is equivalent to "return True"
            policy_file.write('import random\ndef policy(resource): return random.random() <= 1')
        with open(os.path.join(_TMP, 'panther-4.py'), 'w') as policy_file:
            policy_file.write('def policy(resource): return resource.get("hello", "").startswith("wor")')
        with open(os.path.join(_TMP, 'panther-5.py'), 'w') as policy_file:
            policy_file.write('def policy(resource): return 0/0')

        engine.main()
        mock_read.assert_called_once()
        mock_print.assert_called_once()

        output = json.loads(mock_print.call_args[0][0])
        for result in output['resources']:  # sort each list in the results
            result['errored'] = list(sorted(result['errored'], key=lambda x: x['id']))
            result['failed'] = list(sorted(result['failed']))
            result['passed'] = list(sorted(result['passed']))

        expected = {
            'resources':
                [
                    {
                        'id': 'my-trail',
                        'errored':
                            [
                                {
                                    'id': 'panther-1',
                                    'message': 'SyntaxError: invalid syntax (panther-1.py, line 1)'
                                }, {
                                    'id': 'panther-5',
                                    'message': 'ZeroDivisionError: division by zero'
                                }
                            ],
                        'failed': ['panther-4'],
                        'passed': ['panther-2', 'panther-3']
                    },
                    {
                        'id': 'github.com/panther-labs/policy-engine',
                        'errored': [],
                        'failed': ['panther-2'],
                        'passed': ['panther-4']
                    },
                    {
                        'id': 'github.com/panther-labs/policy-api',
                        'errored': [],
                        'failed': ['panther-2', 'panther-4'],
                        'passed': []
                    },
                ]
        }
        self.assertEqual(expected, output)
