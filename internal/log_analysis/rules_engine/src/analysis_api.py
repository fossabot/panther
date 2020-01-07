import json
import os
from typing import Dict, List

import requests
from botocore.auth import SigV4Auth
from botocore.awsrequest import AWSRequest
from botocore.session import Session


class AnalysisAPIClient:
    """Client for interacting with Analysis API."""

    def __init__(self) -> None:
        current_session = Session()
        region = current_session.get_config_variable('region')
        creds = current_session.get_credentials()
        self.signer = SigV4Auth(creds, 'execute-api', region)

        analysis_api_fqdn = os.environ.get('ANALYSIS_API_FQDN')
        analysis_api_path = os.environ.get('ANALYSIS_API_PATH')
        self.url = 'https://' + analysis_api_fqdn + '/' + analysis_api_path

    def get_enabled_rules(self) -> List[Dict[str, str]]:
        """Gets information for all enabled rules."""
        request = AWSRequest(method='GET', url=self.url + '/enabled', params={'type': 'RULE'})
        self.signer.add_auth(request)
        prepped_request = request.prepare()

        response = requests.get(prepped_request.url, headers=prepped_request.headers)
        response.raise_for_status()
        parsed_response = json.loads(response.text)
        return parsed_response['policies']
