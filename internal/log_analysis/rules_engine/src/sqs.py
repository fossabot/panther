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

import json
import os
from datetime import datetime
from typing import List, Dict

import boto3

# Max number of SQS messages inside an SQS batch
_MAX_MESSAGES = 10
# Max size of an SQS batch request
_MAX_MESSAGE_SIZE = 256 * 1000

sqs_resource = boto3.resource('sqs')
queue = sqs_resource.get_queue_by_name(QueueName=os.environ['ALERTS_QUEUE'])


def send_to_sqs(matches: List) -> None:
    """Send a tuple of (rule_id, event) to SQS."""
    messages = [match_to_sqs_entry_message(i) for i in matches]

    current_entries: List[Dict[str, str]] = []
    current_byte_size = 0

    for i in range(len(messages)):
        entry = {'Id': str(i), 'MessageBody': messages[i]}
        projected_size = current_byte_size + len(messages[i])
        projected_num_entries = len(current_entries) + 1
        if projected_num_entries > _MAX_MESSAGES or projected_size > _MAX_MESSAGE_SIZE:
            queue.send_messages(Entries=current_entries)
            current_entries = [entry]
            current_byte_size = len(messages[i])
        else:
            current_entries.append(entry)
            current_byte_size += len(messages[i])

    if len(current_entries) > 0:
        queue.send_messages(Entries=current_entries)

    return


def match_to_sqs_entry_message(match: (str, str)) -> str:
    notification = {'ruleId': match[0], 'event': match[1], 'timestamp': datetime.utcnow().strftime('%Y-%m-%dT%H:%M:%SZ')}
    return json.dumps(notification)
