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

import logging
import os
import shutil
import tempfile
from typing import Any, Dict

from . import engine

_LOGGER = logging.getLogger()
_LOGGER.setLevel('INFO')
_TMP = os.path.join(tempfile.gettempdir(), 'analysis')


def lambda_handler(lambda_event: Dict[str, Any], unused_context: Any) -> Dict[str, Any]:
    """Entry point for the policy engine.

    Args:
        lambda_event: {
            ###### Compliance Evaluation ######
            'policies': [
                {
                    'body': 'def policy(resource): ...',
                    'id': 'BucketEncryptionEnabled',
                    'resourceTypes': ['AWS.S3.Bucket']  # can be empty for all resource types
                }
            ],
            'resources': [
                {
                    'attributes': { ... resource attributes ... },
                    'id': 'arn:aws:s3:::my-bucket',
                    'type': 'AWS.S3.Bucket'
                }
            ]
        }

    Returns:
        {
            ###### Compliance Evaluation ######
            'resources': [
                {
                    'id': 'arn:aws:s3:::my-bucket',
                    'errored': [  # policies which raised a runtime error
                        {
                            'id': 'policy-id-1',
                            'message': 'ZeroDivisionError'
                        }
                    ],
                    'failed': ['policy-id-2', 'policy-id-3'],  # policies which returned False
                    'passed': ['policy-id-3', 'policy-id-4'],  # policies which returned True
                }
            ]
        }
    """
    if lambda_event.get('resources') is not None and lambda_event.get('policies') is not None:
        _LOGGER.info('Scanning %d resources with %d compliance policies', len(lambda_event['resources']), len(lambda_event['policies']))
    else:
        raise ValueError('resources and policies much be specified')

    # Erase /tmp at the beginning of every invocation.
    if not os.path.exists(_TMP):
        os.makedirs(_TMP)

    for name in os.listdir(_TMP):
        path = os.path.join(_TMP, name)
        if os.path.isdir(path):
            shutil.rmtree(path)
        else:
            os.remove(path)

    # Save all policies to /tmp for easy import.
    for policy in lambda_event['policies']:
        # Sanitize filename: replace all special characters with underscores.
        safe_id = ''.join(x if _allowed_char(x) else '_' for x in policy['id'])
        path = os.path.join(_TMP, safe_id + '.py')
        if os.path.exists(path):
            raise NameError('policy with sanitized id {} already exists'.format(safe_id))
        with open(path, 'w') as py_file:
            py_file.write(policy['body'])
        policy['body'] = path  # Replace policy body with file path.

    return engine.analyze(lambda_event)


def _allowed_char(char: str) -> bool:
    """Return true if the character is part of a valid policy ID."""
    return char.isalnum() or char in {' ', '-', '.'}
