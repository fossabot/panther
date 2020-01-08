import React from 'react';
import { Alert, Table, TableProps } from 'pouncejs';
import { Integration } from 'Generated/schema';
import TablePlaceholder from 'Components/table-placeholder';
import { extractErrorMessage } from 'Helpers/utils';
import { QueryResult } from '@apollo/client';

interface BaseSourceTableProps {
  query: QueryResult<{ integrations: Integration[] }, {}>;
  columns: TableProps<Integration>['columns'];
}

const BaseSourceTable: React.FC<BaseSourceTableProps> = ({ query, columns }) => {
  const { loading, error, data } = query;

  if (loading && !data) {
    return <TablePlaceholder />;
  }

  if (error) {
    return (
      <Alert
        variant="error"
        title="Couldn't load your sources"
        description={
          extractErrorMessage(error) ||
          'There was an error when performing your request, please contact support@runpanther.io'
        }
      />
    );
  }

  return (
    <Table<Integration>
      items={data.integrations}
      getItemKey={item => item.integrationId}
      columns={columns}
    />
  );
};

export default BaseSourceTable;
