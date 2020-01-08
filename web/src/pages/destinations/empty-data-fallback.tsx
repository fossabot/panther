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
