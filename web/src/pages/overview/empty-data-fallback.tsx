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
import { Box, Button, Flex, Heading, Text } from 'pouncejs';
import EmptyDataImg from 'Assets/illustrations/empty-box.svg';
import { Link } from 'react-router-dom';
import urls from 'Source/urls';
import { INTEGRATION_TYPES } from 'Source/constants';

const OverviewPageEmptyDataFallback: React.FC = () => (
  <Flex
    height="100%"
    width="100%"
    justifyContent="center"
    alignItems="center"
    flexDirection="column"
  >
    <Box m={10}>
      <img alt="Empty data illustration" src={EmptyDataImg} width="auto" height={400} />
    </Box>
    <Heading size="medium" color="grey400" mb={6}>
      It{"'"}s empty in here
    </Heading>
    <Text size="large" color="grey200" textAlign="center" mb={10}>
      You don{"'"}t seem to have any sources connected to our system. <br />
      When you do, a high level overview of your system{"'"}s health will appear here.
    </Text>
    <Button
      size="large"
      variant="primary"
      is={Link}
      to={urls.account.settings.sources.create(INTEGRATION_TYPES.AWS_INFRA)}
    >
      Add your first source
    </Button>
  </Flex>
);

export default OverviewPageEmptyDataFallback;
