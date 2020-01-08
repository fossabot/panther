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
import { AlertSummary } from 'Generated/schema';
import { generateEnumerationColumn } from 'Helpers/utils';
import { Table } from 'pouncejs';
import columns from 'Pages/list-alerts/columns';

interface ListAlertsTableProps {
  items?: AlertSummary[];
  enumerationStartIndex?: number;
}

const ListAlertsTable: React.FC<ListAlertsTableProps> = ({ items }) => {
  const enumeratedColumns = [generateEnumerationColumn(0), ...columns];
  return (
    <Table<AlertSummary>
      columns={enumeratedColumns}
      getItemKey={alert => alert.alertId}
      items={items}
    />
  );
};

export default React.memo(ListAlertsTable);
