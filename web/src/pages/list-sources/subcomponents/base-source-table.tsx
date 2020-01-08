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
