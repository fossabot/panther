import React from 'react';
import { Box, Flex, Heading, SideSheet, Text } from 'pouncejs';
import DestinationCard from 'Components/destination-card';
import useSidesheet from 'Hooks/useSidesheet';
import slackLogo from 'Assets/slack-minimal-logo.svg';
import msTeamsLogo from 'Assets/ms-teams-minimal-logo.svg';
import opsgenieLogo from 'Assets/opsgenie-minimal-logo.svg';
import githubLogo from 'Assets/github-minimal-logo.svg';
import pagerDutyLogo from 'Assets/pagerduty-minimal-logo.svg';
import jiraLogo from 'Assets/jira-minimal-logo.svg';
import emailLogo from 'Assets/email-minimal-logo.svg';
import snsLogo from 'Assets/aws-sns.svg';
import sqsLogo from 'Assets/aws-sqs.svg';

import { SIDESHEETS } from 'Components/utils/sidesheet-context';
import { DestinationTypeEnum } from 'Generated/schema';

const destinationConfigs = [
  {
    logo: slackLogo,
    title: 'Slack',
    destinationType: DestinationTypeEnum.Slack,
  },
  {
    logo: msTeamsLogo,
    title: 'Microsoft Teams',
    destinationType: DestinationTypeEnum.Msteams,
  },
  {
    logo: opsgenieLogo,
    title: 'Opsgenie',
    destinationType: DestinationTypeEnum.Opsgenie,
  },
  {
    logo: jiraLogo,
    title: 'Jira',
    destinationType: DestinationTypeEnum.Jira,
  },
  {
    logo: githubLogo,
    title: 'Github',
    destinationType: DestinationTypeEnum.Github,
  },
  {
    logo: pagerDutyLogo,
    title: 'PagerDuty',
    destinationType: DestinationTypeEnum.Pagerduty,
  },
  {
    logo: emailLogo,
    title: 'Email',
    destinationType: DestinationTypeEnum.Email,
  },
  {
    logo: snsLogo,
    title: 'AWS SNS',
    destinationType: DestinationTypeEnum.Sns,
  },
  {
    logo: sqsLogo,
    title: 'AWS SQS',
    destinationType: DestinationTypeEnum.Sqs,
  },
];

export const SelectDestinationSidesheet: React.FC = () => {
  const { hideSidesheet, showSidesheet } = useSidesheet();

  return (
    <SideSheet open onClose={hideSidesheet}>
      <Box width={465}>
        <Box mb={8}>
          <Heading size="medium" mb={8}>
            Select an Alert Destination
          </Heading>
          <Text size="large" color="grey400">
            Add a new destination below to deliver alerts to a specific application for further
            triage
          </Text>
        </Box>
        <Flex justifyContent="space-between" flexWrap="wrap">
          {destinationConfigs.map(destinationConfig => (
            <Box width={224} mb={4} key={destinationConfig.title}>
              <DestinationCard
                logo={destinationConfig.logo}
                title={destinationConfig.title}
                onClick={() =>
                  showSidesheet({
                    sidesheet: SIDESHEETS.ADD_DESTINATION,
                    props: {
                      destinationType: destinationConfig.destinationType,
                    },
                  })
                }
              />
            </Box>
          ))}
        </Flex>
      </Box>
    </SideSheet>
  );
};

export default SelectDestinationSidesheet;
