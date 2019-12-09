"""Unit tests for src/main.py"""
import os
import unittest
from unittest import mock

from ..src import main


@mock.patch.object(main, '_LOGGER')
@mock.patch.dict(os.environ, {'LD_LIBRARY_PATH': '.', 'PATH': os.environ['PATH']})
class TestMain(unittest.TestCase):
    """Unit tests for lambda handler."""

    def test_empty_event(self, mock_logger: mock.MagicMock) -> None:
        """Empty Lambda event raises a ValueError"""
        with self.assertRaises(ValueError):
            main.lambda_handler({}, None)

        mock_logger.assert_not_called()

    def test_error(self, mock_logger: mock.MagicMock) -> None:
        """Error is logged and re-raised"""
        lambda_event = {'policies': [], 'resources': [1, 2, 3]}
        with self.assertRaises(TypeError):
            main.lambda_handler(lambda_event, None)

        mock_logger.assert_has_calls([
            mock.call.info('Scanning %d resources with %d compliance policies', 3, 0),
        ])

    def test_policy(self, mock_logger: mock.MagicMock) -> None:
        """Test the execution of policies."""
        lambda_event = {
            'policies':
                [
                    {
                        'body': 'def policy(resource): return True',
                        'id': 'panther-true'
                    },
                    {
                        'body': 'def policy(resource): return False',
                        'id': 'panther-false'
                    },
                ],
            'resources': [{
                'attributes': {
                    'key': 'value'
                },
                'id': 'my-trail',
                'type': 'AWS.CloudTrail'
            }]
        }
        result = main.lambda_handler(lambda_event, None)

        expected = {'resources': [{'id': 'my-trail', 'errored': [], 'failed': ['panther-false'], 'passed': ['panther-true']}]}
        self.assertEqual(expected, result)
        self.assertEqual(1, mock_logger.info.call_count)
        mock_logger.exception.assert_not_called()

    def test_rule(self, mock_logger: mock.MagicMock) -> None:
        """Test the execution of rules."""
        lambda_event = {
            'rules':
                [
                    {
                        'body': 'def rule(event): return True',
                        'id': 'panther-true'
                    },
                    {
                        'body': 'def rule(event): return False',
                        'id': 'panther-false',
                        'logTypes': ['AWS.CloudTrail', 'osquery.snapshot']
                    },
                ],
            'events':
                [
                    {
                        'data': {
                            'key': 'value'
                        },
                        'id': 'my-trail',
                        'type': 'AWS.CloudTrail'
                    }, {
                        'data': {
                            'key': 'value'
                        },
                        'id': 'my-trail-2',
                        'type': 'AWS.CloudTrail'
                    }
                ]
        }
        result = main.lambda_handler(lambda_event, None)

        expected = {
            'events':
                [
                    {
                        'id': 'my-trail',
                        'errored': [],
                        'matched': ['panther-true'],
                        'notMatched': ['panther-false']
                    }, {
                        'id': 'my-trail-2',
                        'errored': [],
                        'matched': ['panther-true'],
                        'notMatched': ['panther-false']
                    }
                ]
        }
        self.assertEqual(expected, result)
        self.assertEqual(1, mock_logger.info.call_count)
        mock_logger.exception.assert_not_called()

    def test_duplicate_id(self, mock_logger: mock.MagicMock) -> None:
        """Policies with duplicate sanitized ids raise an error."""
        lambda_event = {
            'policies':
                [
                    {
                        'body': 'def policy(resource): return True',
                        'id': 'panther/true'
                    },
                    {
                        'body': 'def policy(resource): return True',
                        'id': 'panther%true'
                    },
                ],
            'resources': []
        }
        with self.assertRaises(NameError):
            main.lambda_handler(lambda_event, None)

        mock_logger.info.assert_called()
