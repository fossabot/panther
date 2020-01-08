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
import ContentLoader from 'react-content-loader';
import { Box, Card } from 'pouncejs';

interface PanelPlaceholderProps {
  /** The number of rows that the placeholder component should render. Defaults to 5 */
  rowCount?: number;

  /** The height of each row. Defaults to 15px */
  rowHeight?: number;
}

const PanelPlaceholder: React.FC<PanelPlaceholderProps> = ({ rowCount = 4, rowHeight = 15 }) => (
  <Card
    width={1}
    borderBottom="1px solid"
    borderColor="grey100"
    py={8}
    px={8}
    backgroundColor="#fff"
  >
    <Box pb={8} borderBottom="1px solid" borderColor="grey100">
      <ContentLoader height={10}>
        <rect x="0" y="0" rx="1" ry="1" width="30%" height="10" />
      </ContentLoader>
    </Box>
    <Box mt={8}>
      <ContentLoader height={rowCount * rowHeight}>
        {[...Array(rowCount)].map((__, index) => (
          <rect key={index} x="0" y={index * rowHeight} rx="1" ry="1" width="40%" height="10" />
        ))}
      </ContentLoader>
    </Box>
  </Card>
);

export default PanelPlaceholder;
