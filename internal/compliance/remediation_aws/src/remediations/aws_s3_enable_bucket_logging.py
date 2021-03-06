# Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
# Copyright (C) 2020 Panther Labs Inc
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

from typing import Any, Dict

from boto3 import Session

from ..app import Remediation
from ..app.remediation_base import RemediationBase


@Remediation
class AwsS3EnableBucketLogging(RemediationBase):
    """Remediation that enables S3 bucket access logging"""

    @classmethod
    def _id(cls) -> str:
        return 'S3.EnableBucketLogging'

    @classmethod
    def _parameters(cls) -> Dict[str, str]:
        return {'TargetBucket': '', 'TargetPrefix': ''}

    @classmethod
    def _fix(cls, session: Session, resource: Dict[str, Any], parameters: Dict[str, str]) -> None:
        session.client('s3').put_bucket_logging(
            Bucket=resource['Name'],
            BucketLoggingStatus={
                'LoggingEnabled': {
                    'TargetBucket': parameters['TargetBucket'],
                    'TargetPrefix': parameters['TargetPrefix']
                }
            }
        )
