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
from ...src.remediations.aws_iam_update_account_password_policy import AwsIamUpdateAccountPasswordPolicy


class TestAwsIamUpdateAccountPasswordPolicy(TestCase):

    @mock.patch.object(Session, 'client')
    def test_fix(self, mock_session: mock.MagicMock) -> None:
        mock_client = mock.Mock()
        mock_session.return_value = mock_client
        parameters = {
            'MinimumPasswordLength': '10',
            'RequireSymbols': 'True',
            'RequireNumbers': 'True',
            'RequireUppercaseCharacters': 'True',
            'RequireLowercaseCharacters': 'True',
            'AllowUsersToChangePassword': 'True',
            'MaxPasswordAge': '100',
            'PasswordReusePrevention': '1000'
        }
        AwsIamUpdateAccountPasswordPolicy()._fix(Session, {}, parameters)
        mock_session.assert_called_once_with('iam')
        mock_client.update_account_password_policy.assert_called_once_with(
            MinimumPasswordLength=10,
            RequireSymbols=True,
            RequireNumbers=True,
            RequireUppercaseCharacters=True,
            RequireLowercaseCharacters=True,
            AllowUsersToChangePassword=True,
            MaxPasswordAge=100,
            PasswordReusePrevention=1000
        )
