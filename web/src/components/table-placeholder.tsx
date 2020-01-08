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

interface TablePlaceholderProps {
  /** The number of rows that the placeholder component should render. Defaults to 5 */
  rowCount?: number;

  /** The height of each row. Defaults to 10px */
  rowHeight?: number;

  /** The vertical gap between each row. Defaults to 5px */
  rowGap?: number;
}

const TablePlaceholder: React.FC<TablePlaceholderProps> = ({
  rowCount = 5,
  rowHeight = 10,
  rowGap = 5,
}) => (
  <ContentLoader height={rowCount * (rowHeight + rowGap)}>
    {[...Array(rowCount)].map((__, index) => (
      <rect
        key={index}
        x="0"
        y={index * (rowHeight + rowGap)}
        rx="1"
        ry="1"
        width="100%"
        height={rowHeight}
      />
    ))}
  </ContentLoader>
);

export default TablePlaceholder;
