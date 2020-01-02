import React from 'react';
import { Box, Flex } from 'pouncejs';
import Navigation from 'Components/navigation';
import Header from 'Components/header';

const Layout: React.FC = ({ children }) => {
  return (
    <Flex minHeight="100%" bg="white">
      <Navigation />
      <Box is="main" minHeight={1} flex="1 0 auto" bg="grey50">
        <Box width={1214} mx="auto">
          <Header />
          <Box mt={6}>{children}</Box>
        </Box>
      </Box>
    </Flex>
  );
};

export default Layout;
