import React from 'react';
import { Box, Button, Flex, Heading, Text } from 'pouncejs';
import EmptyDataImg from 'Assets/illustrations/empty-box.svg';
import { Link } from 'react-router-dom';
import urls from 'Source/urls';
import { INTEGRATION_TYPES } from 'Source/constants';

const OverviewPageEmptyDataFallback: React.FC = () => (
  <Flex
    height="100%"
    width="100%"
    justifyContent="center"
    alignItems="center"
    flexDirection="column"
  >
    <Box m={10}>
      <img alt="Empty data illustration" src={EmptyDataImg} width="auto" height={400} />
    </Box>
    <Heading size="medium" color="grey400" mb={6}>
      It{"'"}s empty in here
    </Heading>
    <Text size="large" color="grey200" textAlign="center" mb={10}>
      You don{"'"}t seem to have any sources connected to our system. <br />
      When you do, a high level overview of your system{"'"}s health will appear here.
    </Text>
    <Button
      size="large"
      variant="primary"
      is={Link}
      to={urls.account.settings.sources.create(INTEGRATION_TYPES.AWS_INFRA)}
    >
      Add your first source
    </Button>
  </Flex>
);

export default OverviewPageEmptyDataFallback;
