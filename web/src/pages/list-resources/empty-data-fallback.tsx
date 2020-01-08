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
import BlankCanvasImg from 'Assets/illustrations/blank-canvas.svg';
import urls from 'Source/urls';
import { Link } from 'react-router-dom';
import { INTEGRATION_TYPES } from 'Source/constants';

const ListResourcesPageEmptyDataFallback: React.FC = () => {
  return (
    <Flex justifyContent="center" alignItems="center" flexDirection="column">
      <Box my={10}>
        <img alt="Black Canvas Illustration" src={BlankCanvasImg} width="auto" height={300} />
      </Box>
      <Heading size="medium" color="grey300" mb={6}>
        No resources found
      </Heading>
      <Text size="large" color="grey200" textAlign="center" mb={10}>
        You don{"'"}t have any resources connected to your Panther account
      </Text>
      <Button
        size="large"
        variant="primary"
        to={urls.account.settings.sources.create(INTEGRATION_TYPES.AWS_INFRA)}
        is={Link}
      >
        Get started
      </Button>
    </Flex>
  );
};

export default ListResourcesPageEmptyDataFallback;
