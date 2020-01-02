import React from 'react';
import TablePlaceholder from 'Components/table-placeholder';
import { Card } from 'pouncejs';

const AlertDetailsPageSkeleton: React.FC = () => {
  return (
    <Card p={6}>
      <TablePlaceholder rowCount={2} rowHeight={10} />
    </Card>
  );
};

export default AlertDetailsPageSkeleton;
