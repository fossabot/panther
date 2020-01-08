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
import { Box, Flex, Heading, Card, Label } from 'pouncejs';

interface PanelProps {
  title: string;
  size: 'small' | 'large';
  actions?: React.ReactNode;
}

const Panel: React.FC<PanelProps> = ({ title, actions, size, children }) => {
  return (
    <Card
      is="section"
      width={1}
      borderBottom="1px solid"
      borderColor="grey100"
      p={size === 'large' ? 8 : 6}
    >
      <Flex
        pb={size === 'large' ? 8 : 6}
        borderBottom="1px solid"
        borderColor="grey100"
        justifyContent="space-between"
        alignItems="center"
      >
        {size === 'large' ? (
          <Heading size="medium" is="h2">
            {title}
          </Heading>
        ) : (
          <Label size="large" is="h4">
            {title}
          </Label>
        )}
        {actions}
      </Flex>
      <Box mt={size === 'large' ? 8 : 6}>{children}</Box>
    </Card>
  );
};

export default Panel;
