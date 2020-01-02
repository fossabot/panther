import React from 'react';
import SecurityCheckImg from 'Assets/illustrations/security-check.svg';
import { Box, Flex, Heading, Text } from 'pouncejs';

const ListAlertsPageEmptyDataFallback: React.FC = () => {
  return (
    <Flex
      height="100%"
      width="100%"
      justifyContent="center"
      alignItems="center"
      flexDirection="column"
    >
      <Box m={10}>
        <img
          alt="Shield with checkmark illustration"
          src={SecurityCheckImg}
          width="auto"
          height={350}
        />
      </Box>
      <Heading size="medium" color="grey400" mb={6}>
        It{"'"}s quiet in here
      </Heading>
      <Text size="large" color="grey200" textAlign="center" mb={10}>
        Any suspicious rule-based activity we detect will be listed here
      </Text>
    </Flex>
  );
};

export default ListAlertsPageEmptyDataFallback;
