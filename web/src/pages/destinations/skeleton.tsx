import React from 'react';
import { Card } from 'pouncejs';
import TablePlaceholder from 'Components/table-placeholder';

const DestinationsPageSkeleton: React.FC = () => {
  return (
    <Card p={9}>
      <TablePlaceholder />
    </Card>
  );
};

export default DestinationsPageSkeleton;
