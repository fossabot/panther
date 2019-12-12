"""Panther policy engine."""
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
            ],

            ###### Log Analysis ######
            'rules': [
                {
                    'body': 'def rule(log): ...',
                    'id': 'UnauthorizedAccess',
                    'logTypes': ['AWS.CloudTrail']  # can be empty for all resource types
                }
            ],
            'events': [
                {
                    'data': { ... log data ... },
                    'id': 'abc-123',  # any ID unique for each event in the request
                    'type': 'AWS.CloudTrail'
                }
            ],
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

            ###### Log Analysis ######
            'events': [
                {
                    'id': 'abc-123',
                    'errored': [  # policies which raised a runtime error
                        {
                            'id': 'policy-id-1',
                            'message': 'ZeroDivisionError'
                        }
                    ],
                    'matched': ['policy-id-2', 'policy-id-3'],  # policies which returned True
                    'notMatched': ['policy-id-3', 'policy-id-4'],  # policies which returned False
                }
            ]
        }
    """
    compliance = True
    if lambda_event.get('resources') is not None and lambda_event.get('policies') is not None:
        _LOGGER.info('Scanning %d resources with %d compliance policies', len(lambda_event['resources']), len(lambda_event['policies']))
    elif lambda_event.get('rules') is not None and lambda_event.get('events') is not None:
        _LOGGER.info('Scanning %d events with %d log analysis rules', len(lambda_event['events']), len(lambda_event['rules']))
        compliance = False

        # Convert to the compliance format so we can use the exact same analysis.
        lambda_event['policies'] = [
            {
                'body': rule['body'],
                'id': rule['id'],
                'resourceTypes': rule.get('logTypes')
            } for rule in lambda_event['rules']
        ]
        lambda_event['resources'] = [
            {
                'attributes': event['data'],
                'id': event['id'],
                'type': event['type']
            } for event in lambda_event['events']
        ]
    else:
        raise ValueError('either resources/policies or rules/events must be specified')

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

    response = engine.analyze(lambda_event)

    if not compliance:
        # Convert to log analysis response
        response = {
            'events':
                [
                    {
                        'id': result['id'],
                        'errored': result['errored'],
                        'matched': result['passed'],  # returned True
                        'notMatched': result['failed'],  # returned False
                    } for result in response['resources']
                ]
        }

    return response


def _allowed_char(char: str) -> bool:
    """Return true if the character is part of a valid policy ID."""
    return char.isalnum() or char in {' ', '-', '.'}
