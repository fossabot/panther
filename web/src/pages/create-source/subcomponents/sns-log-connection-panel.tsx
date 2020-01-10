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

import { Button, Text, Flex, Box, Heading } from 'pouncejs';
import React from 'react';
import FormikTextInput from 'Components/fields/text-input';
import { Field } from 'formik';
import { PANTHER_LOG_PROCESSING_ROLE } from 'Source/constants';

// Super important for all these links to have no space or indents in each new line
// The params for the cloudformation are passed in via query parameters using param_<Param_Name>
export const logProcessingCloudformationLink = `https://${process.env.AWS_REGION}.console.aws.amazon.com/cloudformation/home?\
region=${process.env.AWS_REGION}#/stacks/create/review?templateURL=https://s3-us-west-2.amazonaws.com/\
panther-public-cloudformation-templates/${PANTHER_LOG_PROCESSING_ROLE}/latest/template.yml&\
stackName=${PANTHER_LOG_PROCESSING_ROLE}`;

const SnsLogConnectionPanel: React.FC = () => {
  const [isStackLaunched, markStackAsLaunched] = React.useState(false);

  return !isStackLaunched ? (
    <Box>
      <Heading size="medium" m="auto" mb={10} color="grey400">
        Grant us permission to read
      </Heading>
      <Text size="large" color="grey200" is="p">
        By clicking the button below, you will be redirected to the CloudFormation console to launch
        a stack in your account.
        <br />
        <br />
        This stack will create a ReadOnly IAM Role used to read gathered logs that are accumulated
        into the S3 buckets that you specify. By default this role will be able to read logs from
        all your S3 buckets, but you can limit that through the template parameter{' '}
        <b>S3ObjectArns</b>.
        <br />
        <br />
        If your logs are encrypted, please provide the encryption keys via the template parameter{' '}
        <b>KmsKeys</b>
      </Text>
      <Flex mt={6}>
        <Button
          onClick={() => markStackAsLaunched(true)}
          size="large"
          variant="default"
          is="a"
          target="_blank"
          rel="noopener noreferrer"
          href={logProcessingCloudformationLink}
        >
          Launch Stack
        </Button>
      </Flex>
    </Box>
  ) : (
    <Box width={460} m="auto">
      <Heading size="medium" m="auto" mb={10} color="grey400">
        Let us know which role to assume
      </Heading>
      <Text size="large" color="grey200" is="p" mb={10}>
        When you have finished deploying the stack, please fill in the ARN value of the Role that we
        need to assume in order to access your logs on S3.
        <br />
        <br />
        You can find this value as an output of the Cloudformation stack you deployed, under the key{' '}
        <b>RoleArn</b>
      </Text>
      <Field
        name="logProcessingRoleArn"
        as={FormikTextInput}
        label="Assumable Role ARN"
        placeholder="The ARN of the role to read logs from S3"
        aria-required
        mb={6}
      />
    </Box>
  );
};

export default SnsLogConnectionPanel;
