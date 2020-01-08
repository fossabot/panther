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
import { Box, Grid, Flex } from 'pouncejs';
import Panel from 'Components/panel';
import TablePlaceholder from 'Components/table-placeholder';
import CirclePlaceholder from 'Components/circle-placeholder';
import DonutChartWrapper from 'Pages/overview/subcomponents/donut-chart-wrapper';

const ChartPlaceholder: React.FC = () => (
  <Flex height="100%" alignItems="center" justifyContent="center">
    <CirclePlaceholder size={150} />
  </Flex>
);

const OverviewPageSkeleton: React.FC = () => {
  return (
    <Box is="article" mb={6}>
      <Grid
        gridTemplateColumns="repeat(4, 1fr)"
        gridRowGap={3}
        gridColumnGap={3}
        is="section"
        mb={3}
      >
        <DonutChartWrapper title="Policy Overview" icon="policy">
          <ChartPlaceholder />
        </DonutChartWrapper>
        <DonutChartWrapper title="Policy Failure Breakdown" icon="policy">
          <ChartPlaceholder />
        </DonutChartWrapper>
        <DonutChartWrapper title="Resources Platforms" icon="resource">
          <ChartPlaceholder />
        </DonutChartWrapper>
        <DonutChartWrapper title="Resources Health" icon="resource">
          <ChartPlaceholder />
        </DonutChartWrapper>
      </Grid>
      <Grid gridTemplateColumns="1fr 1fr" gridRowGap={2} gridColumnGap={3}>
        <Panel title="Top Failing Policies" size="small">
          <TablePlaceholder />
        </Panel>
        <Panel title="Top Failing Resources" size="small">
          <TablePlaceholder />
        </Panel>
      </Grid>
    </Box>
  );
};

export default OverviewPageSkeleton;
