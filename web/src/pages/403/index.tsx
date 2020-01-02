import React from 'react';
import { Box, Button, Flex, Heading, Text } from 'pouncejs';
import useAuth from 'Hooks/useAuth';
import { Link } from 'react-router-dom';
import AccessDeniedImg from 'Assets/illustrations/authentication.svg';

const Page403: React.FC = () => {
  const { userInfo } = useAuth();

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
        <img alt="Access denied illustration" src={AccessDeniedImg} width="auto" height={400} />
      </Box>
      <Heading size="medium" color="grey300" mb={4}>
        You have no power here, {userInfo ? userInfo.given_name : 'Anonymous'} the Grey
      </Heading>
      <Text size="medium" color="grey200" is="p" mb={10}>
        ( Sarum... Your administrator has restricted your powers )
      </Text>
      <Button size="small" variant="default" is={Link} to="/">
        Back to Shire
      </Button>
    </Flex>
  );
};

export default Page403;
