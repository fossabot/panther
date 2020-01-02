import React from 'react';
import TablePlaceholder from 'Components/table-placeholder';
import { Card } from 'pouncejs';

const ListPoliciesPageSkeleton: React.FC = () => {
  return (
    <Card p={9}>
      <TablePlaceholder />
    </Card>
  );
};

export default ListPoliciesPageSkeleton;
