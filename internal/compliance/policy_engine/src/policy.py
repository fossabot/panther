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
"""Classes to represent a Panther policy and a collection of policies."""
import collections
import sys
from importlib import util as import_util
from typing import Any, Dict, List, Union

AWS_GLOBALS = 'aws_globals'


class Policy:
    """Panther policy metadata and imported module."""

    @staticmethod
    def import_module(module_name: str, path: str) -> Any:
        """Dynamically import a Python module from a file."""
        # https://docs.python.org/3/library/importlib.html#importing-a-source-file-directly
        spec = import_util.spec_from_file_location(module_name, path)
        mod = import_util.module_from_spec(spec)
        spec.loader.exec_module(mod)  # type: ignore
        return mod

    def __init__(self, policy_id: str, path: str) -> None:
        """Import policy contents from /tmp.

        Args:
            policy_id: Unique policy identifier
            path: Path to the /tmp .py file
        """
        self.policy_id = policy_id

        self._import_error = None
        try:
            self._module = self.import_module(policy_id, path)
        except Exception as err:  # pylint: disable=broad-except
            self._import_error = err

    def run(self, resource_attributes: Dict[str, Any]) -> Union[bool, Exception]:
        """Analyze a resource with this policy and return True, False, or an error."""
        if self._import_error:
            return self._import_error

        try:
            # Python source should have a method called "policy"
            matched = self._module.policy(resource_attributes)
        except Exception as err:  # pylint: disable=broad-except
            return err

        if not isinstance(matched, bool):
            return TypeError('policy returned {}, expected bool'.format(type(matched).__name__))

        return matched


# TODO: Support helpers
class PolicySet:
    """A collection of Panther policies."""

    def __init__(self, policies: List[Dict[str, Any]]) -> None:
        """Import all policies."""
        # For efficient lookup, map resource type to list of applicable policies.
        self._policies_by_type: Dict[str, List[Policy]] = collections.defaultdict(list)
        self._global_policies: List[Policy] = []  # List of policies that apply to all log types

        for index, raw_policy in enumerate(policies):
            if raw_policy['id'] == AWS_GLOBALS:
                sys.modules[AWS_GLOBALS] = Policy.import_module(AWS_GLOBALS, raw_policy['body'])
                del policies[index]
                break

        for raw_policy in policies:
            policy = Policy(raw_policy['id'], raw_policy['body'])
            resource_types = raw_policy.get('resourceTypes')

            if resource_types:
                for rtype in resource_types:
                    self._policies_by_type[rtype].append(policy)
            else:
                self._global_policies.append(policy)

    def analyze(self, resource: Dict[str, Any]) -> Dict[str, Any]:
        """Analyze a resource with this policy set.

        Returns:
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
        """
        errored: List[Dict[str, str]] = []
        failed: List[str] = []
        passed: List[str] = []

        for policy in self._policies_by_type[resource['type']] + self._global_policies:
            result = policy.run(resource['attributes'])
            if isinstance(result, Exception):
                errored.append({'id': policy.policy_id, 'message': '{}: {}'.format(type(result).__name__, result)})
            elif result is False:
                failed.append(policy.policy_id)
            else:
                passed.append(policy.policy_id)

        return {'id': resource['id'], 'errored': errored, 'failed': failed, 'passed': passed}