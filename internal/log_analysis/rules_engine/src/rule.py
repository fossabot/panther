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

import os
import sys
import tempfile
from importlib import util as import_util
from pathlib import Path
from typing import Any, Dict, Union

from .logging import get_logger

_RULE_FOLDER = os.path.join(tempfile.gettempdir(), 'rules')

# Rule with ID 'aws_globals' contains common Python logic used by other rules
COMMON_MODULE_RULE_ID = 'aws_globals'


class Rule:
    """Panther rule metadata and imported module."""
    logger = get_logger()

    def __init__(self, rule_id: str, rule_body: str) -> None:
        """Import rule contents from disk.

        Args:
            rule_id: Unique rule identifier
            rule_body: The rule body
        """
        self.rule_id = rule_id

        self._import_error = None
        try:
            self.store_rule(rule_id, rule_body)
            self._module = self.import_rule_as_module(rule_id)
        except Exception as err:  # pylint: disable=broad-except
            self._import_error = err

    def store_rule(self, rule_id: str, rule_body: str) -> None:
        """Stores rule to disk."""
        path = self.rule_id_to_path(rule_id)
        self.logger.debug('storing rule in path {}'.format(path))

        ## Create dir if it doesn't exist
        Path(os.path.dirname(path)).mkdir(parents=True, exist_ok=True)
        with open(path, 'w') as py_file:
            py_file.write(rule_body)

    def import_rule_as_module(self, rule_id: str) -> Any:
        """Dynamically import a Python module from a file.

        See also: https://docs.python.org/3/library/importlib.html#importing-a-source-file-directly
        """

        path = self.rule_id_to_path(rule_id)
        spec = import_util.spec_from_file_location(rule_id, path)
        mod = import_util.module_from_spec(spec)
        spec.loader.exec_module(mod)  # type: ignore
        self.logger.debug('imported module {} from path {}'.format(rule_id, path))
        if rule_id == COMMON_MODULE_RULE_ID:
            self.logger.debug('imported global module {} from path {}'.format(rule_id, path))
            # Importing it as a shared module
            sys.modules[rule_id] = mod
        return mod

    def rule_id_to_path(self, rule_id: str) -> str:
        safe_id = ''.join(x if self.allowed_char(x) else '_' for x in rule_id)
        path = os.path.join(_RULE_FOLDER, safe_id + '.py')
        return path

    def allowed_char(self, char: str) -> bool:
        """Return true if the character is part of a valid rule ID."""
        return char.isalnum() or char in {' ', '-', '.'}

    def run(self, event: Dict[str, Any]) -> Union[bool, Exception]:
        """Analyze a log line with this rule and return True, False, or an error."""
        if self._import_error:
            return self._import_error

        try:
            # Python source should have a method called "rule"
            matched = self._module.rule(event)
        except Exception as err:  # pylint: disable=broad-except
            return err

        if not isinstance(matched, bool):
            return Exception('rule returned {}, expected bool'.format(type(matched).__name__))

        return matched
