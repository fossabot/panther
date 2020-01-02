import React from 'react';
import { Text, Box, Heading, Button, Flex } from 'pouncejs';
import { PANTHER_REAL_TIME } from 'Source/constants';

/*
https://s3-us-west-2.amazonaws.com/panther-public-cloudformation-templates/panther-cloudwatch-events/latest/template.yml
 */

export const getAdminRealTimeCloudformationLink = () => {
  return `https://${process.env.AWS_REGION}.console.aws.amazon.com/cloudformation/home?\
region=${process.env.AWS_REGION}#/stacks/create/review?templateURL=https://s3-${process.env.AWS_REGION}.amazonaws.com/\
panther-public-cloudformation-templates/${PANTHER_REAL_TIME}/latest/\
template.yml&stackName=${PANTHER_REAL_TIME}`;
};

const RealTimeEventPanel: React.FC = () => {
  return (
    <Box>
      <Heading size="medium" m="auto" mb={10} color="grey400">
        Setup Real-Time AWS Resource Scans (Optional)
      </Heading>
      <Text size="large" color="grey200" mb={6} is="p">
        By clicking the button below, you will be redirected to the CloudFormation console to launch
        a stack in your account.
        <br />
        <br />
        This stack will configure Panther to track real-time changes of your AWS Account resources
        when they are created, modified, or deleted. This ensures Panther can detect potential
        security issues as fast as possible. Please visit our{' '}
        <a
          target="_blank"
          rel="noopener noreferrer"
          href="https://docs.runpanther.io/amazon-web-services/aws-setup/real-time-events"
        >
          documentation
        </a>{' '}
        to learn more about this functionality.
      </Text>
      <Flex mt={6}>
        <Button
          size="large"
          variant="default"
          target="_blank"
          is="a"
          rel="noopener noreferrer"
          href={getAdminRealTimeCloudformationLink()}
        >
          Launch Stack
        </Button>
      </Flex>
    </Box>
  );
};

export default RealTimeEventPanel;
