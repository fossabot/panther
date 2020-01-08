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
import { Box, Flex, Heading, Text } from 'pouncejs';
import EmptyNotepadImg from 'Assets/illustrations/empty-notepad.svg';
import PolicyCreateButton from './subcomponents/create-button';

const ListPoliciesPageEmptyDataFallback: React.FC = () => {
  return (
    <Flex justifyContent="center" alignItems="center" flexDirection="column">
      <Box my={10}>
        <img alt="Empty Notepad Illustration" src={EmptyNotepadImg} width="auto" height={300} />
      </Box>
      <Heading size="medium" color="grey300" mb={6}>
        No policies found
      </Heading>
      <Text size="large" color="grey200" textAlign="center" mb={10}>
        Writing policies is the only way to secure your infrastructure against misconfigurations
      </Text>
      <PolicyCreateButton />
    </Flex>
  );
};

export default ListPoliciesPageEmptyDataFallback;
