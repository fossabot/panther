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
