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
import TablePlaceholder from 'Components/table-placeholder';
import { Box, Card } from 'pouncejs';
import Panel from 'Components/panel';

const PolicyDetailsPageSkeleton: React.FC = () => {
  return (
    <React.Fragment>
      <Card p={6}>
        <TablePlaceholder rowCount={2} rowHeight={10} />
      </Card>
      <Box mt={2} mb={6}>
        <Panel size="large" title="Resources">
          <Box mt={6}>
            <TablePlaceholder />
          </Box>
        </Panel>
      </Box>
    </React.Fragment>
  );
};

export default PolicyDetailsPageSkeleton;
