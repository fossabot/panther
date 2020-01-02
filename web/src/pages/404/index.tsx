import React from 'react';
import { Flex, Heading, Text, Button, Box } from 'pouncejs';
import { Link } from 'react-router-dom';
import NotFoundImg from 'Assets/illustrations/not-found.svg';

const Page404: React.FC = () => {
  return (
    <Flex
      justifyContent="center"
      alignItems="center"
      width="100vw"
      height="100vh"
      position="fixed"
      left={0}
      top={0}
      bg="white"
      flexDirection="column"
    >
      <Box mb={10}>
        <img alt="Page not found illustration" src={NotFoundImg} width="auto" height={400} />
      </Box>
      <Heading size="medium" color="grey300" mb={4}>
        Not all who wander are lost...
      </Heading>
      <Text size="large" color="grey200" is="p" mb={10}>
        ( You definitely are though )
      </Text>
      <Button size="small" variant="default" is={Link} to="/">
        Back to Home
      </Button>
    </Flex>
  );
};

export default Page404;
