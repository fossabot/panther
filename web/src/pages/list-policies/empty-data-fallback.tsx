import React from 'react';
import { Box, Flex, Heading, Text } from 'pouncejs';
import EmptyNotepadImg from 'Assets/illustrations/empty-notepad.svg';
import PolicyCreateButton from './subcomponents/create-button';

const ListPoliciesPageEmptyDataFallback: React.FC = () => {
  return (
    <Flex justifyContent="center" alignItems="center" flexDirection="column">
      <Box my={10}>
        <img alt="Empty Notepad Illustration" src={EmptyNotepadImg} width="auto" height={300} />
      </Box>
      <Heading size="medium" color="grey300" mb={6}>
        No policies found
      </Heading>
      <Text size="large" color="grey200" textAlign="center" mb={10}>
        Writing policies is the only way to secure your infrastructure against misconfigurations
      </Text>
      <PolicyCreateButton />
    </Flex>
  );
};

export default ListPoliciesPageEmptyDataFallback;
