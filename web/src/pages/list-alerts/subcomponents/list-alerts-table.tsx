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
