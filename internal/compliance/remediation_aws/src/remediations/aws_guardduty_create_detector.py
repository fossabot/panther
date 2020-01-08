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
class AwsGuardDutyCreateDetector(RemediationBase):
    """Remediation that creates a GuardDuty detector if one doesn't exist"""

    @classmethod
    def _id(cls) -> str:
        return 'GuardDuty.CreateDetector'

    @classmethod
    def _parameters(cls) -> Dict[str, str]:
        return {'FindingPublishingFrequency': 'FIFTEEN_MINUTES'}

    @classmethod
    def _fix(cls, session: Session, resource: Dict[str, Any], parameters: Dict[str, str]) -> None:
        session.client("guardduty").create_detector(Enable=True, FindingPublishingFrequency=parameters['FindingPublishingFrequency'])
