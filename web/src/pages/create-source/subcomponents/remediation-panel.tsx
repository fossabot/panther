/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React from 'react';
import { Box, Heading, Text } from 'pouncejs';
import AddRemediationLambdaForm from 'Components/forms/add-remediation-lambda-form';
import SetupRemediationForm from 'Components/forms/setup-remediation-form';
import {
  PANTHER_REMEDIATION_MASTER_ACCOUNT,
  PANTHER_REMEDIATION_SATELLITE_ACCOUNT,
} from 'Source/constants';

export const adminRemediationCloudformationLink = `https://${process.env.AWS_REGION}.console.aws.amazon.com/cloudformation/home?\
region=${process.env.AWS_REGION}#/stacks/create/review?templateURL=https://s3-us-west-2.amazonaws.com/\
panther-public-cloudformation-templates/${PANTHER_REMEDIATION_MASTER_ACCOUNT}/\
latest/template.yml&stackName=${PANTHER_REMEDIATION_MASTER_ACCOUNT}`;

export const getSatelliteRemediationCloudformationLink = (masterAWSAccountId: string) => {
  return `https://us-west-2.console.aws.amazon.com/cloudformation/home?\
region=us-west-2#/stacks/create/review?templateURL=https://s3-us-west-2.amazonaws.com/\
panther-public-cloudformation-templates/${PANTHER_REMEDIATION_SATELLITE_ACCOUNT}/latest/template.yml&\
stackName=${PANTHER_REMEDIATION_SATELLITE_ACCOUNT}&param_MasterAccountId=${masterAWSAccountId}`;
};

const RemediationPanel: React.FC = () => {
  const [isStackLaunched, markStackAsLaunched] = React.useState(false);

  return (
    <Box>
      <Heading size="medium" m="auto" mb={10} color="grey400">
        Setup AWS Automatic Remediation (Optional)
      </Heading>
      <Text size="large" color="grey200" mb={6} is="p">
        By clicking the button below, you will be redirected to the CloudFormation console to launch
        a stack in your account.
        <br />
        <br />
        This stack will configure Panther to fix misconfigured infrastructure as soon as it is
        detected. Remediations can be configured on a per-policy basis to take any desired actions.
        After a successful deployment, you will have to come back to this page to save the ARN of
        the created lambda. You will be able to edit it afterwards through your Organization{"'"}s
        settings page.
        <br />
        <br />
        If you need more information on the process, please visit our{' '}
        <a
          target="_blank"
          rel="noopener noreferrer"
          href="https://docs.runpanther.io/amazon-web-services/aws-setup/automatic-remediation"
        >
          documentation
        </a>{' '}
        to learn more about this functionality.
      </Text>

      {isStackLaunched ? (
        <AddRemediationLambdaForm />
      ) : (
        <SetupRemediationForm
          onStackLaunch={() => markStackAsLaunched(true)}
          getStackUrl={({ isSatellite, adminAWSAccountId }) =>
            isSatellite
              ? getSatelliteRemediationCloudformationLink(adminAWSAccountId)
              : adminRemediationCloudformationLink
          }
        />
      )}
    </Box>
  );
};

export default RemediationPanel;
