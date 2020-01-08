/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
