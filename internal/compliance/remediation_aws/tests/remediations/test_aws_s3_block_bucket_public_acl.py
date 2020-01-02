from unittest import mock, TestCase
from boto3 import Session
from ...src.remediations.aws_s3_block_bucket_public_acl import AwsS3BlockBucketPublicACL


class TestAwsS3BlockBucketPublicACL(TestCase):

    @mock.patch.object(Session, 'client')
    def test_fix(self, mock_session: mock.MagicMock) -> None:
        mock_client = mock.Mock()
        mock_session.return_value = mock_client
        resource = {'Name': 'TestName'}
        AwsS3BlockBucketPublicACL()._fix(Session, resource, {})
        mock_session.assert_called_once_with('s3')

        mock_client.put_bucket_acl.assert_called_with(
            Bucket='TestName',
            ACL='private',
        )