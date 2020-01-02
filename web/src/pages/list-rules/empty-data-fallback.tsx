import React from 'react';
import { Box, Flex, Heading, Text } from 'pouncejs';
import EmptyNotepadImg from 'Assets/illustrations/empty-notepad.svg';
import RuleCreateButton from './subcomponents/create-button';

const ListRulesPageEmptyDataFallback: React.FC = () => {
  return (
    <Flex justifyContent="center" alignItems="center" flexDirection="column">
      <Box my={10}>
        <img alt="Empty Notepad illustration" src={EmptyNotepadImg} width="auto" height={300} />
      </Box>
      <Heading size="medium" color="grey300" mb={6}>
        No rules found
      </Heading>
      <Text size="large" color="grey200" textAlign="center" mb={10}>
        Writing rules will allow you to get alerts about suspicious activity in your system
      </Text>
      <RuleCreateButton />
    </Flex>
  );
};

export default ListRulesPageEmptyDataFallback;
