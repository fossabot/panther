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

import * as React from 'react';
import { css } from '@emotion/core';
import { Box, Card, Text } from 'pouncejs';

interface ItemCardProps {
  logo: string;
  title: string;
  onClick?: () => void;
}

const DestinationCard: React.FunctionComponent<ItemCardProps> = ({ logo, title, onClick }) => (
  <Card
    is="button"
    onClick={onClick}
    css={css`
      cursor: pointer;
      transition: transform 0.15s ease-in-out;
      &:hover {
        transform: scale3d(1.03, 1.03, 1.03);
      }
    `}
  >
    <Box height={92} px={10}>
      <img src={logo} alt={title} style={{ objectFit: 'contain' }} width="100%" height="100%" />
    </Box>
    <Box borderTopStyle="solid" borderTopWidth="1px" borderColor="grey50">
      <Text size="medium" px={4} py={3} color="grey500" textAlign="left">
        {title}
      </Text>
    </Box>
  </Card>
);

export default DestinationCard;
