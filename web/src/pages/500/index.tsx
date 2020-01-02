import React from 'react';
import { Box, Button, Flex, Heading, Text } from 'pouncejs';
import { Link } from 'react-router-dom';
import WarningImg from 'Assets/illustrations/warning.svg';

const Page500: React.FC = () => {
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
        <img alt="Page crash illustration" src={WarningImg} width="auto" height={350} />
      </Box>
      <Heading size="medium" color="grey300" mb={4}>
        Something went terribly wrong
      </Heading>
      <Text size="medium" color="grey200" is="p" mb={10}>
        This would normally be an internal server error, but we are fully serverless. Feel free to
        laugh.
      </Text>
      <Button size="small" variant="default" is={Link} to="/">
        Back to somewhere stable
      </Button>
    </Flex>
  );
};

export default Page500;
