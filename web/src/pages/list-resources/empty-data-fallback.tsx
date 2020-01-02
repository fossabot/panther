import React from 'react';
import { Box, Button, Flex, Heading, Text } from 'pouncejs';
import BlankCanvasImg from 'Assets/illustrations/blank-canvas.svg';
import urls from 'Source/urls';
import { Link } from 'react-router-dom';
import { INTEGRATION_TYPES } from 'Source/constants';

const ListResourcesPageEmptyDataFallback: React.FC = () => {
  return (
    <Flex justifyContent="center" alignItems="center" flexDirection="column">
      <Box my={10}>
        <img alt="Black Canvas Illustration" src={BlankCanvasImg} width="auto" height={300} />
      </Box>
      <Heading size="medium" color="grey300" mb={6}>
        No resources found
      </Heading>
      <Text size="large" color="grey200" textAlign="center" mb={10}>
        You don{"'"}t have any resources connected to your Panther account
      </Text>
      <Button
        size="large"
        variant="primary"
        to={urls.account.settings.sources.create(INTEGRATION_TYPES.AWS_INFRA)}
        is={Link}
      >
        Get started
      </Button>
    </Flex>
  );
};

export default ListResourcesPageEmptyDataFallback;
