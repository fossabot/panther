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
