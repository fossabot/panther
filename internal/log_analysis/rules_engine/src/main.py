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

import collections
import json
from gzip import GzipFile
from io import TextIOWrapper
from typing import Any, Dict, List

import boto3

from .engine import Engine
from .logging import get_logger
from .sqs import send_to_sqs

s3_client = boto3.client('s3')
rules_engine = Engine()


def lambda_handler(event: Dict[str, Any], unused_context) -> None:
    logger = get_logger()
    # Dictionary containing mapping from log type to list of TextIOWrapper's
    log_type_to_data: Dict[str, TextIOWrapper] = collections.defaultdict(list)

    for record in event['Records']:
        record_body = json.loads(record['body'])
        bucket = record_body['s3Bucket']
        object_key = record_body['s3ObjectKey']
        logger.info("loading object from S3, bucket [{}], key [{}]".format(bucket, object_key))
        log_type_to_data[record_body['id']].append(load_contents(bucket, object_key))

    # List containing tuple of (rule_id, event) for matched events
    matched: List = []

    for log_type, data_streams in log_type_to_data.items():
        for data_stream in data_streams:
            for data in data_stream:
                for matched_rule in rules_engine.analyze(log_type, data):
                    matched.append((matched_rule, data))

    if len(matched) > 0:
        logger.info("sending {} matches".format(len(matched)))
        send_to_sqs(matched)
    else:
        logger.info("no matches found")


# Returns a TextIOWrapper for the S3 data. This makes sure that we don't have to keep all
# contents of S3 object in memory
def load_contents(bucket: str, key: str) -> TextIOWrapper:
    response = s3_client.get_object(Bucket=bucket, Key=key)
    gzipped = GzipFile(None, 'rb', fileobj=response['Body'])
    return TextIOWrapper(gzipped)
