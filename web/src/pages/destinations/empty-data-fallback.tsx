import React from 'react';
import DestinationImg from 'Assets/illustrations/destination.svg';
import { Box, Flex, Heading, Text } from 'pouncejs';
import DestinationCreateButton from './subcomponents/create-button';

const DestinationsPageEmptyDataFallback: React.FC = () => {
  return (
    <Flex
      height="100%"
      width="100%"
      justifyContent="center"
      alignItems="center"
      flexDirection="column"
    >
      <Box m={10}>
        <img alt="Mobile & Envelope illustration" src={DestinationImg} width="auto" height={350} />
      </Box>
      <Heading size="medium" color="grey400" mb={6}>
        Help us reach you
      </Heading>
      <Text size="large" color="grey200" textAlign="center" mb={10}>
        You don{"'"}t seem to have any destinations setup yet. <br />
        Adding destinations will help you get notified when irregularities occur.
      </Text>
      <DestinationCreateButton />
    </Flex>
  );
};

export default DestinationsPageEmptyDataFallback;
