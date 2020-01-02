import React from 'react';
import DonutChart from 'Components/donut-chart';
import { ScannedResources } from 'Generated/schema';

interface ResourcesByPlatformProps {
  resources: ScannedResources;
}

const ResourcesByPlatform: React.FC<ResourcesByPlatformProps> = ({ resources }) => {
  const allResourcesChartData = [
    {
      value: resources.byType.length,
      label: 'AWS',
      color: 'grey500' as const,
    },
  ];

  return (
    <DonutChart data={allResourcesChartData} renderLabel={(data, index) => data[index].value} />
  );
};

export default React.memo(ResourcesByPlatform);
