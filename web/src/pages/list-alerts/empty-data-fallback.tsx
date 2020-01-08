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
import SecurityCheckImg from 'Assets/illustrations/security-check.svg';
import { Box, Flex, Heading, Text } from 'pouncejs';

const ListAlertsPageEmptyDataFallback: React.FC = () => {
  return (
    <Flex
      height="100%"
      width="100%"
      justifyContent="center"
      alignItems="center"
      flexDirection="column"
    >
      <Box m={10}>
        <img
          alt="Shield with checkmark illustration"
          src={SecurityCheckImg}
          width="auto"
          height={350}
        />
      </Box>
      <Heading size="medium" color="grey400" mb={6}>
        It{"'"}s quiet in here
      </Heading>
      <Text size="large" color="grey200" textAlign="center" mb={10}>
        Any suspicious rule-based activity we detect will be listed here
      </Text>
    </Flex>
  );
};

export default ListAlertsPageEmptyDataFallback;
