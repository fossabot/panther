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
import { Box, Card, Flex, Icon, IconProps, Label } from 'pouncejs';
import ErrorBoundary from 'Components/error-boundary';

interface DonutChartWrapperProps {
  title: string;
  icon: IconProps['type'];
}

const DonutChartWrapper: React.FC<DonutChartWrapperProps> = ({ children, title, icon }) => (
  <Card p={6} height={340}>
    <Flex alignItems="center" is="header" mb={6} color="grey500">
      <Icon size="small" type={icon} mr={4} />
      <Label size="large" is="h4">
        {title}
      </Label>
    </Flex>
    <Box height={250}>
      <ErrorBoundary>{children}</ErrorBoundary>
    </Box>
  </Card>
);

export default DonutChartWrapper;
