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
import DonutChart from 'Components/donut-chart';
import { ScannedResources } from 'Generated/schema';
import { countResourcesByStatus } from 'Helpers/utils';

interface ResourcesByStatusChartProps {
  resources: ScannedResources;
}

const ResourcesByStatusChart: React.FC<ResourcesByStatusChartProps> = ({ resources }) => {
  const totalResources = countResourcesByStatus(resources, ['fail', 'error', 'pass']);

  const failingResourcesChartData = [
    {
      value: countResourcesByStatus(resources, ['fail', 'error']),
      label: 'Failing',
      color: 'red200' as const,
    },
    {
      value: countResourcesByStatus(resources, ['pass']),
      label: 'Passing',
      color: 'green100' as const,
    },
  ];

  return (
    <DonutChart
      data={failingResourcesChartData}
      renderLabel={(chartData, index) => {
        const { value: statusGroupingValue } = chartData[index];
        const percentage = Math.round((statusGroupingValue * 100) / totalResources).toFixed(0);

        return `${statusGroupingValue}\n{small|${percentage}% of all}`;
      }}
    />
  );
};

export default React.memo(ResourcesByStatusChart);
