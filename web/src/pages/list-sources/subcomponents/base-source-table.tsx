import React from 'react';
import { Alert, Table, TableProps } from 'pouncejs';
import { IntegrationsByOrganizationResponse, Integration } from 'Generated/schema';
import TablePlaceholder from 'Components/table-placeholder';
import { extractErrorMessage } from 'Helpers/utils';
import { QueryResult } from '@apollo/client';
import { INTEGRATION_TYPES } from 'Source/constants';

interface BaseSourceTableProps {
  query: QueryResult<{ integrations: IntegrationsByOrganizationResponse }, {}>;
  columns: TableProps<Integration>['columns'];

  // TODO: Remove this prop once the backend properly filters by `integrationType`
  integrationType: INTEGRATION_TYPES;
}

const BaseSourceTable: React.FC<BaseSourceTableProps> = ({ query, columns, integrationType }) => {
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

  // FIXME: This filtering is only needed because the backend doesn't do it. It should be removed
  // once the backend is capable of filtering by `integrationType`
  const items = data.integrations.integrations.filter(i => i.integrationType === integrationType);
  return (
    <Table<Integration> items={items} getItemKey={item => item.integrationId} columns={columns} />
  );
};

export default BaseSourceTable;
